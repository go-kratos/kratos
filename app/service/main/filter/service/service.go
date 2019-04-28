package service

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"go-common/app/service/main/filter/conf"
	"go-common/app/service/main/filter/dao"
	"go-common/app/service/main/filter/model/lrulist"
	"go-common/app/service/main/filter/service/ai"
	"go-common/app/service/main/filter/service/area"
	"go-common/app/service/main/filter/service/filter"
	"go-common/app/service/main/filter/service/white"
	spyrpc "go-common/app/service/main/spy/rpc/client"
	"go-common/library/cache"
	"go-common/library/log"
	"go-common/library/log/infoc"
	"go-common/library/stat/prom"

	pcre "github.com/GRbit/go-pcre"
)

// Service struct.
type Service struct {
	dao *dao.Dao

	// 业务列表
	areas *area.Area
	// 业务白名单
	whites *white.White
	// 业务黑名单
	filters *filter.Filter
	// mid白名单
	whiteMids *ai.AI

	hbaseCh *cache.Cache
	cacheCh *cache.Cache
	// lru
	lruMax  int
	lruLock sync.RWMutex
	lruList *lrulist.List

	// prom
	areaBlackHitProm *prom.Prom

	spyRPC *spyrpc.Service
	// ai chan
	aich        chan func()
	aiDelayTick time.Duration
	missch      chan func()
	infoc       *infoc.Infoc
	lastTime    int64
}

// New new service.
func New() (s *Service) {
	var err error
	s = &Service{
		dao:              dao.New(conf.Conf),
		areas:            area.New(),
		whites:           white.New(),
		filters:          filter.New(),
		whiteMids:        ai.New(),
		lruMax:           conf.Conf.Property.LruLen,
		lruList:          lrulist.New(),
		areaBlackHitProm: prom.New().WithCounter("filter_area_black_hit", []string{"name"}),
		spyRPC:           spyrpc.New(nil),
		aich:             make(chan func(), 1024),
		hbaseCh:          cache.New(1, 1024),
		cacheCh:          cache.New(1, 10240),
		aiDelayTick:      time.Duration(conf.Conf.Property.AIDelayTick),
		missch:           make(chan func(), 10240),
		lastTime:         time.Now().Unix(),
	}
	// 初始化业务
	if err = s.load(); err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	// reload 业务列表、敏感词、白名单、ai白名单
	go s.loadproc()
	go s.lrucleanproc()
	go s.aichproc()
	log.Info("pcre_config: %s", pcre.ConfigAll())
	go s.eventproc()
	if conf.Conf.Infoc != nil {
		s.infoc = infoc.New(conf.Conf.Infoc)
	}
	return
}

func (s *Service) load() (err error) {
	if err = s.areas.Load(context.TODO(), s.dao.AreaList); err != nil {
		return
	}
	if err = s.filters.Load(context.TODO(), s.dao.FilterAreas, s.areas); err != nil {
		return
	}
	if err = s.whites.Load(context.TODO(), s.dao.WhiteAreas, s.areas); err != nil {
		return
	}
	if err = s.whiteMids.LoadWhite(context.TODO(), s.dao.AiWhite); err != nil {
		return
	}
	return
}

func (s *Service) loadproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("loadproc panic %v : %s", x, debug.Stack())
			go s.loadproc()
		}
	}()
	var (
		err        error
		lastTime   int64
		minTick    = time.NewTicker(time.Duration(conf.Conf.Property.ReloadTick))
		secondTick = time.NewTicker(time.Duration(conf.Conf.Property.SecondReloadTick))
	)
	for {
		select {
		case <-secondTick.C:
			if lastTime, err = s.dao.AreaLastTime(context.Background()); err != nil {
				log.Error("s.dao.AreaLastTime() err(%+v)", err)
				continue
			}
			if lastTime <= s.lastTime {
				continue
			}
			s.lastTime = lastTime
			log.Info("reload filter by admin update")
			if err = s.filters.Load(context.TODO(), s.dao.FilterAreas, s.areas); err != nil {
				log.Error("s.filters.Load() err(%+v)", err)
				continue
			}
			log.Info("reload filter by admin update end")
		case <-minTick.C:
			log.Info("reload tick (%s) start", time.Duration(conf.Conf.Property.ReloadTick))
			if err = s.load(); err != nil {
				log.Error("load failed (%+v)", err)
			}
			log.Info("reload tick (%s) end", time.Duration(conf.Conf.Property.ReloadTick))
		}
	}
}

func (s *Service) lrucleanproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("lrucleanproc panic %v : %s", x, debug.Stack())
			go s.lrucleanproc()
		}
	}()
	for {
		s.lruLock.Lock()
		s.lruList.Empty()
		s.lruLock.Unlock()
		time.Sleep(time.Duration(conf.Conf.Property.LruTick))
	}
}

// Ping service ping.
func (s *Service) Ping(c context.Context) (err error) {
	if s.dao != nil {
		err = s.dao.Ping(c)
	}
	return
}

// Close close service.
func (s *Service) Close() {
	if s.dao != nil {
		s.dao.Close()
	}
}

func (s *Service) addAICh(f func()) {
	select {
	case s.aich <- f:
	default:
		log.Error("s.AICh is full")
	}
}

func (s *Service) aichproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("aichproc panic %v : %s", x, debug.Stack())
			go s.aichproc()
		}
	}()
	for {
		f := <-s.aich
		time.AfterFunc(s.aiDelayTick, f)
	}
}

func (s *Service) addEvent(f func()) {
	select {
	case s.missch <- f:
	default:
		log.Warn("eventproc chan full")
	}
}

// eventproc is a routine for executing closure.
func (s *Service) eventproc() {
	for {
		f := <-s.missch
		f()
	}
}
