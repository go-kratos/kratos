package service

import (
	"context"

	"go-common/app/admin/main/space/conf"
	"go-common/app/admin/main/space/dao"
	relrpc "go-common/app/service/main/relation/rpc/client"
)

// Service biz service def.
type Service struct {
	c        *conf.Config
	dao      *dao.Dao
	relation *relrpc.Service
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:        c,
		dao:      dao.New(c),
		relation: relrpc.New(c.RelationRPC),
	}
	return s
}

// Ping .
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}
