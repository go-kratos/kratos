package service

import (
	"context"
	"sync/atomic"
	"time"
	"unsafe"

	"go-common/app/service/main/point/conf"
	"go-common/app/service/main/point/dao"
	"go-common/app/service/main/point/model"
	"go-common/library/log"
)

// Service struct
type Service struct {
	c              *conf.Config
	dao            *dao.Dao
	pointConf      map[int64]int64
	configLoadTick time.Duration
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:              c,
		dao:            dao.New(c),
		pointConf:      make(map[int64]int64),
		configLoadTick: time.Duration(c.Property.ConfigLoadTick),
	}
	if err := s.loadallpointconf(); err != nil {
		panic(err)
	}
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) loadallpointconf() (err error) {
	var (
		c  = context.Background()
		es []*model.VipPointConf
	)
	es, err = s.dao.AllPointConfig(c)
	if err != nil {
		log.Error("loadallpointconf allpoint conf error(%v)", err)
		return
	}
	tmp := make(map[int64]int64, len(es))
	for _, e := range es {
		tmp[e.AppID] = e.Point
	}
	p := unsafe.Pointer(&s.pointConf)
	atomic.SwapPointer(&p, unsafe.Pointer(&tmp))
	log.Info("loadallpointconf (%v) load success", tmp)
	return
}

func (s *Service) loadconfigproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("s.loadconfigproc panic(%v)", x)
			go s.loadconfigproc()
			log.Info("s.loadconfigproc recover")
		}
	}()
	for {
		time.Sleep(s.configLoadTick)
		s.loadallpointconf()
	}
}
