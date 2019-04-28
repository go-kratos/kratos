package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	accrpc "go-common/app/service/main/account/rpc/client"
	"go-common/app/service/main/spy/conf"
	"go-common/app/service/main/spy/dao"
	"go-common/app/service/main/spy/model"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

const (
	_score = 100
)

// Service biz service def.
type Service struct {
	c              *conf.Config
	dao            *dao.Dao
	missch         chan func()
	scorech        chan func()
	infomissch     chan func()
	spyConfig      map[string]interface{}
	configLoadTick time.Duration
	promBaseScore  *prom.Prom
	promEventScore *prom.Prom
	promBlockInfo  *prom.Prom
	accRPC         *accrpc.Service3
	allEventName   map[int64]string
	loadEventTick  time.Duration
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:              c,
		dao:            dao.New(c),
		missch:         make(chan func(), 10240),
		scorech:        make(chan func(), 10240),
		infomissch:     make(chan func(), 10240),
		spyConfig:      make(map[string]interface{}),
		configLoadTick: time.Duration(c.Property.ConfigLoadTick),
		promBaseScore:  prom.New().WithCounter("spy_basescore", []string{"name"}),
		promEventScore: prom.New().WithCounter("spy_eventscore", []string{"name"}),
		promBlockInfo:  prom.New().WithCounter("spy_block_info", []string{"name"}),
		accRPC:         accrpc.New3(c.RPC.Account),
		allEventName:   make(map[int64]string),
		loadEventTick:  time.Duration(c.Property.LoadEventTick),
	}
	if err := s.loadSystemConfig(); err != nil {
		panic(err)
	}
	s.loadeventname()

	go s.missproc()
	go s.infomissproc()
	go s.loadconfigproc()
	go s.loadeventproc()
	go s.updatescoreproc()
	return s
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close close all dao.
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) missproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.missproc panic(%v)", x)
			go s.missproc()
			log.Info("service.missproc recover")
		}
	}()
	for {
		f := <-s.missch
		f()
	}
}

func (s *Service) mission(f func()) {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.mission panic(%v)", x)
		}
	}()
	select {
	case s.missch <- f:
	default:
		log.Error("service.missproc chan full")
	}
}

func (s *Service) updatescoreproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.updatescoreproc panic(%v)", x)
			go s.updatescoreproc()
			log.Info("service.updatescoreproc recover")
		}
	}()
	for {
		f := <-s.scorech
		f()
	}
}

func (s *Service) updatescore(f func()) {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.updatescore panic(%v)", x)
		}
	}()
	select {
	case s.scorech <- f:
	default:
		log.Error("service.updatescore chan full")
	}
}

func (s *Service) infomissproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.infomissproc panic(%v)", x)
			go s.infomissproc()
			log.Info("service.infomissproc recover")
		}
	}()
	for {
		f := <-s.infomissch
		f()
	}
}

func (s *Service) infomission(f func()) {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.infomission panic(%v)", x)
		}
	}()
	select {
	case s.infomissch <- f:
	default:
		log.Error("service.infomissch chan full")
	}
}

func (s *Service) loadconfigproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.loadconfig panic(%v)", x)
		}
	}()
	for {
		time.Sleep(s.configLoadTick)
		s.loadSystemConfig()
	}
}

func (s *Service) loadSystemConfig() (err error) {
	var (
		cdb map[string]string
	)
	cdb, err = s.dao.Configs(context.TODO())
	if err != nil {
		log.Error("sys config db get data err(%v)", err)
		return
	}
	if len(cdb) == 0 {
		err = errors.New("sys config no data")
		return
	}
	cs := make(map[string]interface{}, len(cdb))
	for k, v := range cdb {
		switch k {
		case model.LimitBlockCount:
			t, err1 := strconv.ParseInt(v, 10, 64)
			if err1 != nil {
				log.Error("sys config err(%s,%v,%v)", model.LimitBlockCount, t, err1)
				err = err1
				return
			}
			cs[k] = t
		case model.LessBlockScore:
			tmp, err1 := strconv.ParseInt(v, 10, 8)
			if err1 != nil {
				log.Error("sys config err(%s,%v,%v)", model.LessBlockScore, tmp, err1)
				err = err1
				return
			}
			cs[k] = int8(tmp)
		case model.AutoBlock:
			tmp, err1 := strconv.ParseInt(v, 10, 8)
			if err1 != nil {
				log.Error("sys config err(%s,%v,%v)", model.AutoBlock, tmp, err1)
				err = err1
				return
			}
			cs[k] = int8(tmp)
		default:
			cs[k] = v
		}
	}
	s.spyConfig = cs
	log.Info("loadSystemConfig success(%v)", cs)
	return
}

//Config get config.
func (s *Service) Config(key string) (interface{}, bool) {
	if s.spyConfig == nil {
		return nil, false
	}
	v, ok := s.spyConfig[key]
	return v, ok
}

func (s *Service) loadeventname() (err error) {
	var (
		c  = context.Background()
		es []*model.Event
	)
	es, err = s.dao.AllEvent(c)
	if err != nil {
		log.Error("loadeventname allevent error(%v)", err)
		return
	}
	tmp := make(map[int64]string, len(es))
	for _, e := range es {
		tmp[e.ID] = e.NickName
	}
	s.allEventName = tmp
	log.Info("loadeventname (%v) load success", tmp)
	return
}

func (s *Service) loadeventproc() {
	for {
		time.Sleep(s.loadEventTick)
		s.loadeventname()
	}
}
