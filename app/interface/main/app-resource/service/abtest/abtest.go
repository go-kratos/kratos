package abtest

import (
	"context"
	"time"

	"go-common/app/interface/main/app-resource/conf"
	expdao "go-common/app/interface/main/app-resource/dao/abtest"
	"go-common/app/interface/main/app-resource/model/experiment"
	"go-common/app/service/main/resource/model"
	"go-common/library/log"

	farm "github.com/dgryski/go-farm"
)

var (
	_emptyExperiment   = []*experiment.Experiment{}
	_defaultExperiment = map[int8][]*experiment.Experiment{
		model.PlatAndroid: []*experiment.Experiment{
			&experiment.Experiment{
				ID:           10,
				Name:         "默认值",
				Strategy:     "default_value",
				Desc:         "默认值为不匹配处理",
				TrafficGroup: "0",
			},
		},
	}
)

type Service struct {
	dao *expdao.Dao
	// tick
	tick time.Duration
	epm  map[int8][]*experiment.Experiment
	c    *conf.Config
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: expdao.New(c),
		// tick
		tick: time.Duration(c.Tick),
		epm:  map[int8][]*experiment.Experiment{},
	}
	s.loadAbTest()
	go s.tickproc()
	return
}

// TemporaryABTests 临时的各种abtest垃圾需求
func (s *Service) TemporaryABTests(c context.Context, buvid string) (tests *experiment.ABTestV2) {
	id := farm.Hash32([]byte(buvid))
	n := int(id % 100)
	autoPlay := 1
	if n > s.c.ABTest.Range {
		autoPlay = 2
	}
	tests = &experiment.ABTestV2{
		AutoPlay: autoPlay,
	}
	return
}

func (s *Service) Experiment(c context.Context, plat int8, build int) (eps []*experiment.Experiment) {
	if es, ok := s.epm[plat]; ok {
	LOOP:
		for _, ep := range es {
			for _, l := range ep.Limit {
				if model.InvalidBuild(build, l.Build, l.Condition) {
					continue LOOP
				}
			}
			eps = append(eps, ep)
		}
	}
	if eps == nil {
		if es, ok := _defaultExperiment[plat]; ok {
			eps = es
		} else {
			eps = _emptyExperiment
		}
	}
	return
}

// tickproc tick load cache.
func (s *Service) tickproc() {
	for {
		time.Sleep(s.tick)
		s.loadAbTest()
	}
}

func (s *Service) loadAbTest() {
	c := context.TODO()
	lm, err := s.dao.ExperimentLimit(c)
	if err != nil {
		log.Error("s.dao.ExperimentLimit error(%v)", err)
		return
	}
	ids := make([]int64, 0, len(lm))
	for id := range lm {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return
	}
	eps, err := s.dao.ExperimentByIDs(c, ids)
	if err != nil {
		log.Error("s.dao.ExperimentByIDs(%v) error(%v)", ids, err)
		return
	}
	epm := make(map[int8][]*experiment.Experiment, len(eps))
	for _, ep := range eps {
		if l, ok := lm[ep.ID]; ok {
			ep.Limit = l
		}
		epm[ep.Plat] = append(epm[ep.Plat], ep)
	}
	s.epm = epm
}

// AbServer is
func (s *Service) AbServer(c context.Context, buvid, device, mobiAPP, filteredStr string, build int, mid int64) (a interface{}, err error) {
	return s.dao.AbServer(c, buvid, device, mobiAPP, filteredStr, build, mid)
}
