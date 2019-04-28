package rpc

import (
	"go-common/app/service/main/rank/conf"
	"go-common/app/service/main/rank/model"
	"go-common/app/service/main/rank/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC favorite rpc.
type RPC struct {
	c *conf.Config
	s *service.Service
}

// New init rpc.
func New(c *conf.Config, s *service.Service) (svr *rpc.Server) {
	r := &RPC{
		c: c,
		s: s,
	}
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

// Mget .
func (r *RPC) Mget(c context.Context, a *model.MgetReq, res *model.MgetResp) (err error) {
	var v *model.MgetResp
	if v, err = r.s.Mget(c, a); err == nil {
		*res = *v
	}
	return
}

// Sort .
func (r *RPC) Sort(c context.Context, a *model.SortReq, res *model.SortResp) (err error) {
	var v *model.SortResp
	if v, err = r.s.Sort(c, a); err == nil {
		*res = *v
	}
	return
}

// Group .
func (r *RPC) Group(c context.Context, a *model.GroupReq, res *model.GroupResp) (err error) {
	var v *model.GroupResp
	if v, err = r.s.Group(c, a); err == nil {
		*res = *v
	}
	return
}
