package service

import (
	"context"
	"go-common/app/interface/main/shorturl/conf"
	shortdao "go-common/app/interface/main/shorturl/dao"
	"go-common/library/log"
)

// Service service struct
type Service struct {
	shortd *shortdao.Dao
}

// New new service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		shortd: shortdao.New(c),
	}
	return
}

// Ping ping service.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.shortd.Ping(c); err != nil {
		log.Error("s.dao.Ping error(%v)", err)
	}
	return
}
