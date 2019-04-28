package service

import (
	"context"
	"time"

	"go-common/app/service/main/figure/conf"
	figureDao "go-common/app/service/main/figure/dao"
	"go-common/library/log"
)

// Service biz service def.
type Service struct {
	c      *conf.Config
	dao    *figureDao.Dao
	missch chan func()
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		dao:    figureDao.New(c),
		missch: make(chan func(), 1024),
	}
	go s.cacheproc()
	go s.rankproc()
	return s
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close close all dao.
func (s *Service) Close() {
	s.dao.Close()
	return
}

func (s *Service) addMission(f func()) {
	select {
	case s.missch <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

func (s *Service) rankproc() {
	for {
		s.loadRank(context.TODO())
		time.Sleep(time.Duration(s.c.Property.LoadRankPeriod))
	}
}

func (s *Service) cacheproc() {
	for {
		f := <-s.missch
		f()
	}
}
