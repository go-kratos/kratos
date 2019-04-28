package service

import (
	"context"
	"sync"
	"time"

	upgrpc "go-common/app/service/main/up/api/v1"
	"go-common/app/service/main/up/conf"
	"go-common/app/service/main/up/dao/archive"
	"go-common/app/service/main/up/dao/card"
	"go-common/app/service/main/up/dao/data"
	"go-common/app/service/main/up/dao/global"
	"go-common/app/service/main/up/dao/manager"
	"go-common/app/service/main/up/dao/monitor"
	"go-common/app/service/main/up/dao/up"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"
	"go-common/library/sync/pipeline/fanout"
)

// Service is service.
type Service struct {
	c *conf.Config
	// up dao
	up   *up.Dao
	mng  *manager.Dao
	card *card.Dao
	Data *data.Dao
	arc  *archive.Dao
	// monitor
	monitor *monitor.Dao
	upMo    int64
	// wait group
	wg sync.WaitGroup
	// chan for mids
	// midsChan chan map[int64]int
	//for cache func
	missch chan func()
	// prom
	pCacheHit  *prom.Prom
	pCacheMiss *prom.Prom
	// for schdule consume databus data
	upSub               *databus.Databus
	workerQueue         []chan *databus.Message
	workerCount         int
	consumeRate         int64
	specialDBAddRate    int64
	tokenChan           chan int
	specialAddDBLimiter chan int
	// permit
	authSvr      *permit.Permit
	permCheckMap map[int64][]blademaster.HandlerFunc
	// cache worker
	cacheWorker       *fanout.Fanout
	closeSub          bool
	spGroupsCache     map[int64]*upgrpc.UpGroup
	spGroupsMidsCache map[int64][]int64
}

// New is go-common/app/service/videoup service implementation.
func New(c *conf.Config) (s *Service) {
	if c.UpSub.SpecialAddDBLimit <= 0 {
		c.UpSub.SpecialAddDBLimit = 100
	}
	global.Init(c)
	s = &Service{
		c:    c,
		up:   up.New(c),
		mng:  manager.New(c),
		card: card.New(c),
		// midsChan:       make(chan map[int64]int, c.ChanSize),
		missch:     make(chan func(), 1024),
		pCacheHit:  prom.CacheHit,
		pCacheMiss: prom.CacheMiss,
		Data:       data.New(c),
		arc:        archive.New(c),
		//up databus consume.
		upSub:               databus.New(c.UpSub.Config),
		monitor:             monitor.New(c),
		tokenChan:           make(chan int, c.UpSub.ConsumeLimit), //速率缓冲大小控制
		specialAddDBLimiter: make(chan int, c.UpSub.SpecialAddDBLimit),
		consumeRate:         int64(1e9 / c.UpSub.ConsumeLimit), //per second for ConsumeLimit,如果ConsumeLimit1=1则表示1s一个，如果ConsumeLimit1=10则表示100ms一个，以此类推.
		specialDBAddRate:    int64(1e9 / c.UpSub.SpecialAddDBLimit),
		workerQueue:         make([]chan *databus.Message, c.UpSub.RoutineLimit),
		workerCount:         c.UpSub.RoutineLimit,
		// cache worker
		cacheWorker:       global.GetWorker(),
		spGroupsCache:     make(map[int64]*upgrpc.UpGroup),
		spGroupsMidsCache: make(map[int64][]int64),
	}
	// s.mng.HTTPClient = s.up.GetHTTPClient()
	if c.Consume {
		s.wg.Add(1)
		go s.upConsumer()
		for i := 0; i < s.workerCount; i++ {
			c := make(chan *databus.Message)
			s.workerQueue[i] = c
			s.wg.Add(1)
			go s.Start(c)
		}
		go s.generateToken(time.Duration(s.consumeRate), s.tokenChan)
		go s.monitorConsume()
	}
	go s.generateToken(time.Duration(s.specialDBAddRate), s.specialAddDBLimiter)
	s.refreshCache()
	go s.cacheproc()
	return s
}

func (s *Service) generateToken(duration time.Duration, tokenChan chan int) {
	var timer = time.NewTicker(duration)
	var token = 0
	for range timer.C {
		token++
		tokenChan <- token
	}
}

func (s *Service) refreshCache() {
	log.Info("refresh cache")
	s.loadUpGroups()
	s.loadSpGroupsMids()
}

func (s *Service) cacheproc() {
	for {
		time.Sleep(5 * time.Minute)
		s.refreshCache()
	}
}

// Ping service
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.up.Ping(c); err != nil {
		log.Error("up-service s.up.Ping err(%v)", err)
	}
	if err := s.card.Ping(c); err != nil {
		log.Error("up-service s.card.Ping err(%v)", err)
	}
	return
}

func (s *Service) monitorConsume() {
	if env.DeployEnv != env.DeployEnvProd {
		return
	}
	var up int64
	tm := time.Duration(s.c.Monitor.IntervalAlarm)
	for {
		time.Sleep(tm)
		if s.upMo-up == 0 {
			s.monitor.Send(context.TODO(), s.c.Monitor.UserName, "up-service did not consume within "+tm.String()+" minute, moni url"+s.c.Monitor.Moni)
		}
		up = s.upMo
	}
}

func (s *Service) loadUpGroups() {
	ugs, err := s.mng.UpGroups(context.Background())
	if err != nil {
		log.Error("s.mng.UpGroups error(%v)", err)
		return
	}
	s.spGroupsCache = ugs
}

// Close sub.
func (s *Service) Close() {
	s.upSub.Close()
	// close(s.midsChan)
	s.closeSub = true
	time.Sleep(time.Second * 2)
	s.mng.Close()
	s.up.Close()
	s.card.Close()
	global.Close()
	s.wg.Wait()
}

//SetAuthServer set auth
func (s *Service) SetAuthServer(authSvr *permit.Permit) {
	s.authSvr = authSvr
	s.permCheckMap = map[int64][]blademaster.HandlerFunc{
		15: {
			s.authSvr.Permit("PRIORITY_SIGN_UP"),
		},
	}
}
