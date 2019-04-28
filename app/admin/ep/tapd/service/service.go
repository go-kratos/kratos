package service

import (
	"context"

	"go-common/app/admin/ep/tapd/conf"
	"go-common/app/admin/ep/tapd/dao"
	"go-common/library/sync/pipeline/fanout"

	"github.com/robfig/cron"
)

// Service struct
type Service struct {
	c            *conf.Config
	dao          *dao.Dao
	transferChan *fanout.Fanout
	cron         *cron.Cron
}

// New init.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:            c,
		dao:          dao.New(c),
		transferChan: fanout.New("cache", fanout.Worker(5), fanout.Buffer(10240)),
	}

	if s.c.Scheduler.Active {
		s.cron = cron.New()
		if err := s.cron.AddFunc(c.Scheduler.UpdateHookURLCacheTask, func() { s.dao.SaveEnableHookURLToCache() }); err != nil {
			panic(err)
		}
		s.cron.Start()
	}

	return
}

// Close Service.
func (s *Service) Close() {
	s.dao.Close()
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	err = s.dao.Ping(c)
	return
}
