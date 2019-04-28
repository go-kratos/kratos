package ping

import (
	"context"

	"go-common/app/admin/main/app/conf"
	pgdao "go-common/app/admin/main/app/dao/audit"
)

// Service dao
type Service struct {
	pgDao *pgdao.Dao
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		pgDao: pgdao.New(c),
	}
	return
}

// Ping ping
func (s *Service) Ping(c context.Context) (err error) {
	return s.pgDao.PingDB(c)
}
