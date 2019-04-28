package service

import (
	"context"
	"encoding/json"
	"go-common/app/job/live/push-search/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"strconv"
)

func (s *Service) unameNotifyConsumeProc() {
	defer s.waiter.Done()
	for {
		msg, ok := <-s.dao.UserNameDataBus.Messages()
		if !ok {
			log.Error("unameNotifyConsumeProc closed")
			if err := s.dao.UserNameDataBus.Close(); err != nil {
				log.Error("s.dao.UserNameDataBus.Close() error(%v)", err)
			}
			return
		}
		//先提交防止阻塞,关闭时等待任务执行完
		m := &message{data: msg}
		raw := new(model.LiveDatabus)

		if err := json.Unmarshal(msg.Value, raw); err != nil {
			msg.Commit()
			log.Error("[UnameDataBus]json.Unmarshal(%s) error(%v)", string(msg.Value), err)
			continue
		}

		p := new(model.UnameNotifyInfo)
		if err := json.Unmarshal([]byte(raw.MsgContent), p); err != nil {
			msg.Commit()
			log.Error("[UnameDataBus]json.Unmarshal(%s) error(%v)", raw.MsgContent, err)
			continue
		}

		m.object = p
		s.unameMergeChan[p.Uid%int64(s.c.Group.UserInfo.Num)] <- m
	}
}

func (s *Service) unameNotifyHandleProc(c chan *message) {
	defer s.waiterChan.Done()
	for {
		msgData, ok := <-c
		if !ok {
			log.Error("[UnameDataBus]unameNotifyHandleProc closed")
			return
		}
		//先提交防止阻塞,关闭时等待任务执行完
		msgData.data.Commit()

		p, assertOk := msgData.object.(*model.UnameNotifyInfo)

		if !assertOk {
			log.Error("[UnameDataBus]unameNotifyHandleProc msg object type conversion error, msg:%+v", msgData)
			return
		}

		uid := p.Uid

		if uid == 0 {
			log.Error("[UnameDataBus]empty uid, uid:%d", uid)
			continue
		}
		fc := 0
		newMap := &model.TableField{}

		wg := errgroup.Group{}
		wg.Go(func() (err error) {
			fc, err = s.getFc(uid)
			return
		})

		wg.Go(func() (err error) {
			roomInfo, err := s.getBaseRoomInfo(uid)
			if err == nil && roomInfo != nil {
				newMap.RoomId = int(roomInfo.Roomid)
				newMap.ShortId = int(roomInfo.ShortId)
				newMap.Uid = roomInfo.Uid
				newMap.UName = roomInfo.Uname
				newMap.Area = int(roomInfo.Area)
				newMap.Title = roomInfo.Title
				newMap.Tag = roomInfo.Tags
				newMap.TryTime = roomInfo.TryTime
				newMap.Cover = roomInfo.Cover
				newMap.UserCover = roomInfo.UserCover
				newMap.LockStatus = roomInfo.LockStatus
				newMap.HiddenStatus = roomInfo.HiddenStatus
				newMap.Attentions = int(roomInfo.Attentions)
				newMap.Online = int(roomInfo.Online)
				newMap.LiveTime = roomInfo.LiveTime
				newMap.AreaV2Id = int(roomInfo.AreaV2Id)
				newMap.AreaV2ParentId = int(roomInfo.AreaV2ParentId)
				newMap.Virtual = int(roomInfo.Virtual)
				newMap.AreaV2Name = roomInfo.AreaV2Name
				newMap.CTime = roomInfo.Ctime
				newMap.MTime = roomInfo.Mtime
				newMap.RoundStatus = int(roomInfo.RoundStatus)
				newMap.OnFlag = int(roomInfo.OnFlag)
			}
			return
		})

		err := wg.Wait()

		if err == nil && newMap.RoomId != 0 {
			//非uname更新
			if p.Uname == newMap.UName {
				log.Info("[UnameDataBus]uname no change, msg:(%v)", p)
				continue
			}
			ret, retByte := s.generateSearchInfo("update", _tableArchive, newMap, nil)
			if p.Uname != "" {
				ret["new"].(map[string]interface{})["uname"] = p.Uname
				retByte["uname"] = []byte(p.Uname)
			}
			ret["new"].(map[string]interface{})["attentions"] = fc
			ret["new"].(map[string]interface{})["attention"] = fc
			retByte["attentions"] = []byte(strconv.Itoa(fc))
			retByte["attention"] = []byte(strconv.Itoa(fc))
			ret["old"].(map[string]interface{})["uname"] = ""

			wg := errgroup.Group{}
			wg.Go(func() (err error) {
				for i := 0; i < _retry; i++ {
					hbaseErr := s.saveHBase(context.TODO(), s.rowKey(newMap.RoomId), retByte)
					err = hbaseErr
					if hbaseErr != nil {
						continue
					}
					break
				}
				if err != nil {
					log.Error("[UnameDataBus]fail to write hbase, msg:(%v), err:(%v)", p, err)
				}
				return
			})
			wg.Go(func() (err error) {
				err = s.dao.Pub(context.TODO(), int64(newMap.RoomId), ret)
				if err != nil {
					log.Error("[UnameDataBus]fail to pub, msg:(%v), err:(%v)", p, err)
				}
				return
			})
			wg.Wait()
			log.Info("[UnameDataBus]success to handle, error(%v), msg:(%v)", err, ret)
			continue
		}
		log.Error("[UnameDataBus]fail to getData, error(%v),msg:(%v)", err, p)
	}
}
