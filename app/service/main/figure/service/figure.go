package service

import (
	"context"

	"go-common/app/service/main/figure/model"
	"go-common/library/log"
)

func (s *Service) FigureBatchInfo(c context.Context, mids []int64) (fs []*model.Figure, err error) {
	if len(mids) == 0 {
		return
	}
	var (
		cache     = true
		missIndex []int
	)
	if fs, missIndex, err = s.dao.FigureBatchInfoCache(c, mids); err != nil {
		cache = false
		log.Error("%+v", err)
	}
	if len(missIndex) == 0 {
		return
	}
	for _, i := range missIndex {
		if fs[i], err = s.dao.FigureInfo(c, mids[i]); err != nil {
			return
		}
	}
	if cache {
		s.addMission(func() {
			var cerr error
			for _, i := range missIndex {
				if fs[i] == nil {
					continue
				}
				if cerr = s.dao.AddFigureInfoCache(context.TODO(), fs[i]); err != nil {
					log.Error("%+v", cerr)
					return
				}
			}
		})
	}
	return
}
