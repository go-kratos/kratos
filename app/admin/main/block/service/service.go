package service

import (
	"context"

	"go-common/app/admin/main/block/conf"
	"go-common/app/admin/main/block/dao"
	"go-common/library/cache"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// Service struct
type Service struct {
	dao              *dao.Dao
	cache            *cache.Cache
	missch           chan func()
	accountNotifyPub *databus.Databus
}

// New init
func New() (s *Service) {
	s = &Service{
		dao:              dao.New(),
		cache:            cache.New(1, 10240),
		missch:           make(chan func(), 10240),
		accountNotifyPub: databus.New(conf.Conf.AccountNotify),
	}
	go s.missproc()
	return s
}

func (s *Service) missproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.missproc panic(%v)", x)
			go s.missproc()
			log.Info("service.missproc recover")
		}
	}()
	for {
		f := <-s.missch
		f()
	}
}

func (s *Service) mission(f func()) {
	select {
	case s.missch <- f:
	default:
		log.Error("s.missch full")
	}
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}
