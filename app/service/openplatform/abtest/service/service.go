package service

import (
	"context"
	"sync"

	"go-common/app/service/openplatform/abtest/conf"
	"go-common/app/service/openplatform/abtest/dao"
)

// Service struct of service.
type Service struct {
	d *dao.Dao
	// conf
	c *conf.Config

	// groupId => AB.id => AB
	// abCache    map[int]map[int]*model.AB
	// _versionID map[int]int64
	// mutex      sync.RWMutex
	abCache    sync.Map
	_versionID sync.Map
	keyList    sync.Map
	// statCacheO map[int]map[int]int
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
		d: dao.New(c),
	}
	// s.statCache = make(map[int]map[int]int)
	// s.statCacheO = make(map[int]map[int]int)
	go syncStart(s)
	return
}

// Close dao.
func (s *Service) Close() {
	s.d.Close()
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.d.Ping(c); err != nil {
		return
	}
	return
}
