package service

import (
	"context"

	"go-common/app/interface/main/web-feed/conf"
	"go-common/app/interface/main/web-feed/dao"
	accrpc "go-common/app/service/main/account/rpc/client"
	feedrpc "go-common/app/service/main/feed/rpc/client"
	"go-common/library/cache"
)

// Service service struct info
type Service struct {
	c       *conf.Config
	dao     *dao.Dao
	feedRPC *feedrpc.Service
	accRPC  *accrpc.Service3
	cache   *cache.Cache
}

// New .
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:       c,
		dao:     dao.New(c),
		feedRPC: feedrpc.New(c.FeedRPC),
		accRPC:  accrpc.New3(c.AccountRPC),
		cache:   cache.New(1, 1024),
	}
	return s
}

// Close closes dao.
func (s *Service) Close() {
	s.dao.Close()
}

// Ping is check server ping.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}
