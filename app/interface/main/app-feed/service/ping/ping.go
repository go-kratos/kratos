package ping

import (
	"context"

	"go-common/app/interface/main/app-feed/conf"
	arcdao "go-common/app/interface/main/app-feed/dao/archive"
	adtdao "go-common/app/interface/main/app-feed/dao/audit"
)

type Service struct {
	arcDao *arcdao.Dao
	adtDao *adtdao.Dao
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		arcDao: arcdao.New(c),
		adtDao: adtdao.New(c),
	}
	return
}

// Ping is check server ping.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.adtDao.PingDB(c); err != nil {
		return
	}
	err = s.arcDao.PingMC(c)
	return
}
