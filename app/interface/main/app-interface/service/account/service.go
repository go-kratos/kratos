package space

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/app-interface/conf"
	accdao "go-common/app/interface/main/app-interface/dao/account"
	audiodao "go-common/app/interface/main/app-interface/dao/audio"
	bplusdao "go-common/app/interface/main/app-interface/dao/bplus"
	memberdao "go-common/app/interface/main/app-interface/dao/member"
	paydao "go-common/app/interface/main/app-interface/dao/pay"
	reldao "go-common/app/interface/main/app-interface/dao/relation"
	sidedao "go-common/app/interface/main/app-interface/dao/sidebar"
	"go-common/app/interface/main/app-interface/model"
	resmodel "go-common/app/service/main/resource/model"
	"go-common/library/log"
)

// Service is space service
type Service struct {
	c            *conf.Config
	accDao       *accdao.Dao
	audioDao     *audiodao.Dao
	relDao       *reldao.Dao
	bplusDao     *bplusdao.Dao
	payDao       *paydao.Dao
	memberDao    *memberdao.Dao
	sideDao      *sidedao.Dao
	sectionCache map[string][]*SectionItem
	white        map[int8][]*SectionURL
	redDot       map[int8][]*SectionURL
}

// SectionItem is
type SectionItem struct {
	Item  *resmodel.SideBar
	Limit []*resmodel.SideBarLimit
}

// SectionURL is
type SectionURL struct {
	ID  int64
	URL string
}

// CheckLimit is
func (item *SectionItem) CheckLimit(build int) bool {
	for _, l := range item.Limit {
		if model.InvalidBuild(build, l.Build, l.Condition) {
			return false
		}
	}
	return true
}

// New new space
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:            c,
		accDao:       accdao.New(c),
		audioDao:     audiodao.New(c),
		relDao:       reldao.New(c),
		bplusDao:     bplusdao.New(c),
		payDao:       paydao.New(c),
		memberDao:    memberdao.New(c),
		sideDao:      sidedao.New(c),
		sectionCache: map[string][]*SectionItem{},
		white:        map[int8][]*SectionURL{},
		redDot:       map[int8][]*SectionURL{},
	}
	if err := s.loadSidebar(); err != nil {
		panic(fmt.Sprintf("load sidebar error(%+v)", err))
	}
	go s.tickproc()
	return s
}

// tickproc tick load cache.
func (s *Service) tickproc() {
	for {
		time.Sleep(time.Duration(s.c.Tick))
		s.loadSidebar()
	}
}

func (s *Service) loadSidebar() (err error) {
	var sidebar *resmodel.SideBars
	if sidebar, err = s.sideDao.Sidebars(context.TODO()); err != nil {
		log.Error("s.sideDao.SideBars error(%v)", err)
		return
	}
	ss := make(map[string][]*SectionItem)
	white := make(map[int8][]*SectionURL)
	redDot := make(map[int8][]*SectionURL)
	for _, item := range sidebar.SideBar {
		item.Plat = s.convertPlat(item.Plat)
		key := fmt.Sprintf(_initSidebarKey, item.Plat, item.Module)
		ss[key] = append(ss[key], &SectionItem{
			Item:  item,
			Limit: sidebar.Limit[item.ID],
		})
		if item.WhiteURL != "" {
			white[item.Plat] = append(white[item.Plat], &SectionURL{ID: item.ID, URL: item.WhiteURL})
		}
		if item.Red != "" {
			redDot[item.Plat] = append(redDot[item.Plat], &SectionURL{ID: item.ID, URL: item.Red})
		}
	}
	s.sectionCache = ss
	s.white = white
	s.redDot = redDot
	return
}

func (s *Service) convertPlat(oldPlat int8) (newPlat int8) {
	newPlat = oldPlat
	if oldPlat == 9 {
		newPlat = model.PlatAndroidB
	} else if oldPlat == 10 {
		newPlat = model.PlatIPhoneB
	}
	return newPlat
}
