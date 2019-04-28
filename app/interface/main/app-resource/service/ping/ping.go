package ping

import (
	"context"

	"go-common/app/interface/main/app-resource/conf"
	pgdao "go-common/app/interface/main/app-resource/dao/plugin"
)

type Service struct {
	pgDao *pgdao.Dao
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		pgDao: pgdao.New(c),
	}
	return
}

func (s *Service) Ping(c context.Context) (err error) {
	return s.pgDao.PingDB(c)
}
