package service

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"go-common/app/interface/openplatform/monitor-end/conf"
	"go-common/app/interface/openplatform/monitor-end/dao"
	"go-common/app/interface/openplatform/monitor-end/model"
	"go-common/app/interface/openplatform/monitor-end/model/monitor"
	"go-common/app/interface/openplatform/monitor-end/model/prom"
	"go-common/library/log"
)

var (
	_notConsumedErr = errors.New("未启动消费")
	_pausedNowErr   = errors.New("暂停消费中")
	_closedNowErr   = errors.New("已停止消费")
	_consumedNowErr = errors.New("已启动消费")
)

// Service struct
type Service struct {
	c        *conf.Config
	dao      *dao.Dao
	mh       *monitor.MonitorHandler
	consumer *Consumer
	consumed bool
	// settings
	groups      map[int64]*model.Group
	targets     map[int64]*model.Target
	targetKeys  map[string]*model.Target
	products    map[int64]*model.Product
	productKeys map[string]*model.Product
	newTargets  map[string]*model.Target
	naProducts  map[string]bool
	mapMutex    sync.Mutex
	// infoc
	infoCh chan interface{}
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:           c,
		dao:         dao.New(c),
		mh:          monitor.NewMonitor(c.Monitor),
		groups:      make(map[int64]*model.Group),
		targets:     make(map[int64]*model.Target),
		targetKeys:  make(map[string]*model.Target),
		newTargets:  make(map[string]*model.Target),
		products:    make(map[int64]*model.Product),
		productKeys: make(map[string]*model.Product),
		naProducts:  make(map[string]bool),
		infoCh:      make(chan interface{}, 1024),
	}
	for c.NeedConsume && !s.consumed {
		var err error
		if s.consumer, err = NewConsumer(c); err != nil {
			log.Error("s.New.NewConsumer error(%+v)", err)
			time.Sleep(10 * time.Second)
			continue
		}
		go s.consume()
		go s.handleMsg()
		s.consumed = true
		break
	}
	prom.Init(c)
	go s.loadalertsettingproc()
	go s.infocproc()
	go s.loadNAProducts()
	// go s.alertproc()
	return
}

// Ping .
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// StartConsume .
func (s *Service) StartConsume() (err error) {
	if s.consumed {
		return _consumedNowErr
	}
	if s.consumer, err = NewConsumer(s.c); err != nil {
		log.Error("s.New.NewConsumer error(%+v)", err)
		return
	}
	go s.consume()
	go s.handleMsg()
	s.consumed = true
	return
}

// Close .
func (s *Service) Close() {
	s.dao.Close()
	s.StopConsume()
}

// StopConsume .
func (s *Service) StopConsume() error {
	if s.consumed {
		s.consumed = false
	}
	if s.consumer != nil {
		s.consumer.closed = true
	}
	return nil
}

// PauseConsume .
func (s *Service) PauseConsume(t int64) error {
	if !s.consumed {
		return _notConsumedErr
	}
	if s.consumer.closed {
		return _closedNowErr
	}
	if s.consumer.paused {
		return _pausedNowErr
	}
	s.consumer.paused = true
	s.consumer.duration = time.Second * time.Duration(t)
	return nil
}

func (s *Service) loadalertsettingproc() {
	for {
		s.loadalertsettings()
		time.Sleep(time.Minute)
	}
}

func (s *Service) loadNAProducts() {
	for {
		if s.c.Products != "" {
			var m = make(map[string]bool)
			ps := strings.Split(s.c.Products, ",")
			for _, p := range ps {
				m[p] = true
			}
			s.mapMutex.Lock()
			s.naProducts = m
			s.mapMutex.Unlock()
		}
		time.Sleep(time.Minute)
	}
}
