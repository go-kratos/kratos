package service

import (
	"context"
	"time"

	"go-common/app/job/openplatform/open-market/conf"
	"go-common/app/job/openplatform/open-market/dao"
)

// Service struct of service.
type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	d := dao.New(c)
	s = &Service{
		c:   c,
		dao: d,
	}
	go s.fetchData()
	return
}

// Close close service.
func (s *Service) Close() (err error) {
	s.dao.Close()
	return
}

// Ping ping service.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

func (s *Service) fetchData() {
	for {
		s.marketProc()
		time.Sleep(time.Hour * 24)
	}
}
