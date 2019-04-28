package server

import (
	"go-common/app/interface/main/playlist/conf"
	"go-common/app/interface/main/playlist/model"
	"go-common/app/interface/main/playlist/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC struct info.
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

// Ping check connection success
func (r *RPC) Ping(c context.Context, arg *struct{}, res *struct{}) (err error) {
	return
}

// SetStat set all stat cache(redis)
func (r *RPC) SetStat(c context.Context, arg *model.ArgStats, res *struct{}) (err error) {
	err = r.s.SetStat(c, arg.PlStat)
	return
}
