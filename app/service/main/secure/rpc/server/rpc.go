package rpc

import (
	"go-common/app/service/main/secure/conf"
	model "go-common/app/service/main/secure/model"
	"go-common/app/service/main/secure/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC rpc service.
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

// Ping check connection success.
func (r *RPC) Ping(c context.Context, arg *struct{}, res *struct{}) (err error) {
	return
}

// Status rpc status.
func (r *RPC) Status(c context.Context, a *model.ArgSecure, res *model.Msg) (err error) {
	var tmp *model.Msg
	if tmp, err = r.s.Status(c, a.Mid, a.UUID); err == nil && tmp != nil {
		*res = *tmp
	}
	return
}

// CloseNotify rpc close notify.
func (r *RPC) CloseNotify(c context.Context, arg *model.ArgSecure, res *struct{}) (err error) {
	err = r.s.CloseNotify(c, arg.Mid, arg.UUID)
	return
}

// AddFeedBack  rpc add feedback.
func (r *RPC) AddFeedBack(c context.Context, arg *model.ArgFeedBack, res *struct{}) (err error) {
	err = r.s.AddFeedBack(c, arg.Mid, arg.Ts, arg.Type, arg.IP)
	return
}

// ExpectionLoc rpc expection loc list.
func (r *RPC) ExpectionLoc(c context.Context, arg *model.ArgSecure, res *[]*model.Expection) (err error) {
	*res, err = r.s.ExpectionLoc(c, arg.Mid)
	return
}
