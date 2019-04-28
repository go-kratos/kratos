package service

import (
	"context"

	conf "go-common/app/interface/main/kvo/conf"
	"go-common/app/interface/main/kvo/dao"
	"go-common/library/log"

	"go-common/library/stat/prom"
)

// Service kvo main service
type Service struct {
	da        *dao.Dao
	docLimit  int
	sp        *prom.Prom
	cacheUcCh chan *cacheUc
}

type cacheUc struct {
	mid         int64
	moduleKeyID int
}

// New get a kvo service
func New(c *conf.Config) *Service {
	da := dao.New(c)
	s := &Service{
		da: da,
		//  limit data size
		docLimit:  c.Rule.DocLimit,
		cacheUcCh: make(chan *cacheUc, 1024),
		sp:        prom.New().WithCounter("conf_cache", []string{"method"}),
	}
	go s.cacheUcProc()
	return s
}

func (s *Service) updateUcCache(mid int64, moduleKeyID int) {
	select {
	case s.cacheUcCh <- &cacheUc{
		mid:         mid,
		moduleKeyID: moduleKeyID,
	}:
	default:
		log.Info("s.cacheUcCh is full")
	}
}

func (s *Service) cacheUcProc() {
	for cuc := range s.cacheUcCh {
		uc, err := s.da.UserConf(context.Background(), cuc.mid, cuc.moduleKeyID)
		if err != nil {
			log.Error("service.cacheUcProc(%v,%v),err:%v", cuc.mid, cuc.moduleKeyID)
			continue
		}
		if uc != nil {
			s.da.SetUserConfCache(context.Background(), uc)
		}
	}
}

// Ping kvo service check
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.da.Ping(ctx)
}
