package service

import (
	"context"

	"go-common/app/interface/main/feedback/conf"
	"go-common/app/interface/main/feedback/dao"
	locrpc "go-common/app/service/main/location/rpc/client"
)

// Service struct.
type Service struct {
	// dao
	dao *dao.Dao
	// conf
	c *conf.Config
	// rpc
	locationRPC *locrpc.Service
}

// New new Tag service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
		// rpc
		locationRPC: locrpc.New(c.LocationRPC),
	}
	// init dao
	s.dao = dao.New(c)
	return
}

// Ping check server ok
func (s *Service) Ping(c context.Context) (err error) {
	return
}
