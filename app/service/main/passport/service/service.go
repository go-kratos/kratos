package service

import (
	"context"

	"go-common/app/service/main/passport/conf"
	"go-common/app/service/main/passport/dao"
)

// Service struct of service.
type Service struct {
	d *dao.Dao
	// conf
	c *conf.Config
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
		d: dao.New(c),
	}
	return
}

// Close dao.
func (s *Service) Close() {
	s.d.Close()
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.d.Ping(c); err != nil {
		return
	}
	return
}
