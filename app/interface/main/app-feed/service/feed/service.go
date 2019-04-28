package feed

import (
	"time"

	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-card/model/card/rank"
	"go-common/app/interface/main/app-feed/conf"
	accdao "go-common/app/interface/main/app-feed/dao/account"
	addao "go-common/app/interface/main/app-feed/dao/ad"
	arcdao "go-common/app/interface/main/app-feed/dao/archive"
	artdao "go-common/app/interface/main/app-feed/dao/article"
	audiodao "go-common/app/interface/main/app-feed/dao/audio"
	adtdao "go-common/app/interface/main/app-feed/dao/audit"
	bgmdao "go-common/app/interface/main/app-feed/dao/bangumi"
	blkdao "go-common/app/interface/main/app-feed/dao/black"
	bplusdao "go-common/app/interface/main/app-feed/dao/bplus"
	carddao "go-common/app/interface/main/app-feed/dao/card"
	cvgdao "go-common/app/interface/main/app-feed/dao/converge"
	gamedao "go-common/app/interface/main/app-feed/dao/game"
	livdao "go-common/app/interface/main/app-feed/dao/live"
	locdao "go-common/app/interface/main/app-feed/dao/location"
	rankdao "go-common/app/interface/main/app-feed/dao/rank"
	rcmdao "go-common/app/interface/main/app-feed/dao/recommend"
	reldao "go-common/app/interface/main/app-feed/dao/relation"
	rscdao "go-common/app/interface/main/app-feed/dao/resource"
	searchdao "go-common/app/interface/main/app-feed/dao/search"
	showdao "go-common/app/interface/main/app-feed/dao/show"
	specdao "go-common/app/interface/main/app-feed/dao/special"
	tabdao "go-common/app/interface/main/app-feed/dao/tab"
	tagdao "go-common/app/interface/main/app-feed/dao/tag"
	updao "go-common/app/interface/main/app-feed/dao/upper"
	"go-common/app/interface/main/app-feed/model"
	"go-common/app/interface/main/app-feed/model/feed"
	resource "go-common/app/service/main/resource/model"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

var (
	_emptyItem = []*feed.Item{}
)

// Service is show service.
type Service struct {
	c     *conf.Config
	pHit  *prom.Prom
	pMiss *prom.Prom
	// dao
	rcmd  *rcmdao.Dao
	bgm   *bgmdao.Dao
	tg    *tagdao.Dao
	adt   *adtdao.Dao
	blk   *blkdao.Dao
	lv    *livdao.Dao
	ad    *addao.Dao
	rank  *rankdao.Dao
	gm    *gamedao.Dao
	sp    *specdao.Dao
	cvg   *cvgdao.Dao
	show  *showdao.Dao
	tab   *tabdao.Dao
	audio *audiodao.Dao
	// rpc
	arc    *arcdao.Dao
	acc    *accdao.Dao
	rel    *reldao.Dao
	upper  *updao.Dao
	art    *artdao.Dao
	rsc    *rscdao.Dao
	card   *carddao.Dao
	search *searchdao.Dao
	bplus  *bplusdao.Dao
	loc    *locdao.Dao
	// tick
	tick time.Duration
	// audit cache
	auditCache map[string]map[int]struct{} // audit mobi_app builds
	// black cache
	blackCache map[int64]struct{} // black aids
	// ai cache
	rcmdCache []*ai.Item
	// rank
	rankCache []*rank.Rank
	// converge cache
	convergeCache map[int64]*operate.Converge
	// download cache
	downloadCache map[int64]*operate.Download
	// special cache
	specialCache map[int64]*operate.Special
	// follow cache
	followCache   map[int64]*operate.Follow
	liveCardCache map[int64][]*live.Card
	// tab cache
	menuCache  []*operate.Menu
	tabCache   map[int64][]*operate.Active
	coverCache map[int64]string
	// group cache
	groupCache map[int64]int
	// cache
	cacheCh chan func()
	// infoc
	logCh chan interface{}
	// ad
	cmResourceMap map[int8]int64
	// abtest cache
	abtestCache map[string]*resource.AbTest
	// autoplay mids cache
	autoplayMidsCache map[int64]struct{}
	// follow mode list
	followModeList map[int64]struct{}
}

// New new a show service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:     c,
		pHit:  prom.CacheHit,
		pMiss: prom.CacheMiss,
		// dao
		rcmd:   rcmdao.New(c),
		bgm:    bgmdao.New(c),
		tg:     tagdao.New(c),
		adt:    adtdao.New(c),
		blk:    blkdao.New(c),
		lv:     livdao.New(c),
		ad:     addao.New(c),
		rank:   rankdao.New(c),
		gm:     gamedao.New(c),
		cvg:    cvgdao.New(c),
		sp:     specdao.New(c),
		show:   showdao.New(c),
		tab:    tabdao.New(c),
		audio:  audiodao.New(c),
		card:   carddao.New(c),
		search: searchdao.New(c),
		bplus:  bplusdao.New(c),
		// rpc
		arc:   arcdao.New(c),
		rel:   reldao.New(c),
		acc:   accdao.New(c),
		upper: updao.New(c),
		art:   artdao.New(c),
		rsc:   rscdao.New(c),
		loc:   locdao.New(c),
		// tick
		tick: time.Duration(c.Tick),
		// group cache
		groupCache: map[int64]int{},
		// cache
		cacheCh: make(chan func(), 1024),
		// infoc
		logCh: make(chan interface{}, 1024),
		// abtest cache
		abtestCache: map[string]*resource.AbTest{},
		// autoplay mids cache
		autoplayMidsCache: map[int64]struct{}{},
	}
	var err error
	if s.cmResourceMap, err = s.coverCMResource(c.Feed.CMResource); err != nil {
		panic(err)
	}
	s.loadAuditCache()
	s.loadBlackCache()
	s.loadRcmdCache()
	s.loadRankCache()
	s.loadConvergeCache()
	s.loadDownloadCache()
	s.loadSpecialCache()
	s.loadTabCache()
	s.loadGroupCache()
	s.loadUpCardCache()
	s.loadLiveCardCache()
	s.loadABTestCache()
	s.loadAutoPlayMid()
	s.loadFollowModeList()
	go s.cacheproc()
	go s.auditproc()
	go s.blackproc()
	go s.rcmdproc()
	go s.rankproc()
	go s.convergeproc()
	go s.downloadproc()
	go s.specialproc()
	go s.tabproc()
	go s.groupproc()
	go s.infocproc()
	go s.upCardproc()
	go s.liveCardproc()
	go s.loadABTestCacheProc()
	go s.followModeListproc()
	return
}

func (s *Service) coverCMResource(resource map[string]int64) (rscm map[int8]int64, err error) {
	if len(resource) == 0 {
		return
	}
	rscm = make(map[int8]int64, len(resource))
	for mobiApp, r := range resource {
		var plat int8
		if mobiApp == "iphone" {
			plat = model.PlatIPhone
		} else if mobiApp == "android" {
			plat = model.PlatAndroid
		} else if mobiApp == "ipad" {
			plat = model.PlatIPad
		}
		rscm[plat] = r
	}
	return
}

func (s *Service) addCache(f func()) {
	select {
	case s.cacheCh <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

func (s *Service) cacheproc() {
	for {
		f, ok := <-s.cacheCh
		if !ok {
			log.Warn("cache proc exit")
			return
		}
		f()
	}
}
