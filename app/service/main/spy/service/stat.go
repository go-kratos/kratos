package service

import (
	"context"

	"go-common/app/service/main/spy/model"
)

// StatByID spy stat by id or mid.
func (s *Service) StatByID(c context.Context, mid, id int64) (stat []*model.Statistics, err error) {
	if mid != 0 && id != 0 {
		stat, err = s.dao.StatListByIDAndMid(c, mid, id)
	} else {
		if id == 0 {
			stat, err = s.dao.StatListByMid(c, mid)
		} else {
			stat, err = s.dao.StatListByID(c, id)
		}
	}
	if len(stat) == 0 {
		return
	}
	for _, st := range stat {
		st.EventName = s.allEventName[st.EventID]
	}
	return
}

// StatByIDGroupEvent spy stat by id or mid.
func (s *Service) StatByIDGroupEvent(c context.Context, mid, id int64) (res []*model.Statistics, err error) {
	var (
		stat []*model.Statistics
		em   = make(map[int64]*model.Statistics)
	)
	if mid != 0 && id != 0 {
		stat, err = s.dao.StatListByIDAndMid(c, mid, id)
	} else {
		if id == 0 {
			stat, err = s.dao.StatListByMid(c, mid)
		} else {
			stat, err = s.dao.StatListByID(c, id)
		}
	}
	if len(stat) == 0 {
		return
	}
	for _, st := range stat {
		item, ok := em[st.EventID]
		if !ok {
			item = &model.Statistics{EventID: st.EventID, EventName: st.EventName}
		}
		item.Quantity += st.Quantity
		em[st.EventID] = item
	}
	for _, val := range em {
		res = append(res, val)
	}
	return
}
