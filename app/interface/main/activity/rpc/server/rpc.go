package server

import (
	"go-common/app/interface/main/activity/conf"
	matmdl "go-common/app/interface/main/activity/model/like"
	match "go-common/app/interface/main/activity/service/like"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
	"go-common/library/net/rpc/interceptor"
)

// RPC struct info.
type RPC struct {
	s *match.Service
}

// New new rpc server.
func New(c *conf.Config, s *match.Service) (svr *rpc.Server) {
	r := &RPC{s: s}
	svr = rpc.NewServer(c.RPCServer)
	in := interceptor.NewInterceptor("")
	svr.Interceptor = in
	if err := svr.Register(r); err != nil {
		panic(err)
	}
	return
}

// Ping check connection success.
func (r *RPC) Ping(c context.Context, arg *struct{}, res *struct{}) (err error) {
	return
}

// Matchs return matchs by sid
func (r *RPC) Matchs(c context.Context, a *matmdl.ArgMatch, res *[]*matmdl.Match) (err error) {
	*res, err = r.s.Match(c, a.Sid)
	return
}

// SubjectUp up act_subject cahce info and act_subject maxID cache.
func (r *RPC) SubjectUp(c context.Context, a *matmdl.ArgSubjectUp, res *struct{}) (err error) {
	return r.s.SubjectUp(c, a.Sid)
}

// ActSubject get act subject info.
func (r *RPC) ActSubject(c context.Context, a *matmdl.ArgActSubject, res *matmdl.SubjectItem) (err error) {
	var rr *matmdl.SubjectItem
	if rr, err = r.s.ActSubject(c, a.Sid); err == nil {
		*res = *rr
	}
	return
}

// LikeUp up likes cache info and like maxID cache
func (r *RPC) LikeUp(c context.Context, a *matmdl.ArgLikeUp, res *struct{}) (err error) {
	return r.s.LikeUp(c, a.Lid)
}

// AddLikeCtimeCache like ctime cache.
func (r *RPC) AddLikeCtimeCache(c context.Context, a *matmdl.ArgLikeUp, res *struct{}) (err error) {
	return r.s.AddLikeCtimeCache(c, a.Lid)
}

// DelLikeCtimeCache del like ctime cache
func (r *RPC) DelLikeCtimeCache(c context.Context, item *matmdl.ArgLikeItem, res *struct{}) (err error) {
	return r.s.DelLikeCtimeCache(c, item.ID, item.Sid, item.Type)
}

// ActProtocol .
func (r *RPC) ActProtocol(c context.Context, a *matmdl.ArgActProtocol, res *matmdl.SubProtocol) (err error) {
	var rr *matmdl.SubProtocol
	if rr, err = r.s.ActProtocol(c, a); err == nil {
		*res = *rr
	}
	return
}
