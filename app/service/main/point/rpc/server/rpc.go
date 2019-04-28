package server

import (
	"go-common/app/service/main/point/conf"
	"go-common/app/service/main/point/model"
	"go-common/app/service/main/point/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC represent rpc server
type RPC struct {
	s *service.Service
}

// New init rpc.
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

// Auth check connection success.
func (r *RPC) Auth(c context.Context, arg *rpc.Auth, res *struct{}) (err error) {
	return
}

// PointInfo point info.
func (r *RPC) PointInfo(c context.Context, a *model.ArgRPCMid, res *model.PointInfo) (err error) {
	var p *model.PointInfo
	if p, err = r.s.PointInfo(c, a.Mid); err == nil && p != nil {
		*res = *p
	}
	return
}

// ConsumePoint point consume.
func (r *RPC) ConsumePoint(c context.Context, a *model.ArgPointConsume, status *int8) (err error) {
	*status, err = r.s.ConsumePoint(c, a)
	return
}

// PointAddByBp point add by bp.
func (r *RPC) PointAddByBp(c context.Context, arg *model.ArgPointAdd, res *int64) (err error) {
	*res, err = r.s.PointAddByBp(c, arg)
	return
}

// PointAdd point add.
func (r *RPC) PointAdd(c context.Context, a *model.ArgPoint, status *int8) (err error) {
	*status, err = r.s.AddPoint(c, a)
	return
}

// PointHistory point history.
func (r *RPC) PointHistory(c context.Context, arg *model.ArgRPCPointHistory, res *model.PointHistoryResp) (err error) {
	var (
		phs   []*model.OldPointHistory
		total int
	)
	if phs, total, err = r.s.OldPointHistory(c, arg.Mid, arg.PN, arg.PS); err == nil && len(phs) > 0 {
		p := &model.PointHistoryResp{
			Phs:   phs,
			Total: total,
		}
		*res = *p
	}
	return
}
