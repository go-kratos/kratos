package rpc

import (
	"go-common/app/service/main/figure/conf"
	"go-common/app/service/main/figure/model"
	"go-common/app/service/main/figure/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC server def.
type RPC struct {
	s *service.Service
}

// New init rpc.
func New(c *conf.Config, s *service.Service) (svr *rpc.Server) {
	r := &RPC{s: s}
	svr = rpc.NewServer(c.RPCServer)
	if err := svr.Register(r); err != nil {
		panic(err)
	}
	return
}

// Ping check rpc server health.
func (r *RPC) Ping(c context.Context, arg *struct{}, res *struct{}) (err error) {
	return
}

// UserFigure get user figure & rank info.
func (r *RPC) UserFigure(c context.Context, arg *model.ArgUserFigure, res *model.FigureWithRank) (err error) {
	var fr *model.FigureWithRank
	if fr, err = r.s.FigureWithRank(c, arg.Mid); fr != nil {
		*res = *fr
	}
	return
}
