package goblin

import (
	"context"
	"time"

	"go-common/app/interface/main/tv/model"
	"go-common/app/interface/main/tv/model/goblin"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_pgcLabel = 1
	_ugcLabel = 2
)

func (s *Service) labelsproc() {
	for {
		time.Sleep(time.Duration(s.conf.Cfg.IndexLabel.Fre))
		log.Info("Reload Label Data!")
		s.prepareLabels()
	}
}

// prepareLabels refreshes memory labels
func (s *Service) prepareLabels() {
	pgcLbs, errpgc := s.catLabels(_pgcLabel)
	if errpgc != nil {
		log.Error("loadLabels PGC Err %v", errpgc)
		return
	}
	ugcLbs, errugc := s.catLabels(_ugcLabel)
	if errugc != nil {
		log.Error("loadLabels PGC Err %v", errugc)
		return
	}
	if len(pgcLbs) > 0 && len(ugcLbs) > 0 {
		s.labels = &goblin.IndexLabels{
			PGC: pgcLbs,
			UGC: ugcLbs,
		}
	}
}

// catLabels picks ugc/pgc all categories labels
func (s *Service) catLabels(catType int) (result map[int][]*goblin.TypeLabels, err error) {
	var (
		cats    []int
		typeMap map[int32]*model.ArcType
	)
	result = make(map[int][]*goblin.TypeLabels)
	if catType == _pgcLabel {
		cats = s.conf.Cfg.ZonesInfo.PGCZonesID
	} else if catType == _ugcLabel {
		if typeMap, err = s.arcDao.FirstTypes(); err != nil {
			log.Error("loadLabels ArcDao FirstTypes Err %v", err)
			return
		}
		for k := range typeMap {
			cats = append(cats, int(k))
		}
	} else {
		err = ecode.TvDangbeiWrongType
		return
	}
	if len(cats) == 0 {
		err = ecode.TvDangbeiPageNotExist
		return
	}
	for _, category := range cats {
		if result[category], err = s.loadLabels(catType, category); err != nil {
			log.Error("loadLabels Err %v", err)
			return
		}
	}
	return
}

// getTypeLabel builds an typeLabels object from a slice of labels
func getTypeLabel(in []*goblin.Label) *goblin.TypeLabels {
	typeLbs := &goblin.TypeLabels{}
	typeLbs.FromLabels(in)
	return typeLbs
}

// loadLabels
func (s *Service) loadLabels(catType, category int) (result []*goblin.TypeLabels, err error) {
	var (
		ctx       = context.Background()
		labelMap  = make(map[string][]*goblin.Label)
		labels    []*goblin.Label
		showOrder []string
		cfg       = s.conf.Cfg.IndexLabel
	)
	if labels, err = s.dao.Label(ctx, category, catType); err != nil {
		log.Error("loadLabels Dao Label Cat %d, %d, Err %v", category, catType, err)
		return
	}
	for _, v := range labels { // gather labels by their param
		v.TransYear(cfg)
		labelMap[v.Param] = append(labelMap[v.Param], v)
	}
	if catType == _pgcLabel {
		showOrder = cfg.PGCOrder
	} else {
		showOrder = cfg.UGCOrder
	}
	for _, v := range showOrder {
		if line, ok := labelMap[v]; ok {
			result = append(result, getTypeLabel(line))
			delete(labelMap, v)
		}
	}
	if len(labelMap) > 0 {
		for _, v := range labelMap {
			result = append(result, getTypeLabel(v))
		}
	}
	return
}

// Labels picks label
func (s *Service) Labels(c context.Context, catType, category int) (result []*goblin.TypeLabels, err error) {
	var (
		indexLbs map[int][]*goblin.TypeLabels
		ok       bool
	)
	if catType == _pgcLabel {
		indexLbs = s.labels.PGC
	} else {
		indexLbs = s.labels.UGC
	}
	if result, ok = indexLbs[category]; !ok {
		err = ecode.NothingFound
	}
	return
}
