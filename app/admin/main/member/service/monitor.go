package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/member/model"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

const (
	_logActionMonitorAdd = "monitor_user_add"
	_logActionMonitorDel = "monitor_user_del"
)

// Monitors is.
func (s *Service) Monitors(ctx context.Context, arg *model.ArgMonitor) ([]*model.Monitor, int, error) {
	includeDeleted := false
	if arg.Mid > 0 {
		includeDeleted = true
	}
	mns, total, err := s.dao.Monitors(ctx, arg.Mid, includeDeleted, arg.Pn, arg.Ps)
	if err != nil {
		return nil, 0, err
	}
	s.monitorsName(ctx, mns)
	return mns, total, nil
}

// AddMonitor is.
func (s *Service) AddMonitor(ctx context.Context, arg *model.ArgAddMonitor) error {
	remark := fmt.Sprintf("加入监控列表：%s", arg.Remark)
	if err := s.dao.AddMonitor(ctx, arg.Mid, arg.Operator, remark); err != nil {
		return err
	}
	report.Manager(&report.ManagerInfo{
		Uname:    arg.Operator,
		UID:      arg.OperatorID,
		Business: model.ManagerLogID,
		Type:     0,
		Oid:      arg.Mid,
		Action:   _logActionMonitorAdd,
		Ctime:    time.Now(),
		// extra
		Index: []interface{}{},
		Content: map[string]interface{}{
			"remark": remark,
		},
	})
	return nil
}

// DelMonitor is.
func (s *Service) DelMonitor(ctx context.Context, arg *model.ArgDelMonitor) error {
	remark := fmt.Sprintf("移出监控列表：%s", arg.Remark)
	if err := s.dao.DelMonitor(ctx, arg.Mid, arg.Operator, remark); err != nil {
		return err
	}
	report.Manager(&report.ManagerInfo{
		Uname:    arg.Operator,
		UID:      arg.OperatorID,
		Business: model.ManagerLogID,
		Type:     0,
		Oid:      arg.Mid,
		Action:   _logActionMonitorDel,
		Ctime:    time.Now(),
		// extra
		Index: []interface{}{},
		Content: map[string]interface{}{
			"remark": remark,
		},
	})
	return nil
}

func (s *Service) monitorsName(ctx context.Context, mns []*model.Monitor) {
	mids := make([]int64, 0, len(mns))
	for _, mn := range mns {
		mids = append(mids, mn.Mid)
	}
	bs, err := s.dao.Bases(ctx, mids)
	if err != nil {
		log.Error("Failed to fetch bases with mids: %+v: %+v", mids, err)
		return
	}
	for _, mn := range mns {
		b, ok := bs[mn.Mid]
		if !ok {
			continue
		}
		mn.Name = b.Name
	}
}
