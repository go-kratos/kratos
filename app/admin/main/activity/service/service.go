package service

import (
	"context"
	"go-common/app/admin/main/activity/conf"
	"go-common/app/admin/main/activity/dao"
	tagrpc "go-common/app/interface/main/tag/rpc/client"
	artrpc "go-common/app/interface/openplatform/article/rpc/client"
	acccli "go-common/app/service/main/account/api"
	arcclient "go-common/app/service/main/archive/api"

	"github.com/jinzhu/gorm"
)

// Service biz service def.
type Service struct {
	c         *conf.Config
	dao       *dao.Dao
	DB        *gorm.DB
	accClient acccli.AccountClient
	tagRPC    *tagrpc.Service
	artRPC    *artrpc.Service
	arcClient arcclient.ArchiveClient
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		dao:    dao.New(c),
		tagRPC: tagrpc.New2(c.TagRPC),
		artRPC: artrpc.New(c.ArticlrRPC),
	}
	s.DB = s.dao.DB
	var err error
	if s.arcClient, err = arcclient.NewClient(c.ArcClient); err != nil {
		panic(err)
	}
	if s.accClient, err = acccli.NewClient(c.AccClient); err != nil {
		panic(err)
	}
	return s
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Wait wait all closed.
func (s *Service) Wait() {}

// Close close all dao.
func (s *Service) Close() {
	s.dao.Close()
}
