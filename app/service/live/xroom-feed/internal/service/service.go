package service

import (
	"context"
	"sync/atomic"

	"go-common/app/service/live/xroom-feed/internal/dao"
	"go-common/library/conf/paladin"
)

// Service service.
type Service struct {
	ac       *paladin.Map
	dao      *dao.Dao
	recCache atomic.Value
}

// New new a service and return.
func New() (s *Service) {
	var ac = new(paladin.TOML)
	if err := paladin.Watch("application.toml", ac); err != nil {
		panic(err)
	}
	s = &Service{
		ac:  ac,
		dao: dao.New(),
	}
	s.loadPoolConf()
	go s.poolConfProc()
	return s
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}
