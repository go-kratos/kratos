package sidebar

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-common/app/interface/main/app-resource/conf"
	bplusdao "go-common/app/interface/main/app-resource/dao/bplus"
	resdao "go-common/app/interface/main/app-resource/dao/resource"
	whitedao "go-common/app/interface/main/app-resource/dao/white"
	"go-common/app/interface/main/app-resource/model"
	"go-common/app/interface/main/app-resource/model/sidebar"
	resource "go-common/app/service/main/resource/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_initSidebarKey      = "sidebar_%d_%d_%s"
	_defaultLanguageHans = "hans"
	_defaultLanguageHant = "hant"
)

type Service struct {
	c *conf.Config
	//dao
	res  *resdao.Dao
	bdao *bplusdao.Dao
	wdao *whitedao.Dao
	// sidebar
	tick         time.Duration
	sidebarCache map[string][]*sidebar.SideBar
	limitsCahce  map[int64][]*resource.SideBarLimit
	//limit ids
	limitIDs map[int64]struct{}
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:            c,
		res:          resdao.New(c),
		bdao:         bplusdao.New(c),
		wdao:         whitedao.New(c),
		tick:         time.Duration(c.Tick),
		sidebarCache: map[string][]*sidebar.SideBar{},
		limitsCahce:  map[int64][]*resource.SideBarLimit{},
		//limit ids
		limitIDs: map[int64]struct{}{},
	}
	s.loadSidebar()
	s.loadLimit(c.SideBarLimit)
	go s.tickproc()
	return s
}

// SideBar
func (s *Service) SideBar(c context.Context, plat int8, build, module int, mid int64, language string) (ss []*sidebar.SideBar) {
	if key := fmt.Sprintf(_initSidebarKey, plat, module, language); len(s.sidebarCache[fmt.Sprintf(key)]) == 0 || language == "" {
		if model.IsOverseas(plat) {
			key = fmt.Sprintf(_initSidebarKey, plat, module, _defaultLanguageHant)
			if len(s.sidebarCache[fmt.Sprintf(key)]) > 0 {
				language = _defaultLanguageHant
			} else {
				language = _defaultLanguageHans
			}
		} else {
			language = _defaultLanguageHans
		}
	}
	var (
		key    = fmt.Sprintf(_initSidebarKey, plat, module, language)
		verify = map[int64]bool{}
		mutex  sync.Mutex
	)
	if sidebars, ok := s.sidebarCache[key]; ok {
		g, _ := errgroup.WithContext(c)
		for _, v := range sidebars {
			var (
				vid  = v.ID
				vurl = v.WhiteURL
			)
			if vurl != "" && mid > 0 {
				g.Go(func() (err error) {
					var ok bool
					if ok, err = s.wdao.WhiteVerify(context.TODO(), mid, vurl); err != nil {
						log.Error("s.wdao.WhiteVerify uri(%s) error(%v)", vurl, err)
						ok = false
						err = nil
					}
					mutex.Lock()
					verify[vid] = ok
					mutex.Unlock()
					return
				})
			} else if vurl != "" && mid == 0 {
				verify[vid] = false
			}
		}
		g.Wait()
	LOOP:
		for _, v := range sidebars {
			for _, l := range s.limitsCahce[v.ID] {
				if model.InvalidBuild(build, l.Build, l.Condition) {
					continue LOOP
				}
			}
			if verifybool, ok := verify[v.ID]; ok && !verifybool {
				continue LOOP
			}
			ss = append(ss, v)
		}
	}
	return
}

// tickproc tick load cache.
func (s *Service) tickproc() {
	for {
		time.Sleep(s.tick)
		s.loadSidebar()
	}
}

func (s *Service) loadSidebar() {
	sideBars, err := s.res.ResSideBar(context.TODO())
	if err != nil || sideBars == nil {
		log.Error("s.sideDao.SideBar error(%v) or nil", err)
		return
	}
	var (
		tmp = map[int64]struct{}{}
		ss  = map[string][]*sidebar.SideBar{}
	)
	for _, v := range sideBars.SideBar {
		if _, ok := tmp[v.ID]; ok {
			continue
		}
		tmp[v.ID] = struct{}{}
		t := &sidebar.SideBar{}
		t.Change(v)
		key := fmt.Sprintf(_initSidebarKey, t.Plat, t.Module, t.Language)
		ss[key] = append(ss[key], t)
	}
	s.sidebarCache = ss
	s.limitsCahce = sideBars.Limit
	log.Info("loadSidebar cache success")
}

func (s *Service) loadLimit(limit []int64) {
	tmp := map[int64]struct{}{}
	for _, l := range limit {
		tmp[l] = struct{}{}
	}
	s.limitIDs = tmp
}
