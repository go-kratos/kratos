package service

import (
	"context"
	"time"

	"go-common/app/admin/main/spy/model"
)

// UpdateState update spy state.
func (s *Service) UpdateState(c context.Context, state int8, id int64, operater string) (err error) {
	var (
		context string
	)
	if _, err = s.spyDao.UpdateStatState(c, state, id); err != nil {
		return
	}
	if state == 0 {
		context = model.WaiteCheck
	} else {
		context = model.HadCheck
	}
	s.AddLog2(c, &model.Log{
		Name:    operater,
		RefID:   id,
		Module:  model.UpdateStat,
		Context: s.UpdateStateLog(id, context),
		Ctime:   time.Now(),
	})
	return
}

// UpdateStatQuantity update stat quantity.
func (s *Service) UpdateStatQuantity(c context.Context, count, id int64, operater string) (err error) {
	var (
		olds *model.Statistics
	)
	if olds, err = s.spyDao.Statistics(c, id); err != nil || olds == nil {
		return
	}
	if _, err = s.spyDao.UpdateStatQuantity(c, count, id); err != nil {
		return
	}
	s.AddLog2(c, &model.Log{
		Name:    operater,
		RefID:   id,
		Module:  model.UpdateStat,
		Context: s.UpdateStatCountLog(id, olds.Quantity, count),
		Ctime:   time.Now(),
	})
	return
}

// DeleteStat delete stat.
func (s *Service) DeleteStat(c context.Context, isdel int8, id int64, operater string) (err error) {
	if _, err = s.spyDao.DeleteStat(c, isdel, id); err != nil {
		return
	}
	s.AddLog2(c, &model.Log{
		Name:    operater,
		RefID:   id,
		Module:  model.UpdateStat,
		Context: s.DeleteStatLog(id),
		Ctime:   time.Now(),
	})
	return
}

// StatPage  get  stat list.
func (s *Service) StatPage(c context.Context, mid, id int64, t int8, pn, ps int) (page *model.StatPage, err error) {
	var (
		list  []*model.Statistics
		count int64
	)
	page = &model.StatPage{Ps: ps, Pn: pn}
	if t == model.AccountType {
		count, err = s.spyDao.StatCountByMid(c, mid)
		if err != nil || count == 0 {
			return
		}
		list, err = s.spyDao.StatListByMid(c, mid, pn, ps)
	} else {
		count, err = s.spyDao.StatCountByID(c, id, t)
		if err != nil || count == 0 {
			return
		}
		list, err = s.spyDao.StatListByID(c, id, t, pn, ps)
	}
	for _, st := range list {
		st.EventName = s.allEventName[st.EventID]
	}
	page.Items = list
	page.TotalCount = count
	return
}
