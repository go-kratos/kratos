package service

import (
	"context"
	"time"

	"go-common/app/admin/main/filter/conf"
	"go-common/app/admin/main/filter/dao"
	"go-common/app/admin/main/filter/dao/search"
	"go-common/library/cache"
	"go-common/library/log"
)

// Service struct.
type Service struct {
	conf      *conf.Config
	dao       *dao.Dao
	searchDao *search.Dao

	hbaseCh   *cache.Cache
	cacheCh   *cache.Cache
	databusCh *cache.Cache

	eventch chan func()
}

// New new service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf:      c,
		dao:       dao.New(c),
		searchDao: search.New(c),
		eventch:   make(chan func(), 1024),
	}

	// 初始化channel
	s.hbaseCh = cache.New(1, 1024)
	s.cacheCh = cache.New(1, 1024)
	s.databusCh = cache.New(1, 1024)

	go s.expiredproc()
	go s.eventproc()
	return
}

func (s *Service) expiredproc() {
	defer func() {
		if x := recover(); x != nil {
			go s.expiredproc()
		}
	}()
	var err error
	for {
		log.Info("expired check tick (%s)", time.Duration(conf.Conf.Property.ExpiredTick))
		if err = s.expireFilter(); err != nil {
			log.Error("%+v", err)
		}
		time.Sleep(time.Duration(conf.Conf.Property.ExpiredTick))
	}
}

// Ping service ping.
func (s *Service) Ping(c context.Context) (err error) {
	if s.dao != nil {
		err = s.dao.Ping(c)
	}
	return
}
func (s *Service) eventproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.eventproc panic(%v)", x)
			go s.eventproc()
			log.Info("service.eventproc recover")
		}
	}()
	for {
		f := <-s.eventch
		f()
	}
}

func (s *Service) mission(f func()) {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.mission panic(%v)", x)
		}
	}()
	select {
	case s.eventch <- f:
	default:
		log.Error("service.missproc chan full")
	}
}

// Close close service.
func (s *Service) Close() {
	if s.dao != nil {
		s.dao.Close()
	}
}
