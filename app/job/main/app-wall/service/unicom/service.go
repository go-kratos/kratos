package unicom

import (
	"sync"
	"time"

	"go-common/app/job/main/app-wall/conf"
	seqDao "go-common/app/job/main/app-wall/dao/seq"
	unicomDao "go-common/app/job/main/app-wall/dao/unicom"
	"go-common/app/job/main/app-wall/model/unicom"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"
)

type Service struct {
	c        *conf.Config
	dao      *unicomDao.Dao
	seqdao   *seqDao.Dao
	clickSub *databus.Databus
	closed   bool
	// waiter
	waiter    sync.WaitGroup
	cliChan   []chan *unicom.ClickMsg
	dbcliChan []chan *unicom.UserBind
	// infoc
	logCh         []chan interface{}
	packCh        chan interface{}
	packLogCh     chan interface{}
	integralLogCh []chan interface{}
	// prom
	pHit  *prom.Prom
	pMiss *prom.Prom
	// tick
	tick      time.Duration
	lastmonth map[int]bool
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:        c,
		dao:      unicomDao.New(c),
		clickSub: databus.New(c.ReportDatabus),
		seqdao:   seqDao.New(c),
		// infoc
		packCh:    make(chan interface{}, 1024),
		packLogCh: make(chan interface{}, 1024),
		// close
		closed: false,
		// prom
		pHit:      prom.CacheHit,
		pMiss:     prom.CacheMiss,
		lastmonth: map[int]bool{},
		// tick
		tick: time.Duration(c.Tick),
	}
	for i := int64(0); i < s.c.ChanNum; i++ {
		s.cliChan = append(s.cliChan, make(chan *unicom.ClickMsg, 300000))
	}
	for i := int64(0); i < s.c.ChanDBNum; i++ {
		s.dbcliChan = append(s.dbcliChan, make(chan *unicom.UserBind, 1024))
		s.integralLogCh = append(s.integralLogCh, make(chan interface{}, 1024))
		s.logCh = append(s.logCh, make(chan interface{}, 1024))
	}
	for i := int64(0); i < s.c.ChanNum; i++ {
		s.waiter.Add(1)
		go s.cliChanProc(i)
	}
	for i := int64(0); i < s.c.ChanDBNum; i++ {
		s.waiter.Add(1)
		go s.unicomInfocproc(i)
		go s.addUserIntegralLogproc(i)
	}
	for i := int64(0); i < s.c.ChanDBNum; i++ {
		s.waiter.Add(1)
		go s.dbcliChanProc(i)
	}
	s.waiter.Add(1)
	go s.clickConsumer()
	s.waiter.Add(1)
	now := time.Now()
	if s.c.Monthly {
		// s.updatemonth(now)
		s.upBindAll()
	}
	s.waiter.Add(1)
	s.loadUnicomIPOrder(now)
	s.loadUnicomFlow()
	go s.loadproc()
	s.waiter.Add(1)
	go s.unicomPackInfocproc()
	go s.addUserPackLogproc()
	return
}

// Close Service
func (s *Service) Close() {
	s.closed = true
	time.Sleep(time.Second * 2)
	s.clickSub.Close()
	for i := 0; i < len(s.cliChan); i++ {
		close(s.cliChan[i])
	}
	for i := 0; i < len(s.dbcliChan); i++ {
		close(s.dbcliChan[i])
		close(s.integralLogCh[i])
		close(s.logCh[i])
	}
	s.waiter.Wait()
	log.Info("app-wall-job unicom flow closed.")
}

// cacheproc load cache
func (s *Service) loadproc() {
	for {
		time.Sleep(s.tick)
		now := time.Now()
		s.loadUnicomFlow()
		s.updatemonth(now)
		s.loadUnicomIPOrder(now)
	}
}
