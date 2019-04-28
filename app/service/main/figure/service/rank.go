package service

import (
	"context"
	"sync"

	"go-common/app/service/main/figure/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	rs = &ranks{}
)

type ranks struct {
	rs []*model.Rank
	mu sync.RWMutex
}

func (r *ranks) load(datas []*model.Rank) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rs = make([]*model.Rank, 0, 100)
	r.rs = append(r.rs, datas...)
}

func (r *ranks) find(score int32) (percentage int8) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, r := range r.rs {
		if r.ScoreFrom <= score && r.ScoreTo >= score {
			return r.Percentage
		}
	}
	return 100
}

func (s *Service) loadRank(c context.Context) {
	var (
		ranks []*model.Rank
		err   error
	)
	if ranks, err = s.dao.Ranks(c); err != nil {
		log.Error("%+v", err)
		return
	}
	rs.load(ranks)
}

func (s *Service) Rank(c context.Context, score int32) (percentage int8) {
	percentage = rs.find(score)
	return
}

func (s *Service) FigureWithRank(c context.Context, mid int64) (fr *model.FigureWithRank, err error) {
	var fs []*model.Figure
	if fs, err = s.FigureBatchInfo(c, []int64{mid}); err != nil {
		return
	}
	if len(fs) != 1 || fs[0] == nil {
		err = ecode.FigureNotFound
		return
	}
	fr = s.generateFigureWithRank(c, fs[0])
	return
}

func (s *Service) BatchFigureWithRank(c context.Context, mids []int64) (frs []*model.FigureWithRank, err error) {
	if len(mids) == 0 {
		return
	}
	var fs []*model.Figure
	if fs, err = s.FigureBatchInfo(c, mids); err != nil {
		return
	}
	for _, f := range fs {
		if f == nil {
			continue
		}
		frs = append(frs, s.generateFigureWithRank(c, f))
	}
	return
}

func (s *Service) generateFigureWithRank(c context.Context, f *model.Figure) (fr *model.FigureWithRank) {
	fr = &model.FigureWithRank{Figure: f, Percentage: s.Rank(c, f.Score)}
	return
}
