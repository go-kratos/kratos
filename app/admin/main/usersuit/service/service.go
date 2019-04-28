package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/admin/main/usersuit/conf"
	"go-common/app/admin/main/usersuit/dao"
	"go-common/library/log"
	"go-common/library/queue/databus"

	account "go-common/app/service/main/account/api"
)

// Service struct of service.
type Service struct {
	d *dao.Dao
	// wait group
	wg sync.WaitGroup
	// conf
	c             *conf.Config
	accountClient account.AccountClient
	// databus pub
	accountNotifyPub *databus.Databus
	Managers         map[int64]string
	asynch           chan func()
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                c,
		d:                dao.New(c),
		asynch:           make(chan func(), 102400),
		accountNotifyPub: databus.New(c.AccountNotify),
	}
	var err error
	if s.accountClient, err = account.NewClient(c.AccountGRPC); err != nil {
		panic(err)
	}
	s.loadManager()
	s.wg.Add(1)
	go s.asynproc()
	go s.loadmanagerproc()
	return
}

func (s *Service) loadmanagerproc() {
	for {
		time.Sleep(1 * time.Hour)
		s.loadManager()
	}
}

func (s *Service) loadManager() {
	managers, err := s.d.Managers(context.TODO())
	if err != nil {
		log.Error("s.Managers error(%v)", err)
		return
	}
	s.Managers = managers
}

// Close dao.
func (s *Service) Close() {
	s.d.Close()
	close(s.asynch)
	time.Sleep(1 * time.Second)
	s.wg.Wait()
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	err = s.d.Ping(c)
	return
}

func (s *Service) addAsyn(f func()) {
	select {
	case s.asynch <- f:
	default:
		log.Warn("asynproc chan full")
	}
}

// cacheproc is a routine for executing closure.
func (s *Service) asynproc() {
	defer s.wg.Done()
	for {
		f, ok := <-s.asynch
		if !ok {
			return
		}
		f()
	}
}
