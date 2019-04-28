package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/block/conf"
	"go-common/app/job/main/block/dao"
	"go-common/library/cache"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/pkg/errors"
)

// Service struct
type Service struct {
	dao    *dao.Dao
	cache  *cache.Cache
	missch chan func()

	creditSub        *databus.Databus
	accountNotifyPub *databus.Databus
}

// New init
func New() (s *Service) {
	s = &Service{
		dao:              dao.New(),
		cache:            cache.New(1, 10240),
		missch:           make(chan func(), 10240),
		creditSub:        databus.New(conf.Conf.Databus.Credit),
		accountNotifyPub: databus.New(conf.Conf.AccountNotify),
	}

	// 自动解禁检查
	if conf.Conf.Property.Flag.ExpireCheck {
		go s.limitcheckproc()
		go s.creditcheckproc()
	}

	// 小黑屋答题状态订阅
	if conf.Conf.Property.Flag.CreditSub {
		go s.creditsubproc()
	}

	go s.missproc()
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) limitcheckproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.limitcheckproc panic(%v)", x)
			go s.limitcheckproc()
			log.Info("service.limitcheckproc recover")
		}
	}()
	for {
		log.Info("limit check start")
		s.limitExpireHandler(context.TODO())
		log.Info("limit check end")
		time.Sleep(time.Duration(conf.Conf.Property.LimitExpireCheckTick))
	}
}

func (s *Service) creditcheckproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("%+v", errors.WithStack(fmt.Errorf("service.creditcheckproc panic(%v)", x)))
			go s.creditcheckproc()
			log.Info("service.creditcheckproc recover")
		}
	}()
	for {
		log.Info("black house check start")
		s.creditExpireHandler(context.TODO())
		log.Info("black house check end")
		time.Sleep(time.Duration(conf.Conf.Property.CreditExpireCheckTick))
	}
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
