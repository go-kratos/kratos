package service

import (
	"context"

	"go-common/app/job/live/recommend-job/internal/conf"
	"go-common/app/job/live/recommend-job/internal/dao"

	"github.com/robfig/cron"
)

// Service struct
type Service struct {
	c    *conf.Config
	dao  *dao.Dao
	cron *cron.Cron
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	return s
}

// RunCrontab ...
func (s *Service) RunCrontab() {
	s.cron = cron.New()
	if s.c.ItemCFJob.Schedule == "" {
		panic("invalid schedule: " + s.c.ItemCFJob.Schedule)
	}
	if s.c.UserAreaJob.Schedule == "" {
		panic("invalid schedule: " + s.c.UserAreaJob.Schedule)
	}
	s.cron.AddJob(s.c.ItemCFJob.Schedule, &ItemCFJob{Conf: s.c.ItemCFJob, RedisConf: s.c.Redis, HadoopConf: s.c.Hadoop})
	s.cron.AddJob(s.c.UserAreaJob.Schedule, &UserAreaJob{JobConf: s.c.UserAreaJob, RedisConf: s.c.Redis, HadoopConf: s.c.Hadoop})
	s.cron.Start()
}

// Ping Service
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}
