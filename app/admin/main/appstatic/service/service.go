package service

import (
	"context"
	"sync"

	"go-common/app/admin/main/appstatic/conf"
	"go-common/app/admin/main/appstatic/dao"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

// Service biz service def.
type Service struct {
	c         *conf.Config
	dao       *dao.Dao
	DB        *gorm.DB
	waiter    *sync.WaitGroup
	daoClosed bool  // logic close the dao's DB
	MaxSize   int64 // max size supported for the file to upload
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:         c,
		dao:       dao.New(c),
		daoClosed: false,
		waiter:    new(sync.WaitGroup),
	}
	s.DB = s.dao.DB
	if s.c.Cfg.Storage == "nas" {
		s.MaxSize = 200 * 1024 * 1024 // 200M NAS
	} else {
		s.MaxSize = 20 * 1024 * 1024 // 20M BFS
	}
	return s
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return
}

// Wait wait all closed.
func (s *Service) Wait() {
	if s.dao != nil {
		s.daoClosed = true
		log.Info("Dao is logically closed!")
	}
	log.Info("Wait waiter!")
	s.waiter.Wait()
}

// Close close all dao.
func (s *Service) Close() {
	log.Info("Close Dao physically!")
	s.dao.Close()
	log.Info("Service Closed!")
}
