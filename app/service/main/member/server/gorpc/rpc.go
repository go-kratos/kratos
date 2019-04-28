package gorpc

import (
	"go-common/app/service/main/member/conf"
	"go-common/app/service/main/member/model"
	"go-common/app/service/main/member/service"
	"go-common/app/service/main/member/service/block"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC is.
type RPC struct {
	s     *service.Service
	block *block.Service
}

// New new rpc server.
func New(c *conf.Config, s *service.Service) *rpc.Server {
	r := &RPC{
		s:     s,
		block: s.BlockImpl(),
	}
	svr := rpc.NewServer(c.RPCServer)
	if err := svr.Register(r); err != nil {
		panic(err)
	}
	return svr
}

// Ping check connection success.
func (r *RPC) Ping(c context.Context, arg *struct{}, res *struct{}) (err error) {
	return
}

// ---  exp --- //

// Exp  get user exp.
func (r *RPC) Exp(c context.Context, arg *model.ArgMid2, res *model.LevelInfo) (err error) {
	v, err := r.s.Exp(c, arg.Mid)
	if err == nil && v != nil {
		*res = *v
	}
	return
}

// Level  get user exp.
func (r *RPC) Level(c context.Context, arg *model.ArgMid2, res *model.LevelInfo) (err error) {
	v, err := r.s.Level(c, arg.Mid)
	if err == nil && v != nil {
		*res = *v
	}
	return
}

// UpdateExp user exp.
func (r *RPC) UpdateExp(c context.Context, arg *model.ArgAddExp, res *struct{}) (err error) {
	err = r.s.UpdateExp(c, arg)
	return
}

// Log get user exp log.
func (r *RPC) Log(c context.Context, arg *model.ArgMid2, res *[]*model.UserLog) (err error) {
	*res, err = r.s.ExpLog(c, arg.Mid, arg.RealIP)
	return
}

// Stat  get user exp stat.
func (r *RPC) Stat(c context.Context, arg *model.ArgMid2, res *model.ExpStat) (err error) {
	v, err := r.s.Stat(c, arg.Mid)
	if err == nil && v != nil {
		*res = *v
	}
	return
}
