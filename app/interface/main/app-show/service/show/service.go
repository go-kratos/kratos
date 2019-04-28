package show

import (
	"strconv"
	"time"

	clive "go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-card/model/card/rank"
	"go-common/app/interface/main/app-show/conf"
	accdao "go-common/app/interface/main/app-show/dao/account"
	actdao "go-common/app/interface/main/app-show/dao/activity"
	addao "go-common/app/interface/main/app-show/dao/ad"
	arcdao "go-common/app/interface/main/app-show/dao/archive"
	adtdao "go-common/app/interface/main/app-show/dao/audit"
	bgmdao "go-common/app/interface/main/app-show/dao/bangumi"
	carddao "go-common/app/interface/main/app-show/dao/card"
	dbusdao "go-common/app/interface/main/app-show/dao/databus"
	dyndao "go-common/app/interface/main/app-show/dao/dynamic"
	livedao "go-common/app/interface/main/app-show/dao/live"
	locdao "go-common/app/interface/main/app-show/dao/location"
	rcmmdao "go-common/app/interface/main/app-show/dao/recommend"
	regiondao "go-common/app/interface/main/app-show/dao/region"
	reldao "go-common/app/interface/main/app-show/dao/relation"
	resdao "go-common/app/interface/main/app-show/dao/resource"
	showdao "go-common/app/interface/main/app-show/dao/show"
	tagdao "go-common/app/interface/main/app-show/dao/tag"
	"go-common/app/interface/main/app-show/model/card"
	recmod "go-common/app/interface/main/app-show/model/recommend"
	"go-common/app/interface/main/app-show/model/region"
	"go-common/app/interface/main/app-show/model/show"
	creativeAPI "go-common/app/interface/main/creative/api"
	"go-common/app/service/main/archive/api"
	resource "go-common/app/service/main/resource/model"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

type recommend struct {
	key  string
	aids []int64
}

type rcmmndCfg struct {
	Aid   int64  `json:"aid"`
	Goto  string `json:"goto"`
	Title string `json:"title"`
	Cover string `json:"cover"`
}

// Service is show service.
type Service struct {
	c              *conf.Config
	creativeClient creativeAPI.CreativeClient
	dao            *showdao.Dao
	rcmmnd         *rcmmdao.Dao
	ad             *addao.Dao // cptbanner
	bgm            *bgmdao.Dao
	lv             *livedao.Dao
	// bnnr   *bnnrdao.Dao
	adt  *adtdao.Dao
	tag  *tagdao.Dao
	arc  *arcdao.Dao
	dbus *dbusdao.Dao
	dyn  *dyndao.Dao
	res  *resdao.Dao
	// artic   *articledao.Dao
	client *httpx.Client
	rg     *regiondao.Dao
	cdao   *carddao.Dao
	act    *actdao.Dao
	acc    *accdao.Dao
	// relation
	reldao *reldao.Dao
	loc    *locdao.Dao

	tick time.Duration

	rcmmndCache         []*show.Item
	rcmmndOseaCache     []*show.Item
	regionCache         map[string][]*show.Item
	regionOseaCache     map[string][]*show.Item
	regionBgCache       map[string][]*show.Item
	regionBgOseaCache   map[string][]*show.Item
	regionBgEpCache     map[string][]*show.Item
	regionBgEpOseaCache map[string][]*show.Item
	bgmCache            map[int8][]*show.Item
	liveCount           int
	liveMoeCache        []*show.Item // TODO change to liveMoeCache
	liveHotCache        []*show.Item // TODO change to liveHotCache
	bannerCache         map[int8]map[int][]*resource.Banner
	cache               map[string][]*show.Show
	cacheBg             map[string][]*show.Show
	cacheBgEp           map[string][]*show.Show
	tempCache           map[string][]*show.Show
	auditCache          map[string]map[int]struct{} // audit mobi_app builds
	blackCache          map[int64]struct{}          // black aids

	logCh     chan infoc
	logFeedCh chan interface{}
	rcmmndCh  chan recommend
	logPath   string

	// loadfile
	jsonOn bool
	// cpm percentage   0~100
	cpmNum       int
	cpmMid       map[int64]struct{}
	cpmAll       bool
	cpmRcmmndNum int
	cpmRcmmndMid map[int64]struct{}
	cpmRcmmndAll bool
	adIsPost     bool
	// recommend api
	rcmmndOn    bool
	rcmmndGroup map[int64]int    // mid -> group
	rcmmndHosts map[int][]string // group -> hosts
	// region
	reRegionCache map[int]*region.Region
	// ranking
	rankCache     []*show.Item
	rankOseaCache []*show.Item
	// card
	cardCache       map[string][]*show.Show
	columnListCache map[int]*card.ColumnList
	cardSetCache    map[int64]*operate.CardSet
	eventTopicCache map[int64]*operate.EventTopic
	// hot card
	hotTenTabCardCache map[int][]*recmod.CardList
	rankAidsCache      []int64
	rankScoreCache     map[int64]int64
	rankArchivesCache  map[int64]*api.Arc
	// hotCache           []*card.PopularCard
	rcmdCache       []*card.PopularCard
	hottopicsCache  []*clive.TopicHot
	rankCache2      []*rank.Rank
	dynamicHotCache []*clive.DynamicHot
	// prom
	pHit  *prom.Prom
	pMiss *prom.Prom
}

// New new a show service.
func New(c *conf.Config) (s *Service) {
	rcmmndHosts := make(map[int][]string, len(c.Recommend.Host))
	for k, v := range c.Recommend.Host {
		key, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		rcmmndHosts[key] = v
	}
	rcmmndGroup := make(map[int64]int, len(c.Recommend.Group))
	for k, v := range c.Recommend.Group {
		key, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		rcmmndGroup[int64(key)] = v
	}
	s = &Service{
		c:      c,
		dao:    showdao.New(c),
		rcmmnd: rcmmdao.New(c),
		ad:     addao.New(c),
		bgm:    bgmdao.New(c),
		lv:     livedao.New(c),
		// bnnr:   bnnrdao.New(c),
		adt:  adtdao.New(c),
		tag:  tagdao.New(c),
		arc:  arcdao.New(c),
		dbus: dbusdao.New(c),
		dyn:  dyndao.New(c),
		res:  resdao.New(c),
		// artic:   articledao.New(c),
		rg:   regiondao.New(c),
		cdao: carddao.New(c),
		act:  actdao.New(c),
		acc:  accdao.New(c),
		// relation
		reldao: reldao.New(c),
		loc:    locdao.New(c),
		client: httpx.NewClient(c.HTTPData),

		tick: time.Duration(c.Tick),

		jsonOn: false,

		logCh:     make(chan infoc, 1024),
		logFeedCh: make(chan interface{}, 1024),
		rcmmndCh:  make(chan recommend, 1024),
		logPath:   c.ShowLog,

		rcmmndOn:    false,
		rcmmndGroup: rcmmndGroup,
		rcmmndHosts: rcmmndHosts,
		// cpm percentage   0~100
		cpmNum:       0,
		cpmMid:       map[int64]struct{}{},
		cpmAll:       true,
		cpmRcmmndNum: 0,
		cpmRcmmndMid: map[int64]struct{}{},
		cpmRcmmndAll: true,
		adIsPost:     false,
		// region
		reRegionCache: map[int]*region.Region{},
		// ranking
		rankCache:     []*show.Item{},
		rankOseaCache: []*show.Item{},
		//card
		cardCache:       map[string][]*show.Show{},
		columnListCache: map[int]*card.ColumnList{},
		cardSetCache:    map[int64]*operate.CardSet{},
		eventTopicCache: map[int64]*operate.EventTopic{},
		// hot card
		hotTenTabCardCache: make(map[int][]*recmod.CardList),
		rankAidsCache:      []int64{},
		rankScoreCache:     map[int64]int64{},
		rankArchivesCache:  map[int64]*api.Arc{},
		// hotCache:           []*card.PopularCard{},
		rcmdCache:       []*card.PopularCard{},
		hottopicsCache:  []*clive.TopicHot{},
		rankCache2:      []*rank.Rank{},
		dynamicHotCache: []*clive.DynamicHot{},
		// prom
		pHit:  prom.CacheHit,
		pMiss: prom.CacheMiss,
	}
	var err error
	if s.creativeClient, err = creativeAPI.NewClient(nil); err != nil {
		panic("creativeGRPC not found!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	}
	now := time.Now()
	s.loadCache(now)
	go s.cacheproc()
	go s.infocproc()
	go s.rcmmndproc()
	go s.infocfeedproc()
	return s
}

// cacheproc load all cache.
func (s *Service) cacheproc() {
	for {
		time.Sleep(s.tick)
		now := time.Now()
		s.loadCache(now)
	}
}

func (s *Service) loadCache(now time.Time) {
	s.loadRcmmndCache(now)
	s.loadRegionCache(now)
	s.loadBgmCache(now)
	s.loadLiveCache(now)
	s.loadBannerCahce()
	s.loadShowCache()
	s.loadShowTempCache()
	s.loadBlackCache()
	s.loadAuditCache()
	s.loadRegionListCache()
	s.loadRankAllCache()
	s.loadColumnListCache(now)
	s.loadCardCache(now)
	s.loadCardSetCache()
	// hot
	// s.loadPopularCard(now)
	s.loadHotTenTabAids()
	s.loadHotTopicCache()
	for i := 0; i < 10; i++ {
		s.loadHotTenMergeRcmdCache(i)
	}
	s.loadDynamicHotCache()
	s.loadEventTopicCache()
}

// SetRcmmndOn
func (s *Service) SetRcmmndOn(on bool) {
	s.rcmmndOn = on
}

// GetRcmmndOn
func (s *Service) RcmmndOn() bool {
	return s.rcmmndOn
}

// Close dao
func (s *Service) Close() error {
	return s.dao.Close()
}

// SetRcmmndGroup set rcmmnd group data.
func (s *Service) SetRcmmndGroup(m int64, g int) {
	tmp := map[int64]int{}
	tmp[m] = g
	for k, v := range s.rcmmndGroup {
		if k != m {
			tmp[k] = v
		}
	}
	s.rcmmndGroup = tmp
}

// GetRcmmndGroup get rcmmnd group data.
func (s *Service) GetRcmmndGroup() map[string]int {
	tmp := map[string]int{}
	for k, v := range s.rcmmndGroup {
		tmp[strconv.FormatInt(k, 10)] = v
	}
	return tmp
}

// SetRcmmndHost set rcmmnd host data.
func (s *Service) SetRcmmndHost(g int, hosts []string) {
	tmp := map[int][]string{}
	tmp[g] = hosts
	for k, v := range s.rcmmndHosts {
		if k != g {
			tmp[k] = v
		}
	}
	s.rcmmndHosts = tmp
}

// GetRcmmndHost get rcmmnd host data.
func (s *Service) GetRcmmndHost() map[string][]string {
	tmp := map[string][]string{}
	for k, v := range s.rcmmndHosts {
		tmp[strconv.Itoa(k)] = v
	}
	return tmp
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

// RcmmndNum percentage
func (s *Service) RcmmndNum() int {
	return s.cpmRcmmndNum
}

// SetRcmmndNum percentage  0~100
func (s *Service) SetRcmmndNum(num int) {
	s.cpmRcmmndNum = num
	if s.cpmRcmmndNum < 0 {
		s.cpmRcmmndNum = 0
	} else if s.cpmRcmmndNum > 100 {
		s.cpmRcmmndNum = 100
	}
}

// CpmRcmmndMid Mid
func (s *Service) CpmRcmmndMid() []int {
	var mids []int
	for mid, _ := range s.cpmRcmmndMid {
		mids = append(mids, int(mid))
	}
	return mids
}

// SetCpmRcmmndMid Mid
func (s *Service) SetCpmRcmmndMid(mid int64) {
	var mids = map[int64]struct{}{}
	mids[mid] = struct{}{}
	for mid, _ := range s.cpmRcmmndMid {
		if _, ok := mids[mid]; !ok {
			mids[mid] = struct{}{}
		}
	}
	s.cpmRcmmndMid = mids
}

// CpmRcmmnd All
func (s *Service) CpmRcmmndAll() int {
	if s.cpmRcmmndAll {
		return 1
	}
	return 0
}

// SetCpmRcmmnd All
func (s *Service) SetCpmRcmmndAll(isAll bool) {
	s.cpmRcmmndAll = isAll
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
