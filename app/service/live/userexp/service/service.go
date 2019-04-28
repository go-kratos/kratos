package service

import (
	"context"

	"go-common/app/service/live/userexp/conf"
	"go-common/app/service/live/userexp/dao"
	"go-common/library/cache"
)

// Service http service
type Service struct {
	c     *conf.Config
	keys  map[string]string
	dao   *dao.Dao
	cache *cache.Cache
}

// New for new service obj
func New(c *conf.Config) *Service {
	s := &Service{
		c:     c,
		keys:  map[string]string{},
		dao:   dao.New(c),
		cache: cache.New(1, 1024),
	}
	return s
}

// Ping check server ok
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
}
