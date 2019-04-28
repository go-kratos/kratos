package service

import (
	"context"
	"go-common/app/service/openplatform/ticket-sales/conf"
	"go-common/app/service/openplatform/ticket-sales/dao"
)

// Service http service
type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

// New for new service obj
func New(c *conf.Config) *Service {
	s := &Service{
		c:   c,
		dao: dao.New(c),
	}
	return s
}

//Get get config
func (s *Service) Get() (*conf.Config, *dao.Dao) {
	return s.c, s.dao
}

// Ping check server ok
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
}
