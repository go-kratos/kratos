package service

import (
	"context"
	"sort"

	"go-common/app/admin/main/workflow/model"
	"go-common/library/log"
)

// ActivityList will list activities by given conditions
func (s *Service) ActivityList(c context.Context, business int8, cid int64) (acts *model.Activities, err error) {
	var (
		logs   []*model.WLog
		events map[int64]*model.Event
		uids   []int64
		uNames map[int64]string
	)

	if logs, err = s.AllAuditLog(c, cid, []int{model.WLogModuleChallenge}); err != nil {
		log.Error("s.AllAuditLog(%d) error(%v)", cid, err)
	}
	log.Info("audit log cid(%d) logs(%+v)", cid, logs)
	if events, err = s.dao.EventsByCid(c, cid); err != nil {
		log.Error("Failed to s.dao.EventsByCid(%d): %v", cid, err)
		return
	}
	acts = new(model.Activities)
	acts.Events = make([]*model.Event, 0, len(events))
	acts.Logs = logs
	for _, e := range events {
		acts.Events = append(acts.Events, e)
		uids = append(uids, e.AdminID)
	}

	if uNames, err = s.dao.BatchUNameByUID(c, uids); err != nil {
		log.Error("s.dao.SearchUNameByUid(%v) error(%v)", uids, err)
		err = nil
	} else {
		for i := range acts.Events {
			acts.Events[i].Admin = uNames[acts.Events[i].AdminID]
		}
	}

	//sort.Sort(model.LogSlice(acts.Logs))
	//sort.Sort(model.EventSlice(acts.Events))
	sort.Slice(acts.Logs, func(i, j int) bool {
		return acts.Logs[i].CTime < acts.Logs[j].CTime
	})
	sort.Slice(acts.Events, func(i, j int) bool {
		return acts.Events[i].CTime < acts.Events[j].CTime
	})
	return
}
