package tag

import (
	"time"

	"go-common/app/interface/main/app-tag/conf"
	arcdao "go-common/app/interface/main/app-tag/dao/archive"
	bgmdao "go-common/app/interface/main/app-tag/dao/bangumi"
	rcmdao "go-common/app/interface/main/app-tag/dao/recommend"
	rgdao "go-common/app/interface/main/app-tag/dao/region"
	tagdao "go-common/app/interface/main/app-tag/dao/tag"
	"go-common/app/interface/main/app-tag/model/region"
	"go-common/app/interface/main/app-tag/model/tag"
	"go-common/library/stat/prom"
)

const (
	_initRegionKey    = "region_key_%d_%v"
	_initlanguage     = "hans"
	_initRegionTagKey = "region_tag_%d_%d"
	_initTagNameKey   = "tag_name_%v"
	_bangumiSeasonID  = 1
	_bangumiEpisodeID = 2
)

type Service struct {
	c *conf.Config
	// dao
	rcmd *rcmdao.Dao
	tg   *tagdao.Dao
	// rpc
	arc *arcdao.Dao
	// dao
	regiondao *rgdao.Dao
	bgm       *bgmdao.Dao
	// tick
	tick time.Duration
	// prom
	pHit   *prom.Prom
	pMiss  *prom.Prom
	prmobi *prom.Prom
	// regions cache
	regionListCache map[string]map[int]*region.Region
	// similar tag
	similarTagCache map[int64][]*tag.SimilarTag
	regionTagCache  map[int][]*tag.SimilarTag
	// regions cache
	regionCache   map[string][]*region.Region
	reRegionCache map[int]*region.Region
	// tags
	tagsCache     map[string]string
	tagsNameCache map[string]int64
	// tag change detail
	tagsDetailCache     map[int64][]*region.ShowItem
	tagsDetailOseaCache map[int64][]*region.ShowItem
	tagsDetailAidsCache map[int64][]int64
	// tag detail ranking
	tagsDetailRankingAidsCache map[string][]int64
	tagsDetailRankingCache     map[string][]*region.ShowItem
	tagsDetailRankingOseaCache map[string][]*region.ShowItem
	// infoc
	logCh chan interface{}
}

// New a region service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:    c,
		tg:   tagdao.New(c),
		rcmd: rcmdao.New(c),
		// rpc
		arc: arcdao.New(c),
		// dao
		regiondao: rgdao.New(c),
		bgm:       bgmdao.New(c),
		// tick
		tick: time.Duration(c.Tick),
		// regions cache
		regionListCache: map[string]map[int]*region.Region{},
		// prom
		pHit:   prom.CacheHit,
		pMiss:  prom.CacheMiss,
		prmobi: prom.BusinessInfoCount,
		// similar tag
		similarTagCache: map[int64][]*tag.SimilarTag{},
		regionTagCache:  map[int][]*tag.SimilarTag{},
		// regions cache
		regionCache:   map[string][]*region.Region{},
		reRegionCache: map[int]*region.Region{},
		// tags cache
		tagsCache:     map[string]string{},
		tagsNameCache: map[string]int64{},
		// tag Change Detail
		tagsDetailCache:     map[int64][]*region.ShowItem{},
		tagsDetailOseaCache: map[int64][]*region.ShowItem{},
		tagsDetailAidsCache: map[int64][]int64{},
		// tag detail ranking
		tagsDetailRankingAidsCache: map[string][]int64{},
		tagsDetailRankingCache:     map[string][]*region.ShowItem{},
		tagsDetailRankingOseaCache: map[string][]*region.ShowItem{},
		// infoc
		logCh: make(chan interface{}, 1024),
	}
	s.loadRegion()
	s.loadShowChildTagsInfo()
	go s.loadproc()
	go s.infocfeedproc()
	return
}

func (s *Service) loadproc() {
	for {
		time.Sleep(s.tick)
		s.loadRegion()
		s.loadShowChildTagsInfo()
		s.loadShowChildTags()
	}
}
