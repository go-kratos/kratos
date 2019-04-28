package service

import (
	"context"
	"math/rand"
	"sort"
	"sync/atomic"
	"time"

	"go-common/app/service/main/tv/internal/conf"
	"go-common/app/service/main/tv/internal/dao"
	"go-common/app/service/main/tv/internal/model"
	"go-common/library/log"
)

// Service struct
type Service struct {
	c       *conf.Config
	dao     *dao.Dao
	ppcsMap atomic.Value
	r       *rand.Rand
	eventch chan func()
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:       c,
		dao:     dao.New(c),
		r:       rand.New(rand.NewSource(time.Now().Unix())),
		eventch: make(chan func(), 1024),
	}
	s.initPanels()
	go s.ppcproc()
	go s.eventproc()
	go s.makeupproc()
	return s
}

// Ping Service
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) loadPanels() (map[int8][]*model.PanelPriceConfig, error) {
	pcs, _, err := s.dao.PriceConfigsByStatus(context.TODO(), 0)
	if err != nil {
		log.Error("s.dao.PriceConfigsByStatus err(%v)", err)
		return nil, err
	}
	ppcs := make([]*model.PanelPriceConfig, 0, len(pcs))
	for _, pc := range pcs {
		ppc := new(model.PanelPriceConfig)
		ppc.CopyFromPriceConfig(pc)
		ppcs = append(ppcs, ppc)
	}
	spcs, _, err := s.dao.SaledPriceConfigsByStatus(context.TODO(), 0)
	if err != nil {
		log.Error("s.dao.SaledPriceConfigsByStatus err(%v)", err)
		return nil, err
	}
	sppcs := make([]*model.PanelPriceConfig, 0, len(spcs))
	for _, spc := range spcs {
		sppc := new(model.PanelPriceConfig)
		sppc.CopyFromPriceConfig(spc)
		sppcs = append(sppcs, sppc)
	}
	ppcsMap := make(map[int8][]*model.PanelPriceConfig)
	for _, sppc := range sppcs {
		if _, ok := ppcsMap[sppc.SuitType]; !ok {
			ppcsMap[sppc.SuitType] = make([]*model.PanelPriceConfig, 0)
		}
		included := false
		for _, ppc := range ppcs {
			if sppc.Pid == ppc.ID {
				included = true
				break
			}
		}
		if included {
			ppcsMap[sppc.SuitType] = append(ppcsMap[sppc.SuitType], sppc)
		}
	}
	for _, ppc := range ppcs {
		if _, ok := ppcsMap[ppc.SuitType]; !ok {
			ppcsMap[ppc.SuitType] = make([]*model.PanelPriceConfig, 0)
		}
		included := false
		for _, sppc := range sppcs {
			if ppc.ID == sppc.Pid && ppc.SuitType == sppc.SuitType {
				included = true
				sppc.OriginPrice = ppc.Price
				break
			}
		}
		if !included {
			ppcsMap[ppc.SuitType] = append(ppcsMap[ppc.SuitType], ppc)
		}
	}
	for st := range ppcsMap {
		sort.Slice(ppcsMap[st], func(l, r int) bool {
			lp := ppcsMap[st][l]
			rp := ppcsMap[st][r]
			if lp.SubType > rp.SubType {
				return true
			}
			if lp.SubType < rp.SubType {
				return false
			}
			return lp.Month > rp.Month
		})
	}
	return ppcsMap, nil
}

func (s *Service) initPanels() {
	var (
		ppcsMap map[int8][]*model.PanelPriceConfig
		err     error
	)
	if ppcsMap, err = s.loadPanels(); err != nil {
		log.Error("s.initPanels() err(%+v)", err)
		panic(err)
	}
	if ppcsMap == nil {
		panic(nil)
	}
	s.ppcsMap.Store(ppcsMap)
}

// pcproc auto load price configs periodically
func (s *Service) ppcproc() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("s.pcproc err(%v) ", err)
			go s.ppcproc()
			log.Info("s.pcproc recover")
		}
	}()
	log.Info("s.ppcproc start")
	var (
		td  time.Duration
		err error
	)
	if td, err = time.ParseDuration(s.c.Ticker.PanelRefreshDuration); err != nil {
		panic(err)
	}
	ticker := time.NewTicker(td)
	for range ticker.C {
		ppcsMap, err := s.loadPanels()
		if err != nil {
			log.Info("s.loadPanels() err(%+v)", err)
			continue
		}
		log.Info("s.ppcproc updatePcsMap(%+v)", ppcsMap)
		s.ppcsMap.Store(ppcsMap)
	}
}

func (s *Service) eventproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.eventproc panic(%v)", x)
			go s.eventproc()
			log.Info("service.eventproc recover")
		}
	}()
	for {
		f := <-s.eventch
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
	case s.eventch <- f:
	default:
		log.Error("service.missproc chan full")
	}
}

// make up tv-vip order status
func (s *Service) makeupproc() {
	var (
		td  time.Duration
		err error
	)
	if td, err = time.ParseDuration(s.c.Ticker.UnpaidRefreshDuratuion); err != nil {
		panic(err)
	}

	ticker := time.NewTicker(td)
	defer func() {
		if err := recover(); err != nil {
			log.Error("service.makeUpOrderStatus panic(%v)", err)
			go s.makeupproc()
			log.Info("service.makeUpOrderStatus recover")
		}
	}()
	log.Info("TV Vip Order Status Make Up Event Start!")
	for range ticker.C {
		log.Info("TV Vip Order Status Making Up Event!")
		if err := s.MakeUpOrderStatus(); err != nil {
			log.Error("s.MakeUpOrderStatus() err(%+v)", err)
		}
	}
}
