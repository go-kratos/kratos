package server

import (
	"go-common/app/service/main/passport-auth/conf"
	"go-common/app/service/main/passport-auth/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC rpc struct
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

// DelTokenCache query token
func (r *RPC) DelTokenCache(c context.Context, token string, res *struct{}) (err error) {
	err = r.s.DelTokenCache(c, token)
	return
}

// DelCookieCache del cookie
func (r *RPC) DelCookieCache(c context.Context, session string, res *struct{}) (err error) {
	err = r.s.DelCookieCache(c, session)
	return
}
