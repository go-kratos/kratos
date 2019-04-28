package region

import (
	"time"

	"go-common/app/interface/main/app-show/conf"
	actdao "go-common/app/interface/main/app-show/dao/activity"
	addao "go-common/app/interface/main/app-show/dao/ad"
	arcdao "go-common/app/interface/main/app-show/dao/archive"
	adtdao "go-common/app/interface/main/app-show/dao/audit"
	bgmdao "go-common/app/interface/main/app-show/dao/bangumi"
	carddao "go-common/app/interface/main/app-show/dao/card"
	dyndao "go-common/app/interface/main/app-show/dao/dynamic"
	locdao "go-common/app/interface/main/app-show/dao/location"
	rcmmndao "go-common/app/interface/main/app-show/dao/recommend"
	rgdao "go-common/app/interface/main/app-show/dao/region"
	resdao "go-common/app/interface/main/app-show/dao/resource"
	searchdao "go-common/app/interface/main/app-show/dao/search"
	tagdao "go-common/app/interface/main/app-show/dao/tag"
	"go-common/app/interface/main/app-show/model/bangumi"
	"go-common/app/interface/main/app-show/model/card"
	"go-common/app/interface/main/app-show/model/region"
	resource "go-common/app/service/main/resource/model"
	"go-common/library/stat/prom"
)

const (
	_initRegionTagKey = "region_tag_%d_%d"
)

// Service is region service.
type Service struct {
	c *conf.Config
	// prom
	pHit   *prom.Prom
	pMiss  *prom.Prom
	prmobi *prom.Prom
	// dao
	dao *rgdao.Dao
	// bnnr   *bnnrdao.Dao
	rcmmnd *rcmmndao.Dao
	ad     *addao.Dao // cptbanner
	tag    *tagdao.Dao
	adt    *adtdao.Dao
	arc    *arcdao.Dao
	dyn    *dyndao.Dao
	search *searchdao.Dao
	cdao   *carddao.Dao
	act    *actdao.Dao
	bgm    *bgmdao.Dao
	res    *resdao.Dao
	loc    *locdao.Dao
	// tick
	tick time.Duration
	// regions cache
	cache map[string][]*region.Region
	// new region list cache
	cachelist       map[string][]*region.Region
	limitCache      map[int64][]*region.Limit
	configCache     map[int64][]*region.Config
	regionListCache map[string]map[int]*region.Region
	verCache        map[string]string
	// audit cache
	auditCache map[string]map[int]struct{} // audit mobi_app builds
	// region show item cache
	bannerCache    map[int8]map[int][]*resource.Banner
	bannerBmgCache map[int8]map[int][]*bangumi.Banner
	hotCache       map[int][]*region.ShowItem
	newCache       map[int][]*region.ShowItem
	dynamicCache   map[int][]*region.ShowItem
	// overseas
	hotOseaCache     map[int][]*region.ShowItem
	newOseaCache     map[int][]*region.ShowItem
	dynamicOseaCache map[int][]*region.ShowItem
	// region child show item cache
	childHotCache        map[int][]*region.ShowItem
	childNewCache        map[int][]*region.ShowItem
	childHotAidsCache    map[int][]int64
	childNewAidsCache    map[int][]int64
	showDynamicAidsCache map[int][]int64
	// overseas region child show item cache
	childHotOseaCache map[int][]*region.ShowItem
	childNewOseaCache map[int][]*region.ShowItem
	// region tag show item cache
	tagHotCache     map[string][]*region.ShowItem
	tagNewCache     map[string][]*region.ShowItem
	tagHotAidsCache map[string][]int64
	tagNewAidsCache map[string][]int64
	// overseas region tag show item cache
	tagHotOseaCache map[string][]*region.ShowItem
	tagNewOseaCache map[string][]*region.ShowItem
	// new region feed
	regionFeedCache     map[int]*region.Show
	regionFeedOseaCache map[int]*region.Show
	// tags cache
	tagsCache map[string]string
	// region show
	showCache      map[int]*region.Show
	childShowCache map[int]*region.Show
	// overseas region show
	showOseaCache      map[int]*region.Show
	childShowOseaCache map[int]*region.Show
	// region dynamic show
	showDynamicCache      map[int]*region.Show
	childShowDynamicCache map[int]*region.Show
	// overseas region dynamic show
	showDynamicOseaCache      map[int]*region.Show
	childShowDynamicOseaCache map[int]*region.Show
	// similar tag
	similarTagCache map[string][]*region.SimilarTag
	// similar tag
	regionTagCache map[int][]*region.SimilarTag
	// ranking
	rankCache     map[int][]*region.ShowItem
	rankOseaCache map[int][]*region.ShowItem
	// card
	cardCache       map[string][]*region.Head
	columnListCache map[int]*card.ColumnList
	// region
	reRegionCache map[int]*region.Region
	// json tick
	jsonOn       bool
	jsonCh       chan int64
	jsonIdsCache map[int64]struct{} // rid<<32 | tid
	// cpm percentage   0~100
	cpmNum   int
	cpmMid   map[int64]struct{}
	cpmAll   bool
	adIsPost bool
	// infoc
	logCh chan interface{}
}

