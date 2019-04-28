package service

import (
	"context"

	"go-common/app/interface/main/up-rating/conf"
	"go-common/app/interface/main/up-rating/dao"
)

// Service is up-dao service
type Service struct {
	conf *conf.Config
	// dao dao
	dao *dao.Dao
}

// New fn
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf: c,
		dao:  dao.New(c),
	}
	return s
}

// Ping fn
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
}
