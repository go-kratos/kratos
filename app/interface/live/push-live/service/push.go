package service

import (
	"context"
	"go-common/app/interface/live/push-live/model"
	"go-common/library/log"
	"time"
)

// liveMessageConsumeproc Live push related message handler
func (s *Service) liveMessageConsumeproc() {
	defer func() {
		log.Warn("liveMessageConsumeproc exited.")
		s.wg.Done()
	}()
	var (
		liveStartMsgs  = s.liveStartSub.Messages()
		liveCommonMsgs = s.liveCommonSub.Messages()
	)
	for {
		select {
		case msg, ok := <-liveStartMsgs:
			if !ok {
				log.Warn("[service.push|liveMessageConsumeproc] liveStartSub has been closed.")
				return
			}
			log.Info("[service.push|liveMessageConsumeproc] consume liveStartSub key(%s) offset(%d) message(%s)",
				msg.Key, msg.Offset, msg.Value)
			s.LiveStartMessage(context.TODO(), msg)
		case msg, ok := <-liveCommonMsgs:
			if !ok {
				log.Warn("[service.push|liveMessageConsumeproc] liveCommonSub has been closed.")
				return
			}
			log.Info("[service.push|liveMessageConsumeproc] consume liveCommonSub key(%s) offset(%d) message(%s)",
				msg.Key, msg.Offset, msg.Value)
			s.LiveCommonMessage(context.TODO(), msg)
		default:
			time.Sleep(time.Second * 3)
			continue
		}
	}
}

// Push 组装业务参数，调用推送平台接口
func (s *Service) Push(task *model.ApPushTask, midMap map[int][]int64) (total int) {
	var shouldPushCount int
	for t, list := range midMap {
		length := len(list)
		shouldPushCount += length
		if length > 0 {
			// 调用批量推送方法，批量推送逻辑会切分mid与出错重试，最后返回实际推送成功数量
			task.Group = s.GetPushGroup(t, task.Group)
			pushCount := s.dao.BatchPush(&list, task)
			log.Info("[service.push|Push] push type(%d), count(%d), target_id(%v)", t, pushCount, task.TargetID)
			total += pushCount
		}
	}
	if shouldPushCount == 0 {
		log.Info("[service.push|Push] None to push, task(%+v)", task)
		return
	}
	log.Info("[service.push|Push] push done.should(%d), actual(%d), task(%+v).", shouldPushCount, total, task)
	return
}

// CreatePushTask create push task
func (s *Service) CreatePushTask(task *model.ApPushTask, total int) (affected int64, err error) {
	task.Total = total
	affected, err = s.dao.CreateNewTask(context.TODO(), task)
	if err != nil || affected == 0 {
		log.Error("[service.push|CreatePushTask] CreateNewTask error(%v), task(%+v)", err, task)
		return
	}
	log.Info("[service.push|CreatePushTask] CreateNewTask success, task(%+v)", task)
	return
}

// GetPushGroup 获取不同类型的group
// 兼容逻辑: 开播提醒topic有指定的group(并且单次开播需要区分关注与特别关注两个group)，其余common message topic会传group
func (s *Service) GetPushGroup(t int, g string) string {
	var group string
	switch t {
	case model.RelationAttention:
		group = model.AttentionGroup
	case model.RelationSpecial:
		group = model.SpecialGroup
	default:
		group = g
	}
	return group
}
