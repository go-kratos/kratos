package service

import (
	"context"
	"sync"

	"go-common/app/job/main/search/conf"
	"go-common/app/job/main/search/dao/base"
	"go-common/app/job/main/search/model"
	"go-common/library/log"
)

var (
	ctx = context.TODO()
)

const (
	_bulkSize = 5000
)

// Service .
type Service struct {
	c *conf.Config
	// base
	base *base.Base
	//mutex
	mutex *sync.RWMutex
	// stats
	stats map[string]*model.Stat
}

// New .
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:     c,
		base:  base.NewBase(c),
		mutex: new(sync.RWMutex),
		stats: make(map[string]*model.Stat),
	}
	s.incrproc()
	return
}

// incrproc incr data
func (s *Service) incrproc() {
	for appid, e := range s.base.D.AppPool {
		if !s.base.D.BusinessPool[appid].IncrOpen {
			continue
		}
		if e.Business() == s.c.Business.Env && !s.c.Business.Index {
			go s.incr(ctx, e)
		}
	}
}

// Close .
func (s *Service) Close() {
	s.base.D.Close()
}

// Ping .
func (s *Service) Ping(c context.Context) error {
	return s.base.D.Ping(c)
}

// HTTPAction http action
func (s *Service) HTTPAction(ctx context.Context, appid, action string, recoverID int64, writeEntityIndex bool) (msg string, err error) {
	switch action {
	case "repair":
	case "all":
		if _, ok := s.base.D.AppPool[appid]; !ok {
			msg = "appid不在appPool中"
			log.Error("AppPool inclueds (%v)", s.base.D.AppPool)
			return
		}
		s.base.D.SetRecover(ctx, appid, recoverID, "", 0)
		go s.all(context.Background(), appid, writeEntityIndex)
	default:
		return
	}
	return
}

// Stat .
func (s *Service) Stat(ctx context.Context, appid string) (st *model.Stat, err error) {
	st = s.stat(appid)
	return
}
