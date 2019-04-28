package server

import (
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/dynamic/conf"
	"go-common/app/service/main/dynamic/model"
	"go-common/app/service/main/dynamic/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC struct info.
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

//RegionTotal return dynamic region total.
func (r *RPC) RegionTotal(c context.Context, a *model.ArgRegionTotal, res *map[string]int) (err error) {
	*res = r.s.RegionTotal(c)
	return
}

// RegionArcs3 receive aid, then init archive3 info.
func (r *RPC) RegionArcs3(c context.Context, a *model.ArgRegion3, res *model.DynamicArcs3) (err error) {
	var (
		count int
		arcs  []*api.Arc
	)
	if arcs, count, err = r.s.RegionArcs3(c, a.RegionID, a.Pn, a.Ps); err == nil {
		res.Page = &model.Page{Num: a.Pn, Size: a.Ps, Count: count}
		res.Archives = arcs
	}
	return
}

// RegionTagArcs3 receive aid, then init archive info.
func (r *RPC) RegionTagArcs3(c context.Context, a *model.ArgRegionTag3, res *model.DynamicArcs3) (err error) {
	var (
		count int
		arcs  []*api.Arc
	)
	if arcs, count, err = r.s.RegionTagArcs3(c, a.RegionID, a.TagID, a.Pn, a.Ps); err == nil {
		res.Page = &model.Page{Num: a.Pn, Size: a.Ps, Count: count}
		res.Archives = arcs
	}
	return
}

// RegionsArcs3 receive rids and return dynamic archives3.
func (r *RPC) RegionsArcs3(c context.Context, a *model.ArgRegions3, res *map[int32][]*api.Arc) (err error) {
	*res, err = r.s.RegionsArcs3(c, a.RegionIDs, a.Count)
	return
}
