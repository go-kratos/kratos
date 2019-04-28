package service

import (
	"context"
	"sync"
	"time"

	filterrpc "go-common/app/service/main/filter/rpc/client"
	"go-common/app/service/main/push-strategy/conf"
	"go-common/app/service/main/push-strategy/dao"
	"go-common/app/service/main/push-strategy/model"
	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/cache"
	"go-common/library/log"
)

// Service is service.
type Service struct {
	c          *conf.Config
	dao        *dao.Dao
	wg         sync.WaitGroup
	cache      *cache.Cache
	apps       map[int64]*pushmdl.APP
	businesses map[int64]*pushmdl.Business
	filterRPC  *filterrpc.Service
	closed     bool
	settings   map[int64]map[int]int
	midCh      chan *model.MidChan
}

const (
	_dbBatch = 50000
)

// New .
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:          c,
		dao:        dao.New(c),
		cache:      cache.New(1, c.Cfg.CacheSize),
		apps:       make(map[int64]*pushmdl.APP),
		businesses: make(map[int64]*pushmdl.Business),
		filterRPC:  filterrpc.New(c.FilterRPC),
		settings:   make(map[int64]map[int]int),
		midCh:      make(chan *model.MidChan, 5120),
	}
	s.loadApps()
	s.loadBusiness()
	s.loadUserSetting()
	go s.loadAppsproc()
	go s.loadBusinessproc()
	go s.loadUserSettingproc()
	for i := 0; i < s.c.Cfg.HandleMidGoroutines; i++ {
		s.wg.Add(1)
		go s.saveMidproc()
	}
	return s
}

func (s *Service) loadApps() (res map[int64]*pushmdl.APP, err error) {
	if res, err = s.dao.Apps(context.Background()); err != nil {
		return
	}
	if len(res) > 0 {
		s.apps = res
	}
	return
}

func (s *Service) loadAppsproc() {
	for {
		if res, err := s.loadApps(); err != nil || len(res) == 0 {
			log.Error("s.loadApps() no apps error(%v)", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		time.Sleep(time.Duration(s.c.Cfg.LoadBusinessInteval))
	}
}

func (s *Service) loadBusiness() (res map[int64]*pushmdl.Business, err error) {
	if res, err = s.dao.Businesses(context.Background()); err != nil {
		return
	}
	if len(res) > 0 {
		s.businesses = res
	}
	return
}

func (s *Service) loadBusinessproc() {
	for {
		time.Sleep(time.Duration(s.c.Cfg.LoadBusinessInteval))
		if res, err := s.loadBusiness(); err != nil || len(res) == 0 {
			log.Error("s.loadBusiness() no business")
			time.Sleep(time.Second)
		}
	}
}

func (s *Service) loadUserSettingproc() {
	for {
		time.Sleep(time.Duration(s.c.Cfg.LoadSettingsInteval))
		if err := s.loadUserSetting(); err != nil {
			log.Error("s.loadUserSetting() no settings")
			time.Sleep(time.Second)
			continue
		}
	}
}

func (s *Service) loadUserSetting() (err error) {
	maxid, err := s.dao.MaxSettingID(context.Background())
	if err != nil {
		log.Error("s.dao.MaxSettingID() error(%v)", err)
		return
	}
	log.Info("max setting id(%d)", maxid)
	var (
		ss  map[int64]map[int]int
		res = make(map[int64]map[int]int)
	)
	for i := int64(0); i <= maxid; i += _dbBatch {
		for j := 0; j < _retry; j++ {
			if ss, err = s.dao.SettingsByRange(context.Background(), i, i+_dbBatch); err == nil {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if err != nil {
			log.Error("s.dao.SEttingsByRange(%d,%d) error(%v)", i, i+_dbBatch, err)
			return
		}
		if len(ss) == 0 {
			continue
		}
		for mid, data := range ss {
			res[mid] = data
		}
	}
	if len(res) > 0 {
		s.settings = res
	}
	log.Info("loadUserSetting count(%d)", len(res))
	return
}

// Ping .
func (s *Service) Ping(c context.Context) (err error) {
	err = s.dao.Ping(c)
	return
}

// Close .
func (s *Service) Close() {
	s.closed = true
	close(s.midCh)
	s.wg.Wait()
	s.dao.Close()
}
