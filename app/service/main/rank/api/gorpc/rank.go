package rank

import (
	"context"

	"go-common/app/service/main/rank/model"
	"go-common/library/net/rpc"
)

const (
	_mget  = "RPC.Mget"
	_sort  = "RPC.Sort"
	_group = "RPC.Group"
)

const (
	_appid = "main.search.rank-service"
)

// Service .
type Service struct {
	client *rpc.Client2
}

// New .
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// Mget .
func (s *Service) Mget(c context.Context, arg *model.MgetReq) (res *model.MgetResp, err error) {
	err = s.client.Call(c, _mget, arg, &res)
	return
}

// Sort .
func (s *Service) Sort(c context.Context, arg *model.SortReq) (res *model.SortResp, err error) {
	err = s.client.Call(c, _sort, arg, &res)
	return
}

// Group .
func (s *Service) Group(c context.Context, arg *model.GroupReq) (res *model.GroupResp, err error) {
	err = s.client.Call(c, _group, arg, &res)
	return
}
