package rpc

import (
	"go-common/app/infra/config/conf"
	"go-common/app/infra/config/service/v1"
	"go-common/app/infra/config/service/v2"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"

	"go-common/app/infra/config/model"
)

// RPC export rpc service
type RPC struct {
	s  *v1.Service
	s2 *v2.Service
}

// New new rpc server.
func New(c *conf.Config, s *v1.Service, s2 *v2.Service) (svr *rpc.Server) {
	r := &RPC{s: s, s2: s2}
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

// Push push new config change to config-service
func (r *RPC) Push(c context.Context, a *model.ArgConf, res *struct{}) (err error) {
	service := &model.Service{Name: a.App, BuildVersion: a.BuildVer, Version: a.Ver, Env: a.Env}
	err = r.s.Push(c, service)
	return
}

//SetToken update Token
func (r *RPC) SetToken(c context.Context, a *model.ArgToken, res *struct{}) (err error) {
	r.s.SetToken(c, a.App, a.Env, a.Token)
	return
}

// PushV4 push new config change to config-service
func (r *RPC) PushV4(c context.Context, a *model.ArgConf, res *struct{}) (err error) {
	service := &model.Service{Name: a.App, BuildVersion: a.BuildVer, Version: a.Ver}
	err = r.s2.Push(c, service)
	return
}

//SetTokenV4 update Token
func (r *RPC) SetTokenV4(c context.Context, a *model.ArgToken, res *struct{}) (err error) {
	r.s2.SetToken(a.App, a.Token)
	return
}

//Hosts get host list.
func (r *RPC) Hosts(c context.Context, svr string, res *[]*model.Host) (err error) {
	*res, err = r.s2.Hosts(c, svr)
	return
}

//ClearHost clear host.
func (r *RPC) ClearHost(c context.Context, svr string, res *struct{}) error {
	return r.s2.ClearHost(c, svr)
}

// Force push new host config change to config-service
func (r *RPC) Force(c context.Context, a *model.ArgConf, res *struct{}) (err error) {
	service := &model.Service{Name: a.App, BuildVersion: a.BuildVer, Version: a.Ver}
	err = r.s2.Force(c, service, a.Hosts, a.SType)
	return
}
