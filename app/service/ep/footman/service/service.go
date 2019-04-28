package service

import (
	"context"

	"go-common/app/service/ep/footman/conf"
	"go-common/app/service/ep/footman/dao"
	"go-common/library/cache"

	"github.com/robfig/cron"
)

// Service struct
type Service struct {
	c          *conf.Config
	dao        *dao.Dao
	cache      *cache.Cache
	buglyCache *cache.Cache
	cron       *cron.Cron
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:          c,
		dao:        dao.New(c),
		cache:      cache.New(1, 1024000),
		buglyCache: cache.New(1, 10240),
	}

	if c.Scheduler == nil {
		return
	}
	scheduler := c.Scheduler
	s.cron = cron.New()
	if err := s.cron.AddFunc(scheduler.SaveTapdTime, s.SaveFilesForTask); err != nil {
		panic(err)
	}
	s.cron.Start()

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
