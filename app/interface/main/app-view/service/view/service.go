package view

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/app-view/conf"
	accdao "go-common/app/interface/main/app-view/dao/account"
	actdao "go-common/app/interface/main/app-view/dao/act"
	addao "go-common/app/interface/main/app-view/dao/ad"
	aidao "go-common/app/interface/main/app-view/dao/ai"
	arcdao "go-common/app/interface/main/app-view/dao/archive"
	assdao "go-common/app/interface/main/app-view/dao/assist"
	audiodao "go-common/app/interface/main/app-view/dao/audio"
	bandao "go-common/app/interface/main/app-view/dao/bangumi"
	coindao "go-common/app/interface/main/app-view/dao/coin"
	creativedao "go-common/app/interface/main/app-view/dao/creative"
	dmdao "go-common/app/interface/main/app-view/dao/dm"
	elcdao "go-common/app/interface/main/app-view/dao/elec"
	favdao "go-common/app/interface/main/app-view/dao/favorite"
	gamedao "go-common/app/interface/main/app-view/dao/game"
	livedao "go-common/app/interface/main/app-view/dao/live"
	locdao "go-common/app/interface/main/app-view/dao/location"
	managerdao "go-common/app/interface/main/app-view/dao/manager"
	rgndao "go-common/app/interface/main/app-view/dao/region"
	reldao "go-common/app/interface/main/app-view/dao/relation"
	rscdao "go-common/app/interface/main/app-view/dao/resource"
	searchdao "go-common/app/interface/main/app-view/dao/search"
	spdao "go-common/app/interface/main/app-view/dao/special"
	tagdao "go-common/app/interface/main/app-view/dao/tag"
	thumbupdao "go-common/app/interface/main/app-view/dao/thumbup"
	ugcpaydao "go-common/app/interface/main/app-view/dao/ugcpay"
	vipdao "go-common/app/interface/main/app-view/dao/vip"
	"go-common/app/interface/main/app-view/model"
	elecmdl "go-common/app/interface/main/app-view/model/elec"
	"go-common/app/interface/main/app-view/model/live"
	"go-common/app/interface/main/app-view/model/manager"
	"go-common/app/interface/main/app-view/model/region"
	"go-common/app/interface/main/app-view/model/special"
	"go-common/app/interface/main/app-view/model/view"
	"go-common/app/service/main/archive/model/archive"
	shareclient "go-common/app/service/main/share/api"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

var (
	_elecTypeIds = []int16{
		20, 154, 156, // dance
		31, 30, 59, 29, 28, // music
		26, 22, 126, 127, // guichu
		24, 25, 47, 27, // animae
		17, 18, 16, 65, 136, 19, 121, 171, 172, 173, // game
		37, 124, 122, 39, 96, 95, 98, // tech
		71, 137, 131, // yule
		157, 158, 159, 164, // fashion
		82, 128, // movie and tv
		138, 21, 75, 76, 161, 162, 163, 174, // life
		153, 168, // guo man
		85, 86, 182, 183, 184, // film and television
	}
)

// Service is view service
type Service struct {
	c     *conf.Config
	pHit  *prom.Prom
	pMiss *prom.Prom
	prom  *prom.Prom
	// dao
	accDao      *accdao.Dao
	arcDao      *arcdao.Dao
	tagDao      *tagdao.Dao
	favDao      *favdao.Dao
	banDao      *bandao.Dao
	elcDao      *elcdao.Dao
	rgnDao      *rgndao.Dao
	liveDao     *livedao.Dao
	assDao      *assdao.Dao
	adDao       *addao.Dao
	rscDao      *rscdao.Dao
	relDao      *reldao.Dao
	coinDao     *coindao.Dao
	audioDao    *audiodao.Dao
	actDao      *actdao.Dao
	thumbupDao  *thumbupdao.Dao
	gameDao     *gamedao.Dao
	shareClient shareclient.ShareClient
	vipDao      *vipdao.Dao
	mngDao      *managerdao.Dao
	spDao       *spdao.Dao
	dmDao       *dmdao.Dao
	aiDao       *aidao.Dao
	creativeDao *creativedao.Dao
	search      *searchdao.Dao
	ugcpayDao   *ugcpaydao.Dao
	locDao      *locdao.Dao
	// region
	tick   time.Duration
	region map[int8]map[int]*region.Region
	// elec
	allowTypeIds map[int16]struct{}
	// live cache
	liveCache map[int64]*live.Live
	// chan
	inCh     chan interface{}
	dmRegion map[int16]struct{}
	// vip active cache
	vipActiveCache map[int]string
	vipTick        time.Duration
	// mamager cache
	RelateCache  []*manager.Relate
	specialCache map[int64]*special.Card
	specialMids  map[int64]struct{}
	// player icon
	playerIcon *view.PlayerIcon
	// view relate game from AI
	RelateGameCache map[int64]int64
	// bnj caches
	BnjMainView *archive.View3
	BnjLists    []*archive.View3
	BnjElecInfo *elecmdl.Info
	BnjWhiteMid map[int64]struct{}
	BnjIsGrey   bool
}

