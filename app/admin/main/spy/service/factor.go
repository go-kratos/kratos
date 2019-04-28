package service

import (
	"context"

	"go-common/library/log"

	"go-common/app/admin/main/spy/model"
)

// Factors get all factor.
func (s *Service) Factors(c context.Context) (fs []*model.Factors, err error) {
	var (
		fds []*model.Factor
	)
	gs, err := s.spyDao.Groups(c)
	if err != nil {
		log.Error("spyDao.Groups error(%v)", err)
		return
	}
	if len(gs) == 0 {
		return
	}
	for _, v := range gs {
		fds, err = s.spyDao.Factors(c, v.ID)
		if err != nil {
			log.Error("spyDao.Groups error(%v)", err)
			return
		}
		if len(fds) == 0 {
			continue
		}
		for _, f := range fds {
			f := &model.Factors{
				GroupID:   v.ID,
				GroupName: v.Name,
				NickName:  f.NickName,
				FactorVal: f.FactorVal,
				ID:        f.ID,
			}
			fs = append(fs, f)
		}
	}
	return
}

// UpdateFactor update factor.
func (s *Service) UpdateFactor(c context.Context, fs []*model.Factor, name string) (err error) {
	for _, f := range fs {
		_, err = s.spyDao.UpdateFactor(c, f.FactorVal, f.ID)
		if err != nil {
			log.Error("spyDao.UpdateFactor(%v,%d) error(%v)", f.FactorVal, f.ID, err)
			return
		}
		s.AddLog(c, name, model.UpdateFactor, f)
	}
	return
}

// AddFactor add factor
func (s *Service) AddFactor(c context.Context, f *model.Factor) (err error) {
	ret, err := s.spyDao.AddFactor(c, f)
	if err != nil || ret != 1 {
		log.Error("s.spyDao.AddFactor(%v) error(%v,%d)", f, err, ret)
	}
	return
}

// AddEvent add event
func (s *Service) AddEvent(c context.Context, f *model.Event) (err error) {
	ret, err := s.spyDao.AddEvent(c, f)
	if err != nil || ret != 1 {
		log.Error("s.spyDao.AddEvent(%v) error(%v,%d)", f, err, ret)
	}
	return
}

// AddService add service
func (s *Service) AddService(c context.Context, f *model.Service) (err error) {
	ret, err := s.spyDao.AddService(c, f)
	if err != nil || ret != 1 {
		log.Error("s.spyDao.AddService(%v) error(%v,%d)", f, err, ret)
	}
	return
}

// AddGroup add group
func (s *Service) AddGroup(c context.Context, f *model.FactorGroup) (err error) {
	ret, err := s.spyDao.AddGroup(c, f)
	if err != nil || ret != 1 {
		log.Error("s.spyDao.AddGroup(%v) error(%v,%d)", f, err, ret)
	}
	return
}

// UpdateEventName update event name.
func (s *Service) UpdateEventName(c context.Context, event *model.Event) (err error) {
	ret, err := s.spyDao.UpdateEventName(c, event)
	if err != nil || ret != 1 {
		log.Error("s.spyDao.UpdateEventName(%v) error(%v,%d)", event, err, ret)
	}
	return
}
