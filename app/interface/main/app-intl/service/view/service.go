package view

import (
	"context"
	"time"

	"go-common/app/interface/main/app-intl/conf"
	"go-common/app/interface/main/app-intl/model/region"

	accdao "go-common/app/interface/main/app-intl/dao/account"
	arcdao "go-common/app/interface/main/app-intl/dao/archive"
	assdao "go-common/app/interface/main/app-intl/dao/assist"
	audiodao "go-common/app/interface/main/app-intl/dao/audio"
	bandao "go-common/app/interface/main/app-intl/dao/bangumi"
	coindao "go-common/app/interface/main/app-intl/dao/coin"

	// creativedao "go-common/app/interface/main/app-intl/dao/creative"
	dmdao "go-common/app/interface/main/app-intl/dao/dm"
	favdao "go-common/app/interface/main/app-intl/dao/favorite"
	managerdao "go-common/app/interface/main/app-intl/dao/manager"
	rgndao "go-common/app/interface/main/app-intl/dao/region"
	reldao "go-common/app/interface/main/app-intl/dao/relation"

	rscdao "go-common/app/interface/main/app-intl/dao/resource"
	tagdao "go-common/app/interface/main/app-intl/dao/tag"
	thumbupdao "go-common/app/interface/main/app-intl/dao/thumbup"

	locdao "go-common/app/interface/main/app-intl/dao/location"
	vipdao "go-common/app/interface/main/app-intl/dao/vip"
	"go-common/app/interface/main/app-intl/model"
	"go-common/app/interface/main/app-intl/model/manager"
	"go-common/app/interface/main/app-intl/model/view"

	"go-common/library/log"
	"go-common/library/stat/prom"
)

// Service is view service
type Service struct {
	c     *conf.Config
	pHit  *prom.Prom
	pMiss *prom.Prom
	prom  *prom.Prom
	// dao
	accDao     *accdao.Dao
	arcDao     *arcdao.Dao
	tagDao     *tagdao.Dao
	favDao     *favdao.Dao
	banDao     *bandao.Dao
	rgnDao     *rgndao.Dao
	assDao     *assdao.Dao
	audioDao   *audiodao.Dao
	thumbupDao *thumbupdao.Dao
	rscDao     *rscdao.Dao
	relDao     *reldao.Dao
	coinDao    *coindao.Dao
	vipDao     *vipdao.Dao
	mngDao     *managerdao.Dao
	dmDao      *dmdao.Dao
	locDao     *locdao.Dao
	// creativeDao *creativedao.Dao
	// tick
	tick time.Duration
	// region
	region map[int8]map[int16]*region.Region
	// chan
	inCh chan interface{}
	// vip active cache
	vipActiveCache map[int]string
	vipTick        time.Duration
	// mamager cache
	RelateCache []*manager.Relate
	// player icon
	playerIcon *view.PlayerIcon
	// view relate game from AI
	RelateGameCache map[int64]int64
}

// New new archive
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:     c,
		pHit:  prom.CacheHit,
		pMiss: prom.CacheMiss,
		prom:  prom.BusinessInfoCount,
		// dao
		accDao:     accdao.New(c),
		arcDao:     arcdao.New(c),
		tagDao:     tagdao.New(c),
		favDao:     favdao.New(c),
		banDao:     bandao.New(c),
		rgnDao:     rgndao.New(c),
		assDao:     assdao.New(c),
		relDao:     reldao.New(c),
		coinDao:    coindao.New(c),
		audioDao:   audiodao.New(c),
		thumbupDao: thumbupdao.New(c),
		rscDao:     rscdao.New(c),
		vipDao:     vipdao.New(c),
		mngDao:     managerdao.New(c),
		dmDao:      dmdao.New(c),
		locDao:     locdao.New(c),
		// tick
		tick: time.Duration(c.Tick),
		// region
		region: map[int8]map[int16]*region.Region{},
		// chan
		inCh: make(chan interface{}, 1024),
		// vip
		vipActiveCache: make(map[int]string),
		vipTick:        time.Duration(c.View.VipTick),
		// manager
		RelateCache: []*manager.Relate{},
		// player icon
		playerIcon: &view.PlayerIcon{},
	}
	// load data
	s.loadRegion()
	s.loadPlayerIcon()
	s.loadVIPActive()
	s.loadManager()
	go s.infocproc()
	go s.tickproc()
	go s.vipproc()
	return
}

// Ping is dao ping.
func (s *Service) Ping(c context.Context) (err error) {
	return s.arcDao.Ping(c)
}

// tickproc tick load cache.
func (s *Service) tickproc() {
	for {
		time.Sleep(s.tick)
		s.loadPlayerIcon()
		s.loadManager()
	}
}

// vipproc tick load vip cache.
func (s *Service) vipproc() {
	for {
		time.Sleep(s.vipTick)
		s.loadVIPActive()
	}
}

// loadVIPActive tick load vip active cache.
func (s *Service) loadVIPActive() {
	var (
		va  map[int]string
		err error
	)
	va = make(map[int]string)
	if va[view.VIPActiveView], err = s.vipDao.VIPActive(context.TODO(), view.VIPActiveView); err != nil {
		log.Error("s.vipDao.VIPActinve(%d) error(%v)", view.VIPActiveView, err)
		return
	}
	s.vipActiveCache = va
	log.Info("load vip active success")
}

// loadRegion is.
func (s *Service) loadRegion() {
	res, err := s.rgnDao.Seconds(context.TODO())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.region = res
}

// loadManager is.
func (s *Service) loadManager() {
	r, err := s.mngDao.Relate(context.TODO())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.RelateCache = r
}

// loadPlayerIcon is.
func (s *Service) loadPlayerIcon() {
	res, err := s.rscDao.PlayerIcon(context.TODO())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	if res != nil {
		s.playerIcon = &view.PlayerIcon{URL1: res.URL1, Hash1: res.Hash1, URL2: res.URL2, Hash2: res.Hash2, CTime: res.CTime}
	} else {
		s.playerIcon = nil
	}
}

// relateCache is.
func (s *Service) relateCache(c context.Context, plat int8, build int, now time.Time, aid int64, tids []int64, rid int32) (relate *manager.Relate) {
	rs := s.RelateCache
	rls := make([]*manager.Relate, 0, len(rs))
	if len(rs) != 0 {
	LOOP:
		for _, r := range rs {
			if vs, ok := r.Versions[plat]; ok {
				for _, v := range vs {
					if model.InvalidBuild(build, v.Build, v.Condition) {
						continue LOOP
					}
				}
				if (r.STime == 0 || now.After(r.STime.Time())) && (r.ETime == 0 || now.Before(r.ETime.Time())) {
					rls = append(rls, r)
				}
			}
		}
	}
	for _, r := range rls {
		if _, ok := r.Aids[aid]; ok {
			relate = r
			break
		}
		if len(tids) != 0 {
			for _, tid := range tids {
				if _, ok := r.Tids[tid]; ok {
					relate = r
					break
				}
			}
		}
		if _, ok := r.Rids[int64(rid)]; ok {
			relate = r
			break
		}
	}
	return
}
