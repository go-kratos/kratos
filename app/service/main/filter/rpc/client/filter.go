package filter

import (
	"context"

	rpcmodel "go-common/app/service/main/filter/model/rpc"
	"go-common/library/net/rpc"
)

const (
	_filter  = "RPC.Filter"
	_mfilter = "RPC.MFilter"

	_filterArea  = "RPC.FilterArea"
	_mfilterArea = "RPC.MFilterArea"
)

const (
	_appid = "filter.service"
)

// Service struct .
type Service struct {
	client *rpc.Client2
}

// New .
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// Filter .
func (s *Service) Filter(c context.Context, arg *rpcmodel.ArgFilter) (res *rpcmodel.FilterRes, err error) {
	res = new(rpcmodel.FilterRes)
	err = s.client.Call(c, _filter, arg, res)
	return
}

// MFilter .
func (s *Service) MFilter(c context.Context, arg *rpcmodel.ArgMfilter) (res map[string]*rpcmodel.FilterRes, err error) {
	err = s.client.Call(c, _mfilter, arg, &res)
	return
}

// FilterArea .
func (s *Service) FilterArea(c context.Context, arg *rpcmodel.ArgFilter) (res *rpcmodel.FilterRes, err error) {
	res = new(rpcmodel.FilterRes)
	err = s.client.Call(c, _filterArea, arg, res)
	return
}

// MFilterArea .
func (s *Service) MFilterArea(c context.Context, arg *rpcmodel.ArgMfilter) (res map[string]*rpcmodel.FilterRes, err error) {
	err = s.client.Call(c, _mfilterArea, arg, &res)
	return
}
