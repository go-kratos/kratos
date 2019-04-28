package service

import (
	"context"

	"go-common/app/service/main/search/conf"
	"go-common/app/service/main/search/dao"
)

// Service struct of service.
type Service struct {
	// conf
	c *conf.Config
	// dao
	dao *dao.Dao
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	return
}

// Ping
func (s *Service) Ping(c context.Context) error {
	return s.dao.Ping(c)
}
