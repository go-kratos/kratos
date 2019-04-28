package server

import (
	"go-common/app/service/main/seq-server/conf"
	"go-common/app/service/main/seq-server/model"
	"go-common/app/service/main/seq-server/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC rpc.
type RPC struct {
	s *service.Service
}

// New creates rpc server.
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
	if err = r.s.Ping(c); err != model.ErrBusinessNotReady {
		err = nil
	}
	return
}

// ID return id.
func (r *RPC) ID(c context.Context, a *model.ArgBusiness, res *int64) (err error) {
	*res, err = r.s.ID(c, a.BusinessID, a.Token)
	return
}

// ID32 return id32.
func (r *RPC) ID32(c context.Context, a *model.ArgBusiness, res *int32) (err error) {
	*res, err = r.s.ID32(c, a.BusinessID, a.Token)
	return
}

// CheckVersion check db health.
func (r *RPC) CheckVersion(c context.Context, arg *struct{}, res *model.SeqVersion) (err error) {
	var seqVer *model.SeqVersion
	if seqVer, err = r.s.CheckVersion(c); err == nil {
		*res = *seqVer
	}
	return
}
