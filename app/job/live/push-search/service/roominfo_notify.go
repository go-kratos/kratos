package service

import (
	"encoding/json"
	"go-common/app/job/live/push-search/model"
	"go-common/library/log"

	"context"
	roomV1 "go-common/app/service/live/room/api/liverpc/v1"
	"go-common/library/sync/errgroup"
	"strconv"
)

const (
	_updateAct = "update"
	_insertAct = "insert"
)

func (s *Service) roomInfoNotifyConsumeProc() {
	defer s.waiter.Done()
	for {
		msg, ok := <-s.dao.RoomInfoDataBus.Messages()
		// databus关闭chan导致,服务自杀或异常退出
		if !ok {
			log.Error("roomInfoNotifyConsumeProc closed")
			if err := s.dao.RoomInfoDataBus.Close(); err != nil {
				log.Error("s.dao.RoomInfoDataBus.Close() error(%v)", err)
			}
			return
		}

		m := &message{data: msg}
		p := new(model.ApRoomNotifyInfo)

		if err := json.Unmarshal(msg.Value, p); err != nil {
			msg.Commit()
			log.Error("[RoomInfoDataBus]json.Unmarshal(%s) error(%v)", string(msg.Value), err)
			continue
		}

		if p.Action != _insertAct && p.Action != _updateAct {
			msg.Commit()
			log.Error("[RoomInfoDataBus]Action Invalid error(%v)", p.Action)
			continue
		}

		//判断是否是关注or昵称变更,如果是则跳过,顺便解出新旧map
		isAttentionUpdate, oldMap, newMap, err := isAttentionChange(p.Action, p.Old, p.New)
		if err != nil {
			msg.Commit()
			log.Error("[RoomInfoDataBus]isAttentionChange,json.Unmarshal(old:%s, new:%s) error(%v)", string(p.Old), string(p.New), err)
			continue
		}
		if isAttentionUpdate {
			msg.Commit()
			log.Error("[RoomInfoDataBus]attention change pass")
			continue
		}

		//hash chan

		if newMap == nil || newMap.RoomId <= 0 {
			msg.Commit()
			log.Error("[RoomInfoDataBus]roomId type conversion error, roomId:%+v", newMap)
			continue
		}

		dataMap := new(model.DataMap)
		dataMap.Action = p.Action
		dataMap.Table = p.Table
		dataMap.New = newMap
		dataMap.Old = oldMap
		m.object = dataMap

		log.Info("[RoomInfoDataBus]roomInfoNotifyConsumeProc key:%s partition:%d offset:%d", msg.Key, msg.Partition, msg.Offset)

		s.binLogMergeChan[newMap.RoomId%s.c.Group.RoomInfo.Num] <- m
	}
}

func isAttentionChange(action string, old []byte, new []byte) (bool, *model.TableField, *model.TableField, error) {
	newMap := &model.TableField{}
	oldMap := &model.TableField{}
	err := json.Unmarshal(new, newMap)
	if err != nil {
		return false, oldMap, newMap, err
	}

	if action == _updateAct {
		err := json.Unmarshal(old, oldMap)
		if err != nil {
			return false, oldMap, newMap, err
		}
		if oldMap != nil && oldMap.Attentions != newMap.Attentions {
			return true, oldMap, newMap, err
		}
	}

	if action == _insertAct {
		oldMap = nil
	}

	return false, oldMap, newMap, err
}

func (s *Service) roomInfoNotifyHandleProc(c chan *message) {
	defer s.waiterChan.Done()
	for {
		msgData, ok := <-c
		if !ok {
			log.Error("[RoomInfoDataBus]roomInfoNotifyHandleProc closed")
			return
		}

		msgData.data.Commit()

		p, assertOk := msgData.object.(*model.DataMap)

		if !assertOk {
			log.Error("[RoomInfoDataBus]roomInfoNotifyHandleProc msg object type conversion error, msg:%+v", msgData)
			return
		}

		uid := p.New.Uid

		wg := errgroup.Group{}
		uName := ""
		fc := 0
		areaInfo := &roomV1.AreaGetDetailResp_AreaInfo{}

		wg.Go(func() (err error) {
			userInfo, err := s.getMultiUserInfo(uid)
			if err == nil && userInfo != nil && userInfo.Uname != "" {
				uName = userInfo.Uname
			}
			return
		})

		//fc任何错误都要返回,不然fc为0无法判断是接口返回0还是初始化的0!!!!
		wg.Go(func() (err error) {
			fc, err = s.getFc(uid)
			return
		})

		wg.Go(func() (err error) {
			areaInfo, err = s.getAreaV2Detail(p.New.AreaV2Id)
			return
		})

		err := wg.Wait()

		//成功返回则替换,否则输出原数据
		ret, retByte := s.generateSearchInfo(p.Action, p.Table, p.New, p.Old)
		if err == nil {
			if uName != "" {
				ret["new"].(map[string]interface{})["uname"] = uName
				retByte["uname"] = []byte(uName)
			}
			if areaInfo != nil && areaInfo.Name != "" {
				ret["new"].(map[string]interface{})["s_category"] = areaInfo.Name
				retByte["s_category"] = []byte(areaInfo.Name)
			}
			ret["new"].(map[string]interface{})["attentions"] = fc
			ret["new"].(map[string]interface{})["attention"] = fc
			retByte["attentions"] = []byte(strconv.Itoa(fc))
			retByte["attention"] = []byte(strconv.Itoa(fc))
		}
		writeWg := errgroup.Group{}
		writeWg.Go(func() (err error) {
			for i := 0; i < _retry; i++ {
				hbaseErr := s.saveHBase(context.TODO(), s.rowKey(p.New.RoomId), retByte)
				err = hbaseErr
				if hbaseErr != nil {
					continue
				}
				break
			}
			if err != nil {
				log.Error("[RoomInfoDataBus]fail to write hbase, msg:(%v), err:(%v)", p, err)
			}
			return
		})
		writeWg.Go(func() (err error) {
			err = s.dao.Pub(context.TODO(), int64(p.New.RoomId), ret)
			if err != nil {
				log.Error("[RoomInfoDataBus]fail to pub, msg:(%v), err:(%v)", p, err)
			}
			return
		})
		wg.Wait()
		log.Info("[RoomInfoDataBus]success handle, error(%v),msg:(%v)", err, ret)

	}
}