// New new archive
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:     c,
		pHit:  prom.CacheHit,
		pMiss: prom.CacheMiss,
		prom:  prom.BusinessInfoCount,
		// dao
		accDao:      accdao.New(c),
		arcDao:      arcdao.New(c),
		tagDao:      tagdao.New(c),
		favDao:      favdao.New(c),
		banDao:      bandao.New(c),
		elcDao:      elcdao.New(c),
		rgnDao:      rgndao.New(c),
		liveDao:     livedao.New(c),
		assDao:      assdao.New(c),
		adDao:       addao.New(c),
		rscDao:      rscdao.New(c),
		relDao:      reldao.New(c),
		coinDao:     coindao.New(c),
		audioDao:    audiodao.New(c),
		actDao:      actdao.New(c),
		thumbupDao:  thumbupdao.New(c),
		gameDao:     gamedao.New(c),
		vipDao:      vipdao.New(c),
		mngDao:      managerdao.New(c),
		spDao:       spdao.New(c),
		dmDao:       dmdao.New(c),
		aiDao:       aidao.New(c),
		creativeDao: creativedao.New(c),
		search:      searchdao.New(c),
		ugcpayDao:   ugcpaydao.New(c),
		locDao:      locdao.New(c),
		// region
		tick:   time.Duration(c.Tick),
		region: map[int8]map[int]*region.Region{},
		// live cache
		liveCache: map[int64]*live.Live{},
		// chan
		inCh:         make(chan interface{}, 1024),
		allowTypeIds: map[int16]struct{}{},
		dmRegion:     map[int16]struct{}{},
		specialMids:  map[int64]struct{}{},
		// vip
		vipActiveCache: make(map[int]string),
		vipTick:        time.Duration(c.VipTick),
		// manager
		RelateCache:  []*manager.Relate{},
		specialCache: map[int64]*special.Card{},
		// player icon
		playerIcon: &view.PlayerIcon{},
	}
	for _, id := range _elecTypeIds {
		s.allowTypeIds[id] = struct{}{}
	}
	for _, id := range c.DMRegion {
		s.dmRegion[id] = struct{}{}
	}
	var err error
	if s.shareClient, err = shareclient.NewClient(nil); err != nil {
		panic(fmt.Sprintf("env:%s no share-service", env.DeployEnv))
	}
	// load data
	s.loadLive()
	s.loadRegion()
	s.loadPlayerIcon()
	s.loadVIPActive()
	s.loadManager()
	s.loadRelateGame()
	s.loadBnj2019Infos()
	go s.infocproc()
	go s.tickproc()
	go s.vipproc()
	go s.bnjTickproc()
	return s
}

// Ping is dao ping.
func (s *Service) Ping(c context.Context) (err error) {
	return s.arcDao.Ping(c)
}

func (s *Service) bnjTickproc() {
	for {
		time.Sleep(time.Duration(s.c.Bnj2019.Tick))
		err := s.loadBnj2019Infos()
		if err != nil {
			log.Error("bnj load error(%v)", err)
		}
	}
}

// tickproc tick load cache.
func (s *Service) tickproc() {
	for {
		time.Sleep(s.tick)
		s.loadRegion()
		s.loadLive()
		s.loadPlayerIcon()
		s.loadManager()
		s.loadRelateGame()
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

func (s *Service) loadRegion() {
	res, err := s.rgnDao.Seconds(context.TODO())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.region = res
}

func (s *Service) loadLive() {
	if s.c.Env == model.EnvHK {
		return
	}
	living, err := s.liveDao.Living(context.TODO())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	tmp := map[int64]*live.Live{}
	for _, v := range living {
		tmp[v.Mid] = v
	}
	s.liveCache = tmp
}

func (s *Service) loadManager() {
	r, err := s.mngDao.Relate(context.TODO())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.RelateCache = r
	sp, err := s.spDao.Card(context.TODO())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.specialCache = sp
	midsM, err := s.creativeDao.Special(context.Background())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	log.Info("load special mids(%+v)", midsM)
	s.specialMids = midsM
}

func (s *Service) loadRelateGame() {
	g, err := s.aiDao.Av2Game(context.TODO())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.RelateGameCache = g
}

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

func (s Service) relateGame(c context.Context, aid int64) (id int64) {
	id = s.RelateGameCache[aid]
	return
}
