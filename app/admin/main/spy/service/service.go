package service

import (
	"context"
	"time"

	"go-common/app/admin/main/spy/conf"
	spydao "go-common/app/admin/main/spy/dao"
	"go-common/app/admin/main/spy/model"
	"go-common/library/log"
)

// Service biz service def.
type Service struct {
	c             *conf.Config
	spyDao        *spydao.Dao
	allEventName  map[int64]string
	loadEventTick time.Duration
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:             c,
		spyDao:        spydao.New(c),
		allEventName:  make(map[int64]string),
		loadEventTick: time.Duration(c.Property.LoadEventTick),
	}
	s.loadeventname()
	go s.loadeventproc()
	return
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.spyDao.Ping(c)
}

// Wait wait all closed.
func (s *Service) Wait() {
}

// Close close all dao.
func (s *Service) Close() {
	s.spyDao.Close()
}

func (s *Service) loadeventname() (err error) {
	var (
		c  = context.TODO()
		es []*model.Event
	)
	es, err = s.spyDao.AllEvent(c)
	if err != nil {
		log.Error("loadeventname allevent error(%v)", err)
		return
	}
	tmp := make(map[int64]string, len(es))
	for _, e := range es {
		tmp[e.ID] = e.NickName
	}
	s.allEventName = tmp
	log.V(2).Info("loadeventname (%v) load success", tmp)
	return
}

func (s *Service) loadeventproc() {
	for {
		time.Sleep(s.loadEventTick)
		s.loadeventname()
	}
}
