package config

import (
	"context"

	"go-common/app/infra/config/model"
	"go-common/library/net/rpc"
)

const (
	_appid = "config.service"

	_push       = "RPC.Push"
	_setToken   = "RPC.SetToken"
	_pushV4     = "RPC.PushV4"
	_force      = "RPC.Force"
	_setTokenV4 = "RPC.SetTokenV4"
	_hosts      = "RPC.Hosts"
	_clearHost  = "RPC.ClearHost"
)

var (
	_noArg = &struct{}{}
)

//Service2 service.
type Service2 struct {
	client *rpc.Client2
}

// New2 new a config service.
func New2(c *rpc.ClientConfig) (s *Service2) {
	s = &Service2{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// Push push new ver to config-service
func (s *Service2) Push(c context.Context, arg *model.ArgConf) (err error) {
	err = s.client.Boardcast(c, _push, arg, _noArg)
	return
}

// SetToken update token in config-service
func (s *Service2) SetToken(c context.Context, arg *model.ArgToken) (err error) {
	err = s.client.Boardcast(c, _setToken, arg, _noArg)
	return
}

// PushV4 push new ver to config-service
func (s *Service2) PushV4(c context.Context, arg *model.ArgConf) (err error) {
	err = s.client.Boardcast(c, _pushV4, arg, _noArg)
	return
}

// SetTokenV4 update token in config-service
func (s *Service2) SetTokenV4(c context.Context, arg *model.ArgToken) (err error) {
	err = s.client.Boardcast(c, _setTokenV4, arg, _noArg)
	return
}

//Hosts get host list.
func (s *Service2) Hosts(c context.Context, svr string) (hosts []*model.Host, err error) {
	err = s.client.Call(c, _hosts, svr, &hosts)
	return
}

// ClearHost update token in config-service
func (s *Service2) ClearHost(c context.Context, svr string) (err error) {
	err = s.client.Call(c, _clearHost, svr, _noArg)
	return
}

// Force push new host ver to config-service
func (s *Service2) Force(c context.Context, arg *model.ArgConf) (err error) {
	err = s.client.Boardcast(c, _force, arg, _noArg)
	return
}
