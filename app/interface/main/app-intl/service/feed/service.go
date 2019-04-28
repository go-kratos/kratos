package feed

import (
	"time"

	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-card/model/card/rank"
	"go-common/app/interface/main/app-intl/conf"
	accdao "go-common/app/interface/main/app-intl/dao/account"
	arcdao "go-common/app/interface/main/app-intl/dao/archive"
	blkdao "go-common/app/interface/main/app-intl/dao/black"
	carddao "go-common/app/interface/main/app-intl/dao/card"
	locdao "go-common/app/interface/main/app-intl/dao/location"
	rankdao "go-common/app/interface/main/app-intl/dao/rank"
	rcmdao "go-common/app/interface/main/app-intl/dao/recommend"
	reldao "go-common/app/interface/main/app-intl/dao/relation"
	tagdao "go-common/app/interface/main/app-intl/dao/tag"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

// Service is show service.
type Service struct {
	c     *conf.Config
	pHit  *prom.Prom
	pMiss *prom.Prom
	// dao
	rcmd *rcmdao.Dao
	tg   *tagdao.Dao
	blk  *blkdao.Dao
	rank *rankdao.Dao
	card *carddao.Dao
	// rpc
	arc *arcdao.Dao
	acc *accdao.Dao
	rel *reldao.Dao
	loc *locdao.Dao
	// tick
	tick time.Duration
	// black cache
	blackCache map[int64]struct{} // black aids
	// ai cache
	rcmdCache []*ai.Item
	// rank cache
	rankCache []*rank.Rank
	// follow cache
	followCache map[int64]*operate.Follow
	// group cache
	groupCache map[int64]int
	// cache
	cacheCh chan func()
	// infoc
	logCh chan interface{}
}

// New new a show service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:     c,
		pHit:  prom.CacheHit,
		pMiss: prom.CacheMiss,
		// dao
		rcmd: rcmdao.New(c),
		blk:  blkdao.New(c),
		rank: rankdao.New(c),
		tg:   tagdao.New(c),
		card: carddao.New(c),
		// rpc
		arc: arcdao.New(c),
		rel: reldao.New(c),
		acc: accdao.New(c),
		loc: locdao.New(c),
		// tick
		tick: time.Duration(c.Tick),
		// group cache
		groupCache: map[int64]int{},
		// cache
		cacheCh: make(chan func(), 1024),
		// infoc
		logCh: make(chan interface{}, 1024),
	}
	s.loadBlackCache()
	s.loadRcmdCache()
	s.loadRankCache()
	s.loadUpCardCache()
	s.loadGroupCache()
	go s.cacheproc()
	go s.blackproc()
	go s.rcmdproc()
	go s.rankproc()
	go s.upCardproc()
	go s.groupproc()
	go s.infocproc()
	return
}

// addCache is.
func (s *Service) addCache(f func()) {
	select {
	case s.cacheCh <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc is.
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
