package server

import (
	"go-common/app/service/main/assist/conf"
	"go-common/app/service/main/assist/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC rpc.
type RPC struct {
	s *service.Service
}

// New new rpc server.
func New(c *conf.Config, s *service.Service) (svr *rpc.Server) {
	r := &RPC{s: s}
	svr = rpc.NewServer(c.RPCServer)
	if err := svr.Register(r); err != nil {
		panic(err)
	}
	return
}

// Auth check connection success.
func (r *RPC) Auth(c context.Context, arg *rpc.Auth, res *struct{}) (err error) {
	return
}

// Ping check connection success.
func (r *RPC) Ping(c context.Context, arg *struct{}, res *struct{}) (err error) {
	return
}
