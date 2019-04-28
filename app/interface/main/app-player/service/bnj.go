package service

import (
	"context"
	"time"

	arcmdl "go-common/app/interface/main/app-player/model/archive"
	"go-common/library/log"
)

func (s *Service) loadBnjArc() {
	var (
		arcs map[int64]*arcmdl.Info
		err  error
	)
	if arcs, err = s.arcDao.Views(context.Background(), s.c.Bnj.Aids); err != nil {
		log.Error("s.arcDao.Views error(%+v)", err)
		return
	}
	s.bnjArcs = arcs
}

func (s *Service) bnjTickproc() {
	for {
		time.Sleep(time.Duration(s.c.Bnj.Tick))
		s.loadBnjArc()
	}
}
