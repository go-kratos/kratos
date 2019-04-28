package service

import (
	"context"

	"go-common/app/interface/main/tv/conf"
	appDao "go-common/app/interface/main/tv/dao/app"
	arcDao "go-common/app/interface/main/tv/dao/archive"
	auditDao "go-common/app/interface/main/tv/dao/audit"
	cmsDao "go-common/app/interface/main/tv/dao/cms"
	"go-common/app/interface/main/tv/dao/search"
	"go-common/app/interface/main/tv/model"
	seaMdl "go-common/app/interface/main/tv/model/search"
)

var ctx = context.Background()

// Service .
type Service struct {
	// dao
	dao       *appDao.Dao
	cmsDao    *cmsDao.Dao
	auditDao  *auditDao.Dao
	searchDao *search.Dao
	arcDao    *arcDao.Dao
	// cfg
	conf      *conf.Config
	TVAppInfo *conf.TVApp // tv app basic info
	// memory
	HomeData     *model.Homepage                  // homepage data
	ZoneData     map[int][]*model.Card            // zone list data for homepage
	RankData     map[int]map[string][]*model.Card // zone pages data
	HeaderSids   map[int]int                      // use to remove duplicated ones
	ZoneSids     map[int]map[int]int              // same use as HeaderSids
	ZonesInfo    map[int]*conf.PageCfg            // zones information data
	PGCOrigins   map[int][]*model.Card            // pgc zone list data
	UGCOrigins   map[int][]*model.Card            // pgc types data
	PGCIndexShow map[int64]string                 // pgc index show data
	ModPages     map[int][]*model.Module          // module pages
	IdxIntervs   *seaMdl.IdxIntervSave            // index intervention storage
	RegionInfo   []*model.Region                  // region all
	MaxTime      int64
	styleLabel   map[int][]*model.ParamStyle // style label
}

// New .
func New(c *conf.Config) *Service {
	srv := &Service{
		// dao
		dao:       appDao.New(c),
		cmsDao:    cmsDao.New(c),
		auditDao:  auditDao.New(c),
		searchDao: search.New(c),
		arcDao:    arcDao.New(c),
		// config
		conf:      c,
		TVAppInfo: c.TVApp,
		// memory data
		ZoneData:     make(map[int][]*model.Card),
		RankData:     make(map[int]map[string][]*model.Card),
		ZonesInfo:    make(map[int]*conf.PageCfg),
		PGCOrigins:   make(map[int][]*model.Card),
		ZoneSids:     make(map[int]map[int]int),
		ModPages:     make(map[int][]*model.Module),
		PGCIndexShow: make(map[int64]string),
		HeaderSids:   make(map[int]int),
		IdxIntervs: &seaMdl.IdxIntervSave{
			Pgc: make(map[int][]int64),
			Ugc: make(map[int][]int64),
		},
		RegionInfo: make([]*model.Region, 0),
		styleLabel: make(map[int][]*model.ParamStyle),
	}
	// transform string map to int map
	for k, v := range c.Newzone {
		srv.ZonesInfo[atoi(k)] = v
	}
	// not blocking data loading
	srv.indexShow() // pgc index show data
	go srv.indexShowproc()
	srv.loadRegion() // load all dynamic regions
	go srv.loadRegionproc()
	srv.loadPages() // load pages
	go srv.loadPagesproc()
	srv.ppIdxIntev(ctx) // load es index interventions
	go srv.ppIdxIntervproc()
	return srv
}
