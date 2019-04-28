package service

import (
	"context"

	"go-common/app/admin/ep/melloi/conf"
	"go-common/app/admin/ep/melloi/dao"

	"github.com/jinzhu/gorm"
	"github.com/robfig/cron"
)

// Service struct
type Service struct {
	c    *conf.Config
	dao  *dao.Dao
	DB   *gorm.DB
	cron *cron.Cron
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:    c,
		dao:  dao.New(c),
		cron: cron.New(),
	}
	s.cron.Start()
	return s
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

//Ping test interface
func (s *Service) Ping(c context.Context) (err error) {
	err = s.dao.Ping(c)
	return
}
