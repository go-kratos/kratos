package service

import (
	"context"
	"sort"

	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/param"
	"go-common/library/log"
)

// AddEvent will add a event
func (s *Service) AddEvent(c context.Context, ep *param.EventParam) (eid int64, err error) {
	e := &model.Event{
		Cid:         ep.Cid,
		AdminID:     ep.AdminID,
		Content:     ep.Content,
		Attachments: ep.Attachments,
		Event:       ep.Event,
	}

	if err = s.dao.ORM.Create(e).Error; err != nil {
		log.Error("Failed to create event(%v): %v", e, err)
		return
	}
	eid = e.Eid
	s.task(func() {
		var c *model.Chall
		if c, err = s.dao.Chall(context.Background(), e.Cid); err != nil {
			log.Error("s.dao.Chall(%d) error(%v)", e.Cid, err)
			err = nil
			return
		}
		s.afterAddReply(ep, c)
	})

	return
}

// BatchAddEvent will add events to batch chall
func (s *Service) BatchAddEvent(c context.Context, bep *param.BatchEventParam) (eids []int64, err error) {
	if len(bep.Cids) <= 0 {
		return
	}
	eids = make([]int64, 0, len(bep.Cids))

	for _, cid := range bep.Cids {
		e := &model.Event{
			Cid:         cid,
			AdminID:     bep.AdminID,
			Content:     bep.Content,
			Attachments: bep.Attachments,
			Event:       bep.Event,
		}
		if err = s.dao.ORM.Create(e).Error; err != nil {
			log.Error("Failed to create event(%v): %v", e, err)
			return
		}
		eids = append(eids, int64(e.Eid))
	}

	s.task(func() {
		var challs map[int64]*model.Chall
		if challs, err = s.dao.Challs(context.Background(), bep.Cids); err != nil {
			log.Error("s.dao.Challs(%v) error(%v)", bep.Cids, err)
			return
		}
		s.afterAddMultiReply(bep, challs)
	})
	return
}

// ListEvent will add a set of events by challenge id
func (s *Service) ListEvent(c context.Context, cid int64) (eventList model.EventSlice, err error) {
	var (
		events map[int64]*model.Event
	)
	if events, err = s.dao.EventsByCid(c, cid); err != nil {
		log.Error("Failed to s.dao.Events(%d): %v", cid, err)
		return
	}

	eventList = make(model.EventSlice, 0, len(events))
	for _, e := range events {
		eventList = append(eventList, e)
	}

	sort.Slice(eventList, func(i, j int) bool {
		return eventList[i].CTime < eventList[j].CTime
	})
	return
}

// batchLastEvent will return the last log on specified targets
func (s *Service) batchLastEvent(c context.Context, cids []int64) (cEvents map[int64]*model.Event, err error) {
	var (
		eids   []int64
		events map[int64]*model.Event
	)

	if eids, err = s.dao.BatchLastEventIDs(c, cids); err != nil {
		log.Error("s.dao.BatchLastEventIDs(%d) error(%v)", cids, err)
		return
	}
	if events, err = s.dao.EventsByIDs(c, eids); err != nil {
		log.Error("s.dao.EventsByIDs(%d) error(%v)", eids, err)
		return
	}

	cEvents = make(map[int64]*model.Event, len(eids))
	for _, e := range events {
		cEvents[e.Cid] = e
	}

	return
}