// New new a region service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		pHit:   prom.CacheHit,
		pMiss:  prom.CacheMiss,
		prmobi: prom.BusinessInfoCount,
		// dao
		dao: rgdao.New(c),
		// bnnr:    bnnrdao.New(c),
		rcmmnd: rcmmndao.New(c),
		ad:     addao.New(c),
		adt:    adtdao.New(c),
		arc:    arcdao.New(c),
		tag:    tagdao.New(c),
		dyn:    dyndao.New(c),
		search: searchdao.New(c),
		cdao:   carddao.New(c),
		act:    actdao.New(c),
		bgm:    bgmdao.New(c),
		res:    resdao.New(c),
		loc:    locdao.New(c),
		// tick
		tick: time.Duration(c.Tick),
		// audit cache
		auditCache: map[string]map[int]struct{}{},
		// regions cache
		cache: map[string][]*region.Region{},
		// new region list cache
		cachelist:       map[string][]*region.Region{},
		limitCache:      map[int64][]*region.Limit{},
		configCache:     map[int64][]*region.Config{},
		regionListCache: map[string]map[int]*region.Region{},
		verCache:        map[string]string{},
		// region show item cache
		bannerCache:    map[int8]map[int][]*resource.Banner{},
		bannerBmgCache: map[int8]map[int][]*bangumi.Banner{},
		hotCache:       map[int][]*region.ShowItem{},
		newCache:       map[int][]*region.ShowItem{},
		dynamicCache:   map[int][]*region.ShowItem{},
		// overseas
		hotOseaCache:     map[int][]*region.ShowItem{},
		newOseaCache:     map[int][]*region.ShowItem{},
		dynamicOseaCache: map[int][]*region.ShowItem{},
		// region child show item cache
		childHotCache:        map[int][]*region.ShowItem{},
		childNewCache:        map[int][]*region.ShowItem{},
		childHotAidsCache:    map[int][]int64{},
		childNewAidsCache:    map[int][]int64{},
		showDynamicAidsCache: map[int][]int64{},
		// overseas region child show item cache
		childHotOseaCache: map[int][]*region.ShowItem{},
		childNewOseaCache: map[int][]*region.ShowItem{},
		// region tag show item cache
		tagHotCache:     map[string][]*region.ShowItem{},
		tagNewCache:     map[string][]*region.ShowItem{},
		tagHotAidsCache: map[string][]int64{},
		tagNewAidsCache: map[string][]int64{},
		// overseas region tag show item cache
		tagHotOseaCache: map[string][]*region.ShowItem{},
		tagNewOseaCache: map[string][]*region.ShowItem{},
		// new region feed
		regionFeedCache:     map[int]*region.Show{},
		regionFeedOseaCache: map[int]*region.Show{},
		// tags cache
		tagsCache: map[string]string{},
		// region show
		showCache:      map[int]*region.Show{},
		childShowCache: map[int]*region.Show{},
		// overseas region show
		showOseaCache:      map[int]*region.Show{},
		childShowOseaCache: map[int]*region.Show{},
		// region dynamic show
		showDynamicCache:      map[int]*region.Show{},
		childShowDynamicCache: map[int]*region.Show{},
		// overseas region dynamic show
		showDynamicOseaCache:      map[int]*region.Show{},
		childShowDynamicOseaCache: map[int]*region.Show{},
		// similar tag
		similarTagCache: map[string][]*region.SimilarTag{},
		// similar tag
		regionTagCache: map[int][]*region.SimilarTag{},
		// ranking
		rankCache:     map[int][]*region.ShowItem{},
		rankOseaCache: map[int][]*region.ShowItem{},
		// card
		cardCache:       map[string][]*region.Head{},
		columnListCache: map[int]*card.ColumnList{},
		// region
		reRegionCache: map[int]*region.Region{},
		// json tick
		jsonOn:       false,
		jsonCh:       make(chan int64, 128),
		jsonIdsCache: map[int64]struct{}{},
		// cpm percentage   0~100
		cpmNum:   0,
		cpmMid:   map[int64]struct{}{},
		cpmAll:   true,
		adIsPost: false,
		// infoc
		logCh: make(chan interface{}, 1024),
	}
	now := time.Now()
	s.loadRegionlist()
	s.loadRegion()
	s.loadShow()
	s.loadShowChild()
	s.loadShowChildTagsInfo()
	s.loadBanner()
	s.loadbgmBanner()
	s.loadAuditCache()
	s.loadRegionListCache()
	s.loadRankRegionCache()
	s.loadColumnListCache(now)
	s.loadCardCache(now)
	go s.loadproc()
	go s.infocproc()
	return
}

