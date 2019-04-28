package service

import (
	"context"
	"sync"

	"go-common/app/job/main/card/conf"
	"go-common/app/job/main/card/dao"
	cardCli "go-common/app/service/main/card/api/grpc/v1"
	"go-common/library/queue/databus"
)

const (
	_updateAction  = "update"
	_tableUserInfo = "vip_user_info"
)

// Service struct
type Service struct {
	c           *conf.Config
	waiter      *sync.WaitGroup
	dao         *dao.Dao
	vipConsumer *databus.Databus
	// card service
	cardRPC cardCli.CardClient
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		dao:    dao.New(c),
		waiter: new(sync.WaitGroup),
	}
	cardRPC, err := cardCli.NewClient(c.CardRPC)
	if err != nil {
		panic(err)
	}
	s.cardRPC = cardRPC
	if c.Databus.Vip != nil {
		s.vipConsumer = databus.New(c.Databus.Vip)
		s.waiter.Add(1)
		go s.vipchangeproc()
	}
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	if s.c.Databus.Vip != nil {
		s.vipConsumer.Close()
	}
	s.dao.Close()
}

// Wait wait all chan close
func (s *Service) Wait() {
	s.waiter.Wait()
}
