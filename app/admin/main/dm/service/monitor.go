package service

import (
	"context"

	"go-common/app/admin/main/dm/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// MonitorList get monitor list
func (s *Service) MonitorList(c context.Context, tp int32, pid, oid, mid int64, state int32, kw, sort, order string, page, size int64) (res *model.MonitorResult, err error) {
	var attr int32
	if state > 0 {
		if state == model.MonitorBefore {
			attr = int32(model.AttrSubMonitorBefore) + 1
		} else {
			attr = int32(model.AttrSubMonitorAfter) + 1
		}
	}
	data, err := s.dao.SearchMonitor(c, tp, pid, oid, mid, attr, kw, sort, order, page, size)
	if err != nil {
		log.Error("dao.SearchMonitor(pid:%d,oid:%d) error(%v)", pid, oid, err)
		return
	}
	res = &model.MonitorResult{
		Order:    data.Order,
		Sort:     data.Sort,
		Page:     data.Page.Num,
		PageSize: data.Page.Size,
		Total:    data.Page.Total,
		Result:   make([]*model.Monitor, 0, len(data.Result)),
	}
	for _, v := range data.Result {
		m := &model.Monitor{
			ID:     v.ID,
			Type:   v.Type,
			Pid:    v.Pid,
			Oid:    v.Oid,
			MCount: v.MCount,
			Ctime:  v.Ctime,
			Mtime:  v.Mtime,
			Mid:    v.Mid,
			Title:  v.Title,
			Author: v.Author,
		}
		if v.Attr>>model.AttrSubMonitorBefore&1 == model.AttrYes {
			m.State = model.MonitorBefore
		} else {
			m.State = model.MonitorAfter
		}
		res.Result = append(res.Result, m)
	}
	return
}

// UpdateMonitor update monitor state of dm subject.
func (s *Service) UpdateMonitor(c context.Context, tp int32, oids []int64, state int32) (affect int64, err error) {
	var wg errgroup.Group
	subs, err := s.dao.Subjects(c, tp, oids)
	if err != nil {
		return
	}
	for _, v := range subs {
		sub := v
		switch state {
		case model.MonitorClosed:
			sub.AttrSet(model.AttrNo, model.AttrSubMonitorBefore)
			sub.AttrSet(model.AttrNo, model.AttrSubMonitorAfter)
		case model.MonitorBefore:
			sub.AttrSet(model.AttrYes, model.AttrSubMonitorBefore)
			sub.AttrSet(model.AttrNo, model.AttrSubMonitorAfter)
		case model.MonitorAfter:
			sub.AttrSet(model.AttrNo, model.AttrSubMonitorBefore)
			sub.AttrSet(model.AttrYes, model.AttrSubMonitorAfter)
		default:
			err = ecode.RequestErr
			return
		}
		wg.Go(func() (err error) {
			aft, err := s.dao.UpSubjectAttr(context.TODO(), tp, sub.Oid, sub.Attr)
			if err != nil {
				return
			}
			affect = affect + aft
			return
		})
	}
	err = wg.Wait()
	return
}

// updateMonitorCnt update mcount of subject.
func (s *Service) updateMonitorCnt(c context.Context, sub *model.Subject) (err error) {
	var state, mcount int64
	if sub.AttrVal(model.AttrSubMonitorBefore) == model.AttrYes {
		state = int64(model.StateMonitorBefore)
	} else if sub.AttrVal(model.AttrSubMonitorAfter) == model.AttrYes {
		state = int64(model.StateMonitorAfter)
	} else {
		return
	}
	if mcount, err = s.dao.DMCount(c, sub.Type, sub.Oid, []int64{state}); err != nil {
		return
	}
	_, err = s.dao.UpSubjectMCount(c, sub.Type, sub.Oid, mcount)
	return
}