func (s *Service) loadproc() {
	for {
		time.Sleep(s.tick)
		now := time.Now()
		s.loadRegionlist()
		s.loadRegion()
		s.loadShow()
		s.loadShowChild()
		s.loadShowChildTagsInfo()
		s.loadBanner()
		s.loadbgmBanner()
		s.loadAuditCache()
		s.loadRegionListCache()
		s.loadRankRegionCache()
		s.loadColumnListCache(now)
		s.loadCardCache(now)
	}
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
}

// SetCpm percentage  0~100
func (s *Service) SetCpmNum(num int) {
	s.cpmNum = num
	if s.cpmNum < 0 {
		s.cpmNum = 0
	} else if s.cpmNum > 100 {
		s.cpmNum = 100
	}
}

// GetCpm percentage
func (s *Service) CpmNum() int {
	return s.cpmNum
}

// SetCpm percentage  0~100
func (s *Service) SetCpmMid(mid int64) {
	var mids = map[int64]struct{}{}
	mids[mid] = struct{}{}
	for mid, _ := range s.cpmMid {
		if _, ok := mids[mid]; !ok {
			mids[mid] = struct{}{}
		}
	}
	s.cpmMid = mids
}

// GetCpm percentage
func (s *Service) CpmMid() []int {
	var mids []int
	for mid, _ := range s.cpmMid {
		mids = append(mids, int(mid))
	}
	return mids
}

// SetCpm All
func (s *Service) SetCpmAll(isAll bool) {
	s.cpmAll = isAll
}

// GetCpm All
func (s *Service) CpmAll() int {
	if s.cpmAll {
		return 1
	}
	return 0
}

// SetIsPost Get or Post
func (s *Service) SetAdIsPost(isPost bool) {
	s.adIsPost = isPost
}

// IsPost Get or Post
func (s *Service) AdIsPost() int {
	if s.adIsPost {
		return 1
	}
	return 0
}
