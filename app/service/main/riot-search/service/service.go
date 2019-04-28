package service

import (
	"context"
	"runtime"

	"go-common/app/service/main/riot-search/conf"
	"go-common/app/service/main/riot-search/dao"
	"go-common/library/queue/databus"

	"github.com/ivpusic/grpool"
)

// Service struct
type Service struct {
	c       *conf.Config
	dao     *dao.Dao
	databus *databus.Databus
	pool    *grpool.Pool
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:       c,
		dao:     dao.New(c),
		databus: databus.New(c.Databus),
		pool:    grpool.NewPool(runtime.NumCPU(), 10240),
	}
	if c.Riot.LoadPath != "" {
		s.loadproc(c.Riot.LoadPath)
	}
	go s.watcherproc()
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
