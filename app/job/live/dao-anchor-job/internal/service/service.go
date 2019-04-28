package service

import (
	"context"

	"go-common/library/net/trace"

	"go-common/library/log"

	"github.com/robfig/cron"

	"go-common/app/job/live/dao-anchor-job/internal/conf"
	"go-common/app/job/live/dao-anchor-job/internal/dao"
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
		c:    c,
		dao:  dao.New(c),
		cron: cron.New(),
	}
	if err := s.cron.AddFunc(s.c.CoverControl.CoverCron, s.updateKeyFrame); err != nil {
		panic(err)
	}
	if err := s.cron.AddFunc(s.c.Minute3Control.Minute3Cron, s.minuteDataToDB); err != nil {
		panic(err)
	}
	if err := s.cron.AddFunc(s.c.MinuteControl.MinuteCron, s.minuteDataToCacheList); err != nil {
		panic(err)
	}
	s.cron.Start()
	return s
}

// Ping Service
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close Service
func (s *Service) Close() {
	log.Info("Crontab Closed!")
	s.cron.Stop()
	log.Info("Physical Dao Closed!")
	s.dao.Close()
	log.Info("tv-job has been closed.")
}

func GetTraceLogCtx(ctx context.Context, title string) (ctxNew context.Context) {
	t := trace.New(title)
	ctxNew = trace.NewContext(ctx, t)
	return
}
