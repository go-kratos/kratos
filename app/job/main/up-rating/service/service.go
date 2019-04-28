package service

import (
	"context"

	"go-common/app/job/main/up-rating/conf"
	"go-common/app/job/main/up-rating/dao"
	"go-common/library/log"
)

// Service struct
type Service struct {
	conf *conf.Config
	dao  *dao.Dao
}

// New fn
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf: c,
		dao:  dao.New(c),
	}
	log.Info("service start")
	return s
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close close the service
func (s *Service) Close() {
	s.dao.Close()
}
