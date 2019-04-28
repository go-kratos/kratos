package service

import (
	"context"
	"encoding/json"
	"go-common/app/interface/live/push-live/dao"
	"go-common/app/interface/live/push-live/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// LiveCommonMessage 直播通用消息
func (s *Service) LiveCommonMessage(ctx context.Context, msg *databus.Message) (err error) {
	defer msg.Commit()
	var (
		mids   []int64
		mMap   = make(map[int64]bool)  // mid去重
		midMap = make(map[int][]int64) // 最终格式化后的mid map
	)
	m := new(model.LiveCommonMessage)
	if err = json.Unmarshal(msg.Value, &m); err != nil {
		log.Error("[service.common_message|LiveCommonMessage] json Unmarshal error(%v), model(%v)", err, m)
		return
	}
	task := s.InitCommonTask(m)
	if mids, err = s.convertStrToInt64(m.MsgContent.Mids); err != nil {
		log.Error("[service.push|LiveCommonMessage] format Mids error(%v), task(%v), model(%v)", err, task, m)
		return
	}
	// remove duplicated mid
	for _, mid := range mids {
		mMap[mid] = true
	}
	// mid filter
	business := m.MsgContent.Business
	filteredMids := s.midFilter(mMap, business, task)
	midMap[business] = filteredMids
	log.Info("[service.push|LiveCommonMessage] message info: before(%d), after(%d), model(%v), task(%v)",
		len(mMap), len(midMap[business]), m, task)
	total := s.Push(task, midMap)
	// create push task
	go s.CreatePushTask(task, total)
	go s.setPushInterval(business, s.safeGetExpired(), filteredMids, task)
	log.Info("[service.push|LiveCommonMessage] common message push done, total(%d), err(%v)", total, err)
	return
}

// InitCommonTask Init push task by common message model
func (s *Service) InitCommonTask(m *model.LiveCommonMessage) (task *model.ApPushTask) {
	task = &model.ApPushTask{
		Type:       model.LivePushType,
		TargetID:   0,
		AlertTitle: m.MsgContent.AlertTitle,
		AlertBody:  m.MsgContent.AlertBody,
		MidSource:  m.MsgContent.Business,
		LinkType:   m.MsgContent.LinkType,
		LinkValue:  m.MsgContent.LinkValue,
		ExpireTime: m.MsgContent.ExpireTime,
		Group:      m.MsgContent.Group,
	}
	return task
}

// setPushInterval 活动预约，对每个mid设置推送平滑key
func (s *Service) setPushInterval(business int, expired int32, mids []int64, task *model.ApPushTask) (total int, err error) {
	if business != 111 {
		return
	}
	var conn redis.Conn
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	// redis conn
	conn, err = redis.Dial(s.c.Redis.PushInterval.Proto, s.c.Redis.PushInterval.Addr, s.dao.RedisOption()...)
	if err != nil {
		log.Error("[service.common_message|setPushInterval] redis.Dial error(%v), task(%v), mids(%d)",
			err, task, len(mids))
		return
	}
	for _, mid := range mids {
		key := dao.GetIntervalKey(mid)
		_, err = conn.Do("SET", key, task.LinkValue, "EX", expired)
		if err != nil {
			log.Error("[service.common_message|setPushInterval] set redis error(%v), task(%v), mid(%d)",
				err, task, mid)
			continue
		}
		total++
	}
	return
}
