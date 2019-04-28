package ping

import (
	"context"

	"go-common/app/interface/main/app-tag/conf"
	regiondao "go-common/app/interface/main/app-tag/dao/region"
)

type Service struct {
	regionDao *regiondao.Dao
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		regionDao: regiondao.New(c),
	}
	return
}

// Ping is check server ping.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.regionDao.Ping(c); err != nil {
		return
	}
	return
}
