package plugin

import (
	"context"
	"sort"
	"time"

	"go-common/app/interface/main/app-resource/conf"
	pgdao "go-common/app/interface/main/app-resource/dao/plugin"
	"go-common/app/interface/main/app-resource/model/plugin"
	"go-common/library/log"
)

type Service struct {
	pgDao       *pgdao.Dao
	tick        time.Duration
	pluginCache map[string][]*plugin.Plugin
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		pgDao:       pgdao.New(c),
		tick:        time.Duration(c.Tick),
		pluginCache: map[string][]*plugin.Plugin{},
	}
	s.load()
	go s.loadproc()
	return
}

func (s *Service) Plugin(build, baseCode, seed int, name string) (pg *plugin.Plugin) {
	if build == 0 || seed == 0 || name == "" {
		return
	}
	if ps, ok := s.pluginCache[name]; ok {
		for _, p := range ps {
			if ((p.Policy == 1 && baseCode == p.BaseCode) || p.Policy == 2 && baseCode >= p.BaseCode) && seed%100 <= p.Coverage && build >= p.MinBuild && ((p.MaxBuild == 0) || (p.MaxBuild != 0 && build <= p.MaxBuild)) {
				pg = p
				break
			}
		}
	}
	return
}

// load cache data
func (s *Service) load() {
	psm, err := s.pgDao.All(context.TODO())
	if err != nil {
		log.Error("s.pgDao.All() error(%v)", err)
		return
	}
	pgCache := make(map[string][]*plugin.Plugin, len(psm))
	for name, ps := range psm {
		sort.Sort(plugin.Plugins(ps))
		pgCache[name] = ps
	}
	s.pluginCache = pgCache
}

// cacheproc load cache data
func (s *Service) loadproc() {
	for {
		time.Sleep(s.tick)
		s.load()
	}
}
