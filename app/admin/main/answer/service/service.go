package service

import (
	"context"

	"go-common/app/admin/main/answer/conf"
	"go-common/app/admin/main/answer/dao"
	"go-common/library/cache"
)

// Service struct of service.
type Service struct {
	c         *conf.Config
	dao       *dao.Dao
	eventChan *cache.Cache
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:         c,
		dao:       dao.New(c),
		eventChan: cache.New(1, 10240),
	}
	s.generate(context.Background(), x, 0, len(x)-1)
	return
}

// Close dao.
func (s *Service) Close() {
	s.dao.Close()
}
