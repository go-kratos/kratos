package service

import (
	"context"

	"go-common/app/interface/main/growup/conf"
	"go-common/app/interface/main/growup/dao"
)

// Service is growup service
type Service struct {
	conf *conf.Config
	dao  *dao.Dao
	sf   *SnowFlake
}

// New fn
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf: c,
		dao:  dao.New(c),
		sf:   NewSnowFlake(),
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
