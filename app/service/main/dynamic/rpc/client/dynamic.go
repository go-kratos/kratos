package client

import (
	"context"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/dynamic/model"
	"go-common/library/net/rpc"
)

const (
	_regionArcs3    = "RPC.RegionArcs3"
	_regionTagArcs3 = "RPC.RegionTagArcs3"
	_regionsArcs3   = "RPC.RegionsArcs3"
	_regionTotal    = "RPC.RegionTotal"
)

const (
	_appid = "archive.service.dynamic"
)

// Service struct info.
type Service struct {
	client *rpc.Client2
}

// New new service instance and return.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// RegionTotal receive real ip,then return dynamic region total.
func (s *Service) RegionTotal(c context.Context, arg *model.ArgRegionTotal) (res map[string]int, err error) {
	err = s.client.Call(c, _regionTotal, arg, &res)
	return
}

// RegionArcs3 receive ArgRegion contains regionId and real ip, then return dynamic archives.
func (s *Service) RegionArcs3(c context.Context, arg *model.ArgRegion3) (res *model.DynamicArcs3, err error) {
	res = new(model.DynamicArcs3)
	err = s.client.Call(c, _regionArcs3, arg, res)
	return
}

// RegionTagArcs3 receive ArgRegionTag contains tagId and regionId and real ip, then return dynamic archives3.
func (s *Service) RegionTagArcs3(c context.Context, arg *model.ArgRegionTag3) (res *model.DynamicArcs3, err error) {
	res = new(model.DynamicArcs3)
	err = s.client.Call(c, _regionTagArcs3, arg, res)
	return
}

// RegionsArcs3 receive ArgRegion contains regionIds and real ip, then return dynamic archives3.
func (s *Service) RegionsArcs3(c context.Context, arg *model.ArgRegions3) (res map[int32][]*api.Arc, err error) {
	err = s.client.Call(c, _regionsArcs3, arg, &res)
	return
}
