package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/service/main/account-recovery/conf"
	"go-common/app/service/main/account-recovery/dao"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// Service struct
type Service struct {
	c *conf.Config
	//用户申诉dao
	d *dao.Dao

	mailch chan func()
	// wait
	wg sync.WaitGroup

	userActLogPub *databus.Databus
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:             c,
		d:             dao.New(c),
		mailch:        make(chan func(), c.ChanSize.MailMsg),
		userActLogPub: databus.New(c.DataBus.UserActLog),
	}
	s.wg.Add(1)
	go s.mailproc()
	return s
}

// AddMailch .
func (s *Service) AddMailch(f func()) {
	select {
	case s.mailch <- f:
	default:
		log.Warn("AddMailch chan full")
	}
}

// mailproc send mail
func (s *Service) mailproc() {
	defer s.wg.Done()
	for {
		f, ok := <-s.mailch
		if !ok {
			log.Warn("s.mailch chan is close")
			return
		}
		f()
	}
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.d.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.d.Close()
	close(s.mailch)
	time.Sleep(1 * time.Second)
	s.wg.Wait()

}
