package service

import (
	"context"

	"go-common/app/admin/openplatform/sug/conf"
	"go-common/app/admin/openplatform/sug/dao"
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
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
