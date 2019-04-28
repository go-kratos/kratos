package service

import (
	"context"
	"time"

	"go-common/app/interface/main/tv/model"
	"go-common/library/log"
)

// Regions .
func (s *Service) Regions(ctx context.Context) (res []*model.Region, err error) {
	res = s.RegionInfo
	return
}

func (s *Service) loadRegionproc() {
	for {
		time.Sleep(time.Duration(s.conf.Region.StopSpan))
		s.loadRegion()
	}
}

func (s *Service) loadRegion() {
	var (
		m   int64
		err error
		res []*model.Region
	)
	if res, err = s.dao.Regions(ctx); err != nil {
		log.Error("s.dao.Regions error(%v)", err)
	}
	if len(res) != 0 && err == nil {
		s.RegionInfo = res
	}
	if m, err = s.dao.FindLastMtime(ctx); err != nil {
		log.Error("s.dao.FindLastMtime error(%v)", err)
	}
	s.MaxTime = m
}
