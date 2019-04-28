package service

import (
	"context"

	"go-common/app/service/live/dao-anchor/conf"
	"go-common/app/service/live/dao-anchor/dao"
	consumerV1 "go-common/app/service/live/dao-anchor/service/consumer/v1"
	"go-common/app/service/live/dao-anchor/service/v0"
	"go-common/app/service/live/dao-anchor/service/v1"
)

// Service struct
type Service struct {
	c              *conf.Config
	dao            *dao.Dao
	v1svc          *v1.DaoAnchorService
	consumerSvc    *consumerV1.ConsumerService
	createCacheSvc *v0.CreateDataService
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:              c,
		dao:            dao.New(c),
		v1svc:          v1.NewDaoAnchorService(c),
		consumerSvc:    consumerV1.NewConsumerService(c),
		createCacheSvc: v0.NewCreateDataService(c),
	}
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

// CreateDataSvc return dao-anchor CreateDataService
func (s *Service) CreateDataSvc() *v0.CreateDataService {
	return s.createCacheSvc
}

// V1Svc return dao-anchor v1 service
func (s *Service) V1Svc() *v1.DaoAnchorService {
	return s.v1svc
}
