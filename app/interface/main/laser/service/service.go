package service

import (
	"context"

	"go-common/app/interface/main/laser/conf"
	"go-common/app/interface/main/laser/dao"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

type Service struct {
	conf       *conf.Config
	dao        *dao.Dao
	pCacheHit  *prom.Prom
	pCacheMiss *prom.Prom
	missch     chan func()
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf:       c,
		dao:        dao.New(c),
		pCacheHit:  prom.CacheHit,
		pCacheMiss: prom.CacheMiss,
		missch:     make(chan func(), 1024),
	}
	go s.cacheproc()
	return
}

func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

func (s *Service) Close() {
	s.dao.Close()
}

// AddCache add to chan for cache
func (s *Service) addCache(f func()) {
	select {
	case s.missch <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc is a routine for execute closure.
func (s *Service) cacheproc() {
	for {
		f := <-s.missch
		f()
	}
}
