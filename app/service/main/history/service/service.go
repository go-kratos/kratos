package service

import (
	"context"
	"errors"

	"go-common/app/service/main/history/conf"
	"go-common/app/service/main/history/dao"
	"go-common/app/service/main/history/model"
	"go-common/library/ecode"
	"go-common/library/sync/pipeline"
	"go-common/library/sync/pipeline/fanout"
)

const asyncProcNum = 100

// Service struct
type Service struct {
	c             *conf.Config
	dao           *dao.Dao
	businesses    map[int64]*model.Business
	businessNames map[string]*model.Business
	merge         *pipeline.Pipeline
	asyncChan     chan func()
	cache         *fanout.Fanout
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:         c,
		dao:       dao.New(c),
		cache:     fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		asyncChan: make(chan func(), 1024),
	}
	s.businesses = s.dao.Businesses
	s.businessNames = s.dao.BusinessNames
	s.initMerge()
	for i := 0; i < asyncProcNum; i++ {
		go s.asyncFuncproc()
	}
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.merge.Close()
	s.dao.Close()
}

// checkBusiness .
func (s *Service) checkBusiness(bs string) (err error) {
	if s.businessNames[bs] == nil {
		err = ecode.AppDenied
	}
	return
}

func (s *Service) asyncFuncproc() {
	for {
		fn := <-s.asyncChan
		fn()
	}
}

func (s *Service) asyncFunc(f func()) (err error) {
	select {
	case s.asyncChan <- f:
	default:
		err = errors.New("async full")
	}
	return
}
