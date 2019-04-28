package client

import (
	"context"

	"go-common/app/service/main/location/model"
	"go-common/library/net/rpc"
)

const (
	_archive  = "RPC.Archive"
	_archive2 = "RPC.Archive2"
	_group    = "RPC.Group"
	_authPIDs = "RPC.AuthPIDs"
	// new
	_info          = "RPC.Info"
	_infos         = "RPC.Infos"
	_infoComplete  = "RPC.InfoComplete"
	_infosComplete = "RPC.InfosComplete"
	// app id
	_appid = "location.service"
)

// Service is resource rpc client.
type Service struct {
	client *rpc.Client2
}

// New new a resource rpc client.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// Archive get the aid auth.
func (s *Service) Archive(c context.Context, arg *model.Archive) (res *int64, err error) {
	res = new(int64)
	err = s.client.Call(c, _archive, arg, res)
	return
}

// Archive2 get the aid auth.
func (s *Service) Archive2(c context.Context, arg *model.Archive) (res *model.Auth, err error) {
	res = new(model.Auth)
	err = s.client.Call(c, _archive2, arg, res)
	return
}

// Group get the gip auth.
func (s *Service) Group(c context.Context, arg *model.Group) (res *model.Auth, err error) {
	res = new(model.Auth)
	err = s.client.Call(c, _group, arg, res)
	return
}

// AuthPIDs check if ip in pids.
func (s *Service) AuthPIDs(c context.Context, arg *model.ArgPids) (res map[int64]*model.Auth, err error) {
	err = s.client.Call(c, _authPIDs, arg, &res)
	return
}

// Info get the ip info.
func (s *Service) Info(c context.Context, arg *model.ArgIP) (res *model.Info, err error) {
	res = new(model.Info)
	err = s.client.Call(c, _info, arg, res)
	return
}

// Infos get the ips info.
func (s *Service) Infos(c context.Context, arg []string) (res map[string]*model.Info, err error) {
	err = s.client.Call(c, _infos, arg, &res)
	return
}

// InfoComplete get the whold ip info.
func (s *Service) InfoComplete(c context.Context, arg *model.ArgIP) (res *model.InfoComplete, err error) {
	res = new(model.InfoComplete)
	err = s.client.Call(c, _infoComplete, arg, res)
	return
}

// InfosComplete get the whold ips infos.
func (s *Service) InfosComplete(c context.Context, arg []string) (res map[string]*model.InfoComplete, err error) {
	err = s.client.Call(c, _infosComplete, arg, &res)
	return
}
