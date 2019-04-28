package server

import (
	"go-common/app/service/main/passport/conf"
	"go-common/app/service/main/passport/model"
	"go-common/app/service/main/passport/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC server struct
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

// Ping check connection success.
func (r *RPC) Ping(c context.Context, arg *struct{}, res *struct{}) (err error) {
	return
}

// LoginLogs get the latest limit login logs.
func (r *RPC) LoginLogs(c context.Context, arg *model.ArgLoginLogs, res *[]*model.LoginLog) (err error) {
	if ms, err := r.s.LoginLogs(c, arg.Mid, arg.Limit); err == nil {
		*res = ms
	}
	return
}
