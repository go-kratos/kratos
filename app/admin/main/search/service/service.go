package service

import (
	"context"

	"go-common/app/admin/main/search/conf"
	"go-common/app/admin/main/search/dao"
	"go-common/app/admin/main/search/model"
)

// Service struct of service.
type Service struct {
	c         *conf.Config
	dao       *dao.Dao
	queryConf map[string]*model.QueryConfDetail
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	s.loadQueryConf()
	go s.loadQueryConfproc()
	return
}

// Ping .
func (s *Service) Ping(c context.Context) error {
	return s.dao.Ping(c)
}
