package service

import (
	"context"

	"go-common/app/job/bbq/recall/internal/conf"
	"go-common/app/job/bbq/recall/internal/dao"
	"go-common/library/log"

	"github.com/robfig/cron"
)

// Service struct
type Service struct {
	c    *conf.Config
	dao  *dao.Dao
	sche *cron.Cron
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:    c,
		dao:  dao.New(c),
		sche: cron.New(),
	}
	return s
}

// InitCron .
func (s *Service) InitCron() {
	s.sche.AddFunc("@every 3s", s.HeartBeat)
	s.sche.AddFunc(s.c.Job.ForwardIndex.Schedule, s.GenForwardIndex)
	s.sche.AddFunc(s.c.Job.BloomFilter.Schedule, s.GenBloomFilter)
	s.sche.Start()
}

// RunSrv .
func (s *Service) RunSrv(name string) {
	log.Info("run job{%s}", name)
	switch name {
	case s.c.Job.ForwardIndex.JobName:
		s.GenForwardIndex()
	case s.c.Job.BloomFilter.JobName:
		s.GenBloomFilter()
	default:
		s.HeartBeat()
	}
}

// HeartBeat .
func (s *Service) HeartBeat() {
	log.Info("alive...")
}

// Ping Service
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}
