package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/interface/live/push-live/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/sync/errgroup"
	"sync"

	"github.com/pkg/errors"
)

// LiveStartMessage 直播开播提醒推送消息
func (s *Service) LiveStartMessage(ctx context.Context, msg *databus.Message) (err error) {
	defer msg.Commit()
	var total int
	// message
	m := new(model.StartLiveMessage)
	if err = json.Unmarshal(msg.Value, &m); err != nil {
		log.Error("[service.start_live|LiveStartMessage] json Unmarshal error(%v)", err)
		return
	}
	task := s.InitPushTask(m)
	midMap := s.GetMids(ctx, task)
	// do push
	total = s.Push(task, midMap)
	// create push task
	go s.CreatePushTask(task, total)
	log.Info("[service.push|LiveStartMessage] start live push done, total(%d), task(%v), model(%v), err(%v)",
		total, task, m, err)
	return
}

// InitPushTask 初始化开播提醒推送task
func (s *Service) InitPushTask(m *model.StartLiveMessage) (task *model.ApPushTask) {
	s.mutex.RLock()
	currentPushTypes := s.pushTypes
	s.mutex.RUnlock()
	// push task model
	task = &model.ApPushTask{
		Type:       model.LivePushType,
		TargetID:   m.TargetID,
		AlertTitle: m.Uname,
		AlertBody:  m.RoomTitle,
		MidSource:  s.getSourceByTypes(currentPushTypes),
		LinkType:   s.c.Push.LinkType,
		LinkValue:  m.LinkValue,
		ExpireTime: m.ExpireTime,
	}
	return task
}

// GetMids 开播提醒，根据配置的策略从不同来源获取需要推送的用户id
func (s *Service) GetMids(c context.Context, task *model.ApPushTask) map[int][]int64 {
	var (
		mutex        sync.Mutex
		group        = errgroup.Group{}
		fans         = make(map[int64]bool)
		fansSP       = make(map[int64]bool)
		midMap       = make(map[int][]int64)
		midBlackList = make(map[int64]bool)
	)

	// 获取黑名单
	mb, err := s.dao.GetBlackList(c, task)
	if err != nil {
		log.Error("[service.start_live|GetMids] get black list error(%v), task(%+v)", err, task)
	} else {
		midBlackList = mb
		log.Info("[service.start_live|GetMids] get black list len(%d), task(%+v)", len(midBlackList), task)
	}

	// try get latest push options and expired time
	s.mutex.RLock()
	currentPushTypes := s.pushTypes
	s.mutex.RUnlock()
	// 开多个协程获取后求并集
	for _, t := range currentPushTypes {
		tp := string(t)
		group.Go(func() (e error) {
			var mFans, mSpe map[int64]bool
			switch tp {
			case model.StrategySwitch:
				// 直播开关
				mFans, mSpe, e = s.GetFansBySwitch(context.TODO(), task.TargetID)
			case model.StrategySpecial:
				// 只获取特别关注
				mFans, mSpe, e = s.dao.Fans(context.TODO(), task.TargetID, model.RelationSpecial)
			case model.StrategyFans:
				// 只获取普通关注
				mFans, mSpe, e = s.dao.Fans(context.TODO(), task.TargetID, model.RelationAttention)
			case model.StrategySwitchSpecial:
				// 只获取特别关注(直播开关中的特别关注)
				mFans, mSpe, e = s.GetFansBySwitchAndSpecial(context.TODO(), task.TargetID)
			default:
				log.Error("[service.mids|GetMids] strategy invalid, type(%s), task(%+v)", tp, task)
				e = fmt.Errorf("[service.mids|GetMids] strategy invalid, type(%s), task(%+v)", tp, task)
				return e
			}
			if e != nil {
				log.Error("[service.mids|GetMids] get mid error(%v), type(%s), task(%+v)", e, tp, task)
				return e
			}

			// 来源之间求并集，并过滤重复出现的id
			// filter by black list
			mutex.Lock()
			for fansID := range mFans {
				if _, ok := midBlackList[fansID]; !ok {
					fans[fansID] = true
				}
			}
			for fansID := range mSpe {
				if _, ok := midBlackList[fansID]; !ok {
					fansSP[fansID] = true
				}
			}
			mutex.Unlock()
			log.Info("[service.mids|GetMids] get mids by type(%s), task(%+v), common(%d), special(%d)",
				tp, task, len(mFans), len(mSpe))
			return e
		})
	}
	group.Wait()

	if len(fansSP) > 0 {
		midMap[model.RelationSpecial] = s.midFilter(fansSP, model.StartLiveBusiness, task)
	}
	if len(fans) > 0 {
		midMap[model.RelationAttention] = s.midFilter(fans, model.StartLiveBusiness, task)
	}
	return midMap
}

// GetFansBySwitch 开播提醒，获取开关mids
func (s *Service) GetFansBySwitch(c context.Context, targetID int64) (fans map[int64]bool, fansSP map[int64]bool, err error) {
	// 获取直播侧开关数据(可能包含普通关注与特别关注)
	m, err := s.dao.GetFansBySwitch(c, targetID)
	if err != nil {
		err = errors.WithStack(err)
		log.Error("[service.mids|GetMidsBySwitch] get switch mids error(%v), targetID(%v)", err, targetID)
		return
	}
	// 区分普通关注与特别关注
	fans, fansSP, err = s.dao.SeparateFans(c, targetID, m)
	return
}

// GetFansBySwitchAndSpecial 开播提醒，获取开关用户与特别关注用户的交集
func (s *Service) GetFansBySwitchAndSpecial(c context.Context, targetID int64) (fans map[int64]bool, fansSP map[int64]bool, err error) {
	// 获取直播侧开关数据(可能包含普通关注与特别关注)
	m, err := s.dao.GetFansBySwitch(c, targetID)
	if err != nil {
		err = errors.WithStack(err)
		log.Error("[service.mids|GetMidsBySwitch] get switch mids error(%v), targetID(%v)", err, targetID)
		return
	}
	// 从开关数据中获取到特别关注的部分
	_, fansSP, err = s.dao.SeparateFans(c, targetID, m)
	return
}

// getSourceByTypes 根据不同的推送策略构造Task.MidSource字段
func (s *Service) getSourceByTypes(types []string) int {
	var source, midSource int
	for _, t := range types {
		switch t {
		case model.StrategySwitch:
			source = model.TaskSourceSwitch
		case model.StrategySpecial:
			source = model.TaskSourceSpecial
		case model.StrategyFans:
			source = model.TaskSourceFans
		case model.StrategySwitchSpecial:
			source = model.TaskSourceSwitchSpe
		default:
			source = 0
		}
		midSource = midSource ^ source
	}
	return midSource
}
