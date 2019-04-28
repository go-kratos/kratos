package service

import (
	"context"
	"sync"

	"go-common/app/admin/main/creative/conf"
	"go-common/app/admin/main/creative/dao"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

// Service str
type Service struct {
	conf      *conf.Config
	dao       *dao.Dao
	DB        *gorm.DB
	DBArchive *gorm.DB
	wg        sync.WaitGroup
	asynch    chan func()
	closed    bool
}

// New fn
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf:   c,
		dao:    dao.New(c),
		asynch: make(chan func(), 10240),
	}

	s.DB = s.dao.DB
	s.DBArchive = s.dao.DBArchive
	s.wg.Add(1)
	go s.asynproc()
	return
}

// Ping fn
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
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
		if s.closed {
			return
		}
		f, ok := <-s.asynch
		if !ok {
			return
		}
		f()
	}
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
	s.closed = true
	s.wg.Wait()
}
