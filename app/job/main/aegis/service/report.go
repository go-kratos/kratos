package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/aegis/model"
	"go-common/library/log"
)

func (s *Service) reportSubmit(c context.Context, old, new *model.Task) {
	s.reportTaskFinish(c, new)
	stfield := fmt.Sprintf(model.Submit, new.State, old.UID)
	s.dao.IncresByField(c, new.BusinessID, new.FlowID, new.UID, stfield, 1)
	s.dao.IncresByField(c, new.BusinessID, new.FlowID, new.UID, model.UseTime, new.Utime)

	//统计资源的通过，打回什么的，只统计任务列表操作;  异步统计，免得干扰缓存的同步速度
	if old.UID == new.UID && new.State == model.TaskStateSubmit {
		select {
		case s.chanReport <- &model.RIR{
			BizID:  new.BusinessID,
			FlowID: new.FlowID,
			UID:    new.UID,
			RID:    new.RID,
		}:
		case <-time.NewTimer(time.Millisecond * 10).C:
			log.Error("reportSubmit chanfull")
		}
	}
}

func (s *Service) reportResource(c context.Context, bizid, flowid, rid, uid int64) {
	st, err := s.dao.RscState(c, rid)
	if err != nil {
		log.Error("reportResource RscState(%d) error(%v)", rid, err)
		return
	}
	field := fmt.Sprintf(model.RscState, st)
	s.dao.IncresByField(c, bizid, flowid, uid, field, 1)
}

func (s *Service) syncReport(c context.Context) {
	datas, err := s.dao.FlushReport(c)
	if err != nil {
		log.Error("FlushReport error(%v)", err)
		return
	}
	if len(datas) == 0 {
		return
	}

	for key, val := range datas {
		tp, bizid, flowid, uid, err := model.ParseKey(key)
		if err != nil {
			log.Error("syncReport ParseKey(%s)", key)
			continue
		}

		rt := &model.Report{
			BusinessID: int64(bizid),
			FlowID:     int64(flowid),
			UID:        int64(uid),
			TYPE:       tp,
			Content:    val,
		}
		s.dao.Report(c, rt)
	}
}

func (s *Service) reportTaskCreate(c context.Context, new *model.Task) {
	s.dao.IncresTaskInOut(c, new.BusinessID, new.FlowID, "in")
}

func (s *Service) reportTaskFinish(c context.Context, new *model.Task) {
	s.dao.IncresTaskInOut(c, new.BusinessID, new.FlowID, "out")
}
