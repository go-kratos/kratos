package service

import (
	"context"

	"go-common/app/admin/main/open/conf"
	"go-common/app/admin/main/open/dao"
)

// Service biz service def.
type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
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
