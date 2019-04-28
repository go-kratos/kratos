package channel

import (
	"context"
	"time"

	"go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-channel/conf"
	accdao "go-common/app/interface/main/app-channel/dao/account"
	actdao "go-common/app/interface/main/app-channel/dao/activity"
	arcdao "go-common/app/interface/main/app-channel/dao/archive"
	artdao "go-common/app/interface/main/app-channel/dao/article"
	audiodao "go-common/app/interface/main/app-channel/dao/audio"
	adtdao "go-common/app/interface/main/app-channel/dao/audit"
	bgmdao "go-common/app/interface/main/app-channel/dao/bangumi"
	carddao "go-common/app/interface/main/app-channel/dao/card"
	convergedao "go-common/app/interface/main/app-channel/dao/converge"
	gamedao "go-common/app/interface/main/app-channel/dao/game"
	livdao "go-common/app/interface/main/app-channel/dao/live"
	locdao "go-common/app/interface/main/app-channel/dao/location"
	rgdao "go-common/app/interface/main/app-channel/dao/region"
	reldao "go-common/app/interface/main/app-channel/dao/relation"
	shopdao "go-common/app/interface/main/app-channel/dao/shopping"
	specialdao "go-common/app/interface/main/app-channel/dao/special"
	tabdao "go-common/app/interface/main/app-channel/dao/tab"
	tagdao "go-common/app/interface/main/app-channel/dao/tag"
	"go-common/app/interface/main/app-channel/model/card"
	"go-common/app/interface/main/app-channel/model/channel"
	"go-common/app/interface/main/app-channel/model/tab"
)

// Service channel
type Service struct {
	c *conf.Config
	// dao
	acc   *accdao.Dao
	arc   *arcdao.Dao
	act   *actdao.Dao
	art   *artdao.Dao
	adt   *adtdao.Dao
	bgm   *bgmdao.Dao
	audio *audiodao.Dao
	rel   *reldao.Dao
	sp    *shopdao.Dao
	tg    *tagdao.Dao
	cd    *carddao.Dao
	ce    *convergedao.Dao
	g     *gamedao.Dao
	sl    *specialdao.Dao
	rg    *rgdao.Dao
	lv    *livdao.Dao
	loc   *locdao.Dao
	tab   *tabdao.Dao
	// tick
	tick time.Duration
	// cache
	cardCache         map[int64][]*card.Card
	cardPlatCache     map[string][]*card.CardPlat
	upCardCache       map[int64]*operate.Follow
	convergeCardCache map[int64]*operate.Converge
	gameDownloadCache map[int64]*operate.Download
	specialCardCache  map[int64]*operate.Special
	liveCardCache     map[int64][]*live.Card
	cardSetCache      map[int64]*operate.CardSet
	menuCache         map[int64][]*tab.Menu
	// new region list cache
	cachelist   map[string][]*channel.Region
	limitCache  map[int64][]*channel.RegionLimit
	configCache map[int64][]*channel.RegionConfig
	// audit cache
	auditCache map[string]map[int]struct{} // audit mobi_app builds
	// infoc
	logCh chan interface{}
}

// New channel
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:     c,
		arc:   arcdao.New(c),
		acc:   accdao.New(c),
		adt:   adtdao.New(c),
		art:   artdao.New(c),
		act:   actdao.New(c),
		bgm:   bgmdao.New(c),
		sp:    shopdao.New(c),
		tg:    tagdao.New(c),
		cd:    carddao.New(c),
		ce:    convergedao.New(c),
		g:     gamedao.New(c),
		sl:    specialdao.New(c),
		rg:    rgdao.New(c),
		audio: audiodao.New(c),
		lv:    livdao.New(c),
		rel:   reldao.New(c),
		loc:   locdao.New(c),
		tab:   tabdao.New(c),
		// tick
		tick: time.Duration(c.Tick),
		// cache
		cardCache:         map[int64][]*card.Card{},
		cardPlatCache:     map[string][]*card.CardPlat{},
		upCardCache:       map[int64]*operate.Follow{},
		convergeCardCache: map[int64]*operate.Converge{},
		gameDownloadCache: map[int64]*operate.Download{},
		specialCardCache:  map[int64]*operate.Special{},
		cachelist:         map[string][]*channel.Region{},
		limitCache:        map[int64][]*channel.RegionLimit{},
		configCache:       map[int64][]*channel.RegionConfig{},
		liveCardCache:     map[int64][]*live.Card{},
		cardSetCache:      map[int64]*operate.CardSet{},
		menuCache:         map[int64][]*tab.Menu{},
		// audit cache
		auditCache: map[string]map[int]struct{}{},
		// infoc
		logCh: make(chan interface{}, 1024),
	}
	s.loadCache()
	go s.loadCacheproc()
	go s.infocproc()
	return
}

func (s *Service) loadCacheproc() {
	for {
		time.Sleep(s.tick)
		s.loadCache()
	}
}

func (s *Service) loadCache() {
	now := time.Now()
	s.loadAuditCache()
	s.loadRegionlist()
	// card
	s.loadCardCache(now)
	s.loadConvergeCache()
	s.loadSpecialCache()
	s.loadLiveCardCache()
	s.loadGameDownloadCache()
	s.loadCardSetCache()
	s.loadMenusCache(now)
}

// Ping is check server ping.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.cd.PingDB(c); err != nil {
		return
	}
	return
}
