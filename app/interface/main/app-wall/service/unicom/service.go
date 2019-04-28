package unicom

import (
	"sync"
	"time"

	"go-common/app/interface/main/app-wall/conf"
	accdao "go-common/app/interface/main/app-wall/dao/account"
	liveDao "go-common/app/interface/main/app-wall/dao/live"
	seqDao "go-common/app/interface/main/app-wall/dao/seq"
	shopDao "go-common/app/interface/main/app-wall/dao/shopping"
	unicomDao "go-common/app/interface/main/app-wall/dao/unicom"
	"go-common/app/interface/main/app-wall/model/unicom"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"
)

const (
	_initIPlimitKey = "iplimit_%v_%v"
)

type Service struct {
	c                *conf.Config
	dao              *unicomDao.Dao
	live             *liveDao.Dao
	seqdao           *seqDao.Dao
	accd             *accdao.Dao
	shop             *shopDao.Dao
	tick             time.Duration
	unicomIpCache    []*unicom.UnicomIP
	unicomIpSQLCache map[string]*unicom.UnicomIP
	operationIPlimit map[string]struct{}
	unicomPackCache  []*unicom.UserPack
	// infoc
	logCh      chan interface{}
	packCh     chan interface{}
	packLogCh  chan interface{}
	userBindCh chan interface{}
	// waiter
	waiter sync.WaitGroup
	// databus
	userbindPub *databus.Databus
	// prom
	pHit  *prom.Prom
	pMiss *prom.Prom
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                c,
		dao:              unicomDao.New(c),
		live:             liveDao.New(c),
		seqdao:           seqDao.New(c),
		accd:             accdao.New(c),
		shop:             shopDao.New(c),
		tick:             time.Duration(c.Tick),
		unicomIpCache:    []*unicom.UnicomIP{},
		unicomIpSQLCache: map[string]*unicom.UnicomIP{},
		operationIPlimit: map[string]struct{}{},
		unicomPackCache:  []*unicom.UserPack{},
		// databus
		userbindPub: databus.New(c.UnicomDatabus),
		// infoc
		logCh:      make(chan interface{}, 1024),
		packCh:     make(chan interface{}, 1024),
		packLogCh:  make(chan interface{}, 1024),
		userBindCh: make(chan interface{}, 1024),
		// prom
		pHit:  prom.CacheHit,
		pMiss: prom.CacheMiss,
	}
	// now := time.Now()
	s.loadIPlimit(c)
	s.loadUnicomIP()
	// s.loadUnicomIPOrder(now)
	s.loadUnicomPacks()
	// s.loadUnicomFlow()
	go s.loadproc()
	go s.unicomInfocproc()
	go s.unicomPackInfocproc()
	go s.addUserPackLogproc()
	s.waiter.Add(1)
	go s.userbindConsumer()
	return
}

// cacheproc load cache
func (s *Service) loadproc() {
	for {
		time.Sleep(s.tick)
		// now := time.Now()
		s.loadUnicomIP()
		// s.loadUnicomIPOrder(now)
		s.loadUnicomPacks()
		// s.loadUnicomFlow()
	}
}
