package rpc

import (
	cmmdl "go-common/app/service/main/spy/model"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"

	"go-common/app/service/main/spy/conf"
	"go-common/app/service/main/spy/model"
	"go-common/app/service/main/spy/service"
)

// RPC server def.
type RPC struct {
	s *service.Service
}

// New create instance of rpc server2 and return.
func New(c *conf.Config, s *service.Service) (svr *rpc.Server) {
	r := &RPC{s: s}
	svr = rpc.NewServer(c.RPCServer)
	if err := svr.Register(r); err != nil {
		panic(err)
	}
	return
}

// Ping check rpc server health.
func (r *RPC) Ping(c context.Context, arg *struct{}, res *struct{}) (err error) {
	return
}

// UserScore rpc req for getting user spy score.
func (r *RPC) UserScore(c context.Context, arg *cmmdl.ArgUserScore, res *cmmdl.UserScore) (err error) {
	var (
		ui *model.UserInfo
	)
	if ui, err = r.s.UserInfo(c, arg.Mid, arg.IP); ui != nil && err == nil {
		*res = cmmdl.UserScore{Mid: ui.Mid, Score: ui.Score}
	}
	return
}

// HandleEvent rpc req for handling spy event , maybe from spy-job.
func (r *RPC) HandleEvent(c context.Context, arg *cmmdl.ArgHandleEvent, res *struct{}) (err error) {
	var event = &model.EventMessage{
		Time:      arg.Time,
		IP:        arg.IP,
		Service:   arg.Service,
		Event:     arg.Event,
		ActiveMid: arg.ActiveMid,
		TargetMid: arg.TargetMid,
		TargetID:  arg.TargetID,
		Args:      arg.Args,
		Result:    arg.Result,
		Effect:    arg.Effect,
		RiskLevel: arg.RiskLevel,
	}
	return r.s.HandleEvent(c, event)
}

// ReBuildPortrait rpc.
func (r *RPC) ReBuildPortrait(c context.Context, arg *cmmdl.ArgReBuild, res *struct{}) (err error) {
	return r.s.ReBuildPortrait(c, arg.Mid, arg.Reason)
}

// UpdateBaseScore rpc.
func (r *RPC) UpdateBaseScore(c context.Context, arg *cmmdl.ArgReset, res *struct{}) (err error) {
	return r.s.UpdateBaseScore(c, arg)
}

// UpdateEventScore rpc.
func (r *RPC) UpdateEventScore(c context.Context, arg *cmmdl.ArgReset, res *struct{}) (err error) {
	return r.s.UpdateEventScore(c, arg)
}

// UserInfo rpc.
func (r *RPC) UserInfo(c context.Context, arg *cmmdl.ArgUser, res *cmmdl.UserInfo) (err error) {
	var (
		ui *model.UserInfo
	)
	if ui, err = r.s.UserInfo(c, arg.Mid, arg.IP); ui != nil && err == nil {
		*res = *ui
	}
	return
}

// ClearReliveTimes rpc.
func (r *RPC) ClearReliveTimes(c context.Context, arg *cmmdl.ArgReset, res *struct{}) (err error) {
	return r.s.ClearReliveTimes(c, arg)
}

// StatByID rpc.
func (r *RPC) StatByID(c context.Context, arg *cmmdl.ArgStat, res *[]*cmmdl.Statistics) (err error) {
	var (
		stat []*model.Statistics
	)
	if stat, err = r.s.StatByIDGroupEvent(c, arg.Mid, arg.ID); stat != nil && err == nil {
		*res = stat
	}
	return
}

// RefreshBaseScore rpc.
func (r *RPC) RefreshBaseScore(c context.Context, arg *cmmdl.ArgReset, res *struct{}) (err error) {
	return r.s.RefreshBaseScore(c, arg)
}
