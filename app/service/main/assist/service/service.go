package service

import (
	"context"
	"sync"

	"go-common/app/service/main/assist/conf"
	"go-common/app/service/main/assist/dao/account"
	"go-common/app/service/main/assist/dao/assist"
	"go-common/app/service/main/assist/dao/message"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// Service assist.
type Service struct {
	c   *conf.Config
	ass *assist.Dao
	acc *account.Dao
	msg *message.Dao
	// databus sub
	relationSub *databus.Databus
	// chan
	cacheChan chan func()
	// wait group
	wg sync.WaitGroup
}

// New get assist service.
func New(c *conf.Config) *Service {
	s := &Service{
		c:   c,
		ass: assist.New(c),
		acc: account.New(c),
		msg: message.New(c),
		// chan
		cacheChan: make(chan func(), 1024),
		// databus
		relationSub: databus.New(c.RelationSub),
	}
	s.wg.Add(1)
	go s.relationConsumer()
	s.wg.Add(1)
	go s.cacheproc()
	return s
}

// Ping service.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.ass.Ping(c); err != nil {
		log.Error("s.ass.Dao.Ping err(%v)", err)
	}
	return
}

// asyncCache add to chan for cache
func (s *Service) asyncCache(f func()) {
	select {
	case s.cacheChan <- f:
	default:
		log.Warn("assist cacheproc chan full")
	}
}

// cacheproc is a routine for execute closure.
func (s *Service) cacheproc() {
	for {
		f := <-s.cacheChan
		f()
	}
}

// Close func
func (s *Service) Close() {
	s.relationSub.Close()
	s.wg.Wait()
}
