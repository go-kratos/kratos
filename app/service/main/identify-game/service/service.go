package service

import (
	"context"

	"go-common/app/service/main/identify-game/conf"
	"go-common/app/service/main/identify-game/dao"
	"go-common/app/service/main/identify-game/model"
	"go-common/library/log"
	"go-common/library/stat"
	"go-common/library/stat/prom"
)

// Service is a identify service.
type Service struct {
	c                  *conf.Config
	d                  *dao.Dao
	missch             chan func()
	regionInfos        []*model.RegionInfo
	dispatcherErrStats stat.Stat
}

// New new a identify service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                  c,
		d:                  dao.New(c),
		missch:             make(chan func(), 10240),
		regionInfos:        c.Dispatcher.RegionInfos,
		dispatcherErrStats: prom.BusinessErrCount,
	}
	go s.cacheproc()
	return
}

// Close dao.
func (s *Service) Close() error {
	return s.d.Close()
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	return s.d.Ping(c)
}

func (s *Service) addCache(f func()) {
	select {
	case s.missch <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc is a routine for executing closure.
func (s *Service) cacheproc() {
	for {
		f := <-s.missch
		f()
	}
}
