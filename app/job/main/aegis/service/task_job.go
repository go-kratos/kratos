package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/aegis/model"
)

func (s *Service) taskProc() {
	for {
		// 检索超时任务，进行释放
		s.dao.TaskRelease(context.Background(), time.Now().Add(-10*time.Minute))
		// 检索过期登陆用户，进行踢出
		s.checkKickOut(context.Background())
		time.Sleep(10 * time.Minute)
		s.syncReport(context.Background())
		s.taskClear()
	}
}

func (s *Service) checkKickOut(c context.Context) {
	s.ccMux.RLock()
	defer s.ccMux.RUnlock()
	for bizfwid, uidm := range s.consumerCache {
		for uid := range uidm {
			pos := strings.Index(bizfwid, "-")
			bizid, _ := strconv.Atoi(bizfwid[:pos])
			flowid, _ := strconv.Atoi(bizfwid[pos+1:])

			if on, err := s.dao.IsConsumerOn(c, bizid, flowid, uid); err == nil && !on {
				delete(s.consumerCache[bizfwid], uid)
				s.KickOut(c, int64(bizid), int64(flowid), uid)
			}
		}
	}
}

// KickOut 踢出过期用户并释放任务
func (s *Service) KickOut(c context.Context, bizid, flowid, uid int64) {
	// 1. 踢出用户
	s.dao.KickOutConsumer(c, int64(bizid), int64(flowid), uid)
	s.sendTaskLog(c, &model.Task{BusinessID: bizid, FlowID: flowid}, model.LogTypeTaskConsumer, "kickout", uid, "")

	// 2. 释放任务
	s.dao.ReleaseByConsumer(c, bizid, flowid, uid)
}

func (s *Service) taskClear() {
	mt := time.Now().Add(-3 * 24 * time.Hour)
	for {
		rows, err := s.dao.TaskClear(context.Background(), mt, 1000)
		if err != nil || rows == 0 {
			break
		}
		time.Sleep(time.Second)
	}
}
