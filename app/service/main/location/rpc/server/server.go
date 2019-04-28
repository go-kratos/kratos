package server

import (
	"go-common/app/service/main/location/conf"
	"go-common/app/service/main/location/model"
	"go-common/app/service/main/location/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC struct
type RPC struct {
	s *service.Service
}

// New init rpc server.
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

// Archive get aid auth.
func (r *RPC) Archive(c context.Context, a *model.Archive, res *int64) (err error) {
	var play int64
	if play, err = r.s.Auth(c, a.Aid, a.Mid, a.IP, a.CIP); err == nil {
		*res = play
	}
	return
}

// Archive2 get aid auth.
func (r *RPC) Archive2(c context.Context, a *model.Archive, res *model.Auth) (err error) {
	var auth *model.Auth
	if auth, err = r.s.Archive2(c, a.Aid, a.Mid, a.IP, a.CIP); err == nil {
		*res = *auth
	}
	return
}

// Group get gid auth.
func (r *RPC) Group(c context.Context, a *model.Group, res *model.Auth) (err error) {
	auth := r.s.AuthGID(c, a.Gid, a.Mid, a.IP, a.CIP)
	if auth != nil {
		*res = *auth
	}
	return
}

// AuthPIDs check if ip in pids.
func (r *RPC) AuthPIDs(c context.Context, a *model.ArgPids, res *map[int64]*model.Auth) (err error) {
	*res, err = r.s.AuthPIDs(c, a.Pids, a.IP, a.CIP)
	return
}

// Info get IP info
func (r *RPC) Info(c context.Context, a *model.ArgIP, res *model.Info) (err error) {
	var ti *model.Info
	if ti, err = r.s.Info(c, a.IP); err == nil {
		*res = *ti
	}
	return
}

// Infos get IPs infos.
func (r *RPC) Infos(c context.Context, a []string, res *map[string]*model.Info) (err error) {
	*res, err = r.s.Infos(c, a)
	return
}

// InfoComplete get IP whole info.
func (r *RPC) InfoComplete(c context.Context, a *model.ArgIP, res *model.InfoComplete) (err error) {
	var info *model.InfoComplete
	if info, err = r.s.InfoComplete(c, a.IP); err == nil {
		*res = *info
	}
	return
}

// InfosComplete get ips whole infos.
func (r *RPC) InfosComplete(c context.Context, a []string, res *map[string]*model.InfoComplete) (err error) {
	*res, err = r.s.InfosComplete(c, a)
	return
}
