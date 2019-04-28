package service

import (
	"context"
	"encoding/json"
	"go-common/app/job/live/push-search/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"strconv"
)

const (
	_retry = 3
)

func (s *Service) attentionNotifyConsumeProc() {
	defer s.waiter.Done()
	for {
		msg, ok := <-s.dao.AttentionDataBus.Messages()
		if !ok {
			log.Error("attentionNotifyConsumeProc closed")
			if err := s.dao.AttentionDataBus.Close(); err != nil {
				log.Error("s.dao.AttentionDataBus.Close() error(%v)", err)
			}
			return
		}

		m := &message{data: msg}

		p := new(model.LiveDatabusAttention)

		if err := json.Unmarshal(msg.Value, p); err != nil {
			msg.Commit()
			log.Error("[AttentionDataBus]json.Unmarshal(%s) error(%v)", string(msg.Value), err)
			continue
		}

		if p.MsgContent == nil {
			log.Error("[AttentionDataBus]attentionNotifyConsumeProc msg object msgContent is nil, msg:%+v", string(msg.Value))
			return
		}

		m.object = p
		s.attentionMergeChan[p.MsgContent.UpUid%int64(s.c.Group.Attention.Num)] <- m
	}
}

func (s *Service) attentionNotifyHandleProc(c chan *message) {
	defer s.waiterChan.Done()
	for {
		msgData, ok := <-c
		if !ok {
			log.Error("[AttentionDataBus]attentionNotifyHandleProc closed")
			return
		}
		//先提交防止阻塞,关闭时等待任务执行完
		msgData.data.Commit()

		p, assertOk := msgData.object.(*model.LiveDatabusAttention)

		if !assertOk {
			log.Error("[AttentionDataBus]attentionNotifyHandleProc msg object type conversion error, msg:%+v", msgData)
			return
		}

		uid := p.MsgContent.UpUid
		uName := ""
		newMap := &model.TableField{}

		wg := errgroup.Group{}
		wg.Go(func() (err error) {
			userInfo, err := s.getMultiUserInfo(uid)
			if err == nil && userInfo != nil && userInfo.Uname != "" {
				uName = userInfo.Uname
			}
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
			ret, retByte := s.generateSearchInfo("update", _tableArchive, newMap, nil)
			if uName != "" {
				ret["new"].(map[string]interface{})["uname"] = uName
				retByte["uname"] = []byte(uName)
			}
			if p.MsgContent.ExtInfo != nil {
				ret["new"].(map[string]interface{})["attentions"] = p.MsgContent.ExtInfo.UpUidFans
				ret["new"].(map[string]interface{})["attention"] = p.MsgContent.ExtInfo.UpUidFans
				retByte["attentions"] = []byte(strconv.Itoa(p.MsgContent.ExtInfo.UpUidFans))
				retByte["attention"] = []byte(strconv.Itoa(p.MsgContent.ExtInfo.UpUidFans))
			}
			//构造假old
			ret["old"].(map[string]interface{})["attention"] = 0
			ret["old"].(map[string]interface{})["attentions"] = 0

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
					log.Error("[AttentionDataBus]fail to write hbase, msg:(%v), err:(%v)", p, err)
				}
				return
			})
			wg.Go(func() (err error) {
				err = s.dao.Pub(context.TODO(), int64(newMap.RoomId), ret)
				if err != nil {
					log.Error("[AttentionDataBus]fail to pub, msg:(%v), err:(%v)", p, err)
				}
				return
			})
			wg.Wait()
			log.Info("[AttentionDataBus]success to handle, error(%v), msg:(%v)", err, ret)
			continue
		}
		log.Error("[AttentionDataBus]fail to getData, error(%v),msg:(%v)", err, p)
	}
}
