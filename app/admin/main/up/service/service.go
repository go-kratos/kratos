package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/admin/main/up/conf"
	"go-common/app/admin/main/up/dao/global"
	"go-common/app/admin/main/up/dao/manager"
	"go-common/app/admin/main/up/service/upcrmservice"
	"go-common/app/admin/main/up/util"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

// Service is service.
type Service struct {
	c   *conf.Config
	mng *manager.Dao
	// databus sub
	//upSub *databus.Databus
	// wait group
	wg sync.WaitGroup
	// chan for mids
	//midsChan chan map[int64]int
	//for cache func
	missch chan func()
	//prom
	pCacheHit  *prom.Prom
	pCacheMiss *prom.Prom
	tokenMX    sync.Mutex
	changeMX   sync.Mutex
	tokenChan  chan int
	httpClient *bm.Client
	Crmservice *upcrmservice.Service
}

// New is go-common/business/service/videoup service implementation.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		mng: manager.New(c),
		//upSub:          databus.New(c.UpSub.Config),
		//midsChan:       make(chan map[int64]int, c.ChanSize),
		//upChan:         make(chan *model.Msg, c.UpSub.UpChanSize),
		//consumeToken:   time.Now().UnixNano(),
		//consumeRate:    int64(1e9 / c.UpSub.ConsumeLimit),
		tokenMX:    sync.Mutex{},
		missch:     make(chan func(), 1024),
		pCacheHit:  prom.CacheHit,
		pCacheMiss: prom.CacheMiss,
		changeMX:   sync.Mutex{},
		tokenChan:  make(chan int, 10),
		httpClient: bm.NewClient(c.HTTPClient.Normal),
		Crmservice: upcrmservice.New(c),
	}
	s.mng.HTTPClient = s.httpClient
	s.Crmservice.SetHTTPClient(s.httpClient)
	global.Init(c)
	// load for first time
	s.refreshCache()
	go s.cacheproc()
	go s.timerproc()
	return s
}

func (s *Service) refreshCache() {
	log.Info("refresh cache")
}

func (s *Service) cacheproc() {
	for {
		time.Sleep(5 * time.Minute)
		s.refreshCache()
	}
}

// Ping service
func (s *Service) Ping(c context.Context) (err error) {
	return
}

// Close sub.
func (s *Service) Close() {
	//close(s.midsChan)
	s.wg.Wait()
	s.mng.Close()
}

// AddCache add to chan for cache

func (s *Service) timerproc() {
	var t = time.NewTicker(time.Second)
	for now := range t.C {
		util.GlobalTimer.Advance(now)
	}
}
