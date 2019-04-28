package service

import (
	"context"

	"go-common/app/service/main/open/conf"
	"go-common/app/service/main/open/dao"
)

// Service biz service def.
type Service struct {
	c          *conf.Config
	dao        *dao.Dao
	appsecrets map[string]string
	appIDs     map[string]int64 //map[appkey]appid
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	s.appIDs = map[string]int64{}
	s.loadAppSecrets()
	go s.loadAppSecretsproc()
	return
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close close all dao.
func (s *Service) Close() {
	s.dao.Close()
}
