package service

import (
	"context"

	memrpc "go-common/app/service/main/member/api/gorpc"
	"go-common/app/service/main/relation/conf"
	"go-common/app/service/main/relation/dao"
	"go-common/app/service/main/relation/model"
	"go-common/library/log"
)

var (
	_emptyFollowings   = make([]*model.Following, 0)
	_emptyFollowingMap = make(map[int64]*model.Following)
)

const (
	// UserBlockedStatus -2 is blocked.
	UserBlockedStatus = -2
	// UserRank value
	UserRank = 5000
)

// Service struct of service.
type Service struct {
	dao *dao.Dao
	// conf
	c *conf.Config
	// cache
	missch    chan func()
	inCh      chan interface{}
	memberRPC *memrpc.Service
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:         c,
		dao:       dao.New(c),
		missch:    make(chan func(), 10240),
		inCh:      make(chan interface{}, 10240),
		memberRPC: memrpc.New(c.RPCClient2.Member),
	}
	go s.cacheproc()
	go s.infocproc()
	return
}

// Ping check server ok
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
}

// addCache
func (s *Service) addCache(f func()) {
	select {
	case s.missch <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc is a routine for executing closure.
func (s *Service) cacheproc() {
	for {
		f := <-s.missch
		f()
	}
}
