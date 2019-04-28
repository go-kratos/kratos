package server

import (
	"go-common/app/interface/main/dm2/conf"
	"go-common/app/interface/main/dm2/model"
	"go-common/app/interface/main/dm2/service"
	"go-common/library/ecode"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC rpc server.
type RPC struct {
	s *service.Service
}

// New new a rpc server.
func New(c *conf.Config, s *service.Service) (svr *rpc.Server) {
	r := &RPC{s: s}
	svr = rpc.NewServer(c.RPCServer)
	if err := svr.Register(r); err != nil {
		panic(err)
	}
	return
}

// Auth check connection success.
func (r *RPC) Auth(c context.Context, arg *rpc.Auth, res *struct{}) (err error) {
	return
}

// Ping checks connection success.
func (r *RPC) Ping(c context.Context, arg *struct{}, res *struct{}) (err error) {
	return
}

// SubjectInfos multi get dm subject info by type and oids.
func (r *RPC) SubjectInfos(c context.Context, a *model.ArgOids, res *map[int64]*model.SubjectInfo) (err error) {
	*res, err = r.s.SubjectInfos(c, a.Type, a.Plat, a.Oids)
	return
}

// EditDMState edit dm state.
func (r *RPC) EditDMState(c context.Context, a *model.ArgEditDMState, res *struct{}) (err error) {
	err = r.s.EditDMState(c, a.Type, a.Mid, a.Oid, a.State, a.Dmids, a.Source, a.OperatorType)
	return
}

// EditDMPool edit dm pool.
func (r *RPC) EditDMPool(c context.Context, a *model.ArgEditDMPool, res *struct{}) (err error) {
	err = r.s.EditDMPool(c, a.Type, a.Mid, a.Oid, a.Pool, a.Dmids, a.Source, a.OperatorType)
	return
}

// EditDMAttr edit dm attr
func (r *RPC) EditDMAttr(c context.Context, a *model.ArgEditDMAttr, res *struct{}) (err error) {
	_, err = r.s.EditDMAttr(c, a.Type, a.Mid, a.Oid, a.Bit, a.Value, a.Dmids, a.Source, a.OperatorType)
	return
}

// BuyAdvance 购买高级弹幕
func (r *RPC) BuyAdvance(c context.Context, a *model.ArgAdvance, res *struct{}) (err error) {
	err = r.s.BuyAdvance(c, a.Mid, a.Cid, a.Mode)
	return
}

// AdvanceState 高级弹幕状态
func (r *RPC) AdvanceState(c context.Context, a *model.ArgAdvance, res *model.AdvState) (err error) {
	var v *model.AdvState
	if v, err = r.s.AdvanceState(c, a.Mid, a.Cid, a.Mode); err == nil {
		*res = *v
	}
	return
}

// Advances 高级弹幕申请列表
func (r *RPC) Advances(c context.Context, a *model.ArgMid, res *[]*model.Advance) (err error) {
	*res, err = r.s.Advances(c, a.Mid)
	return
}

// PassAdvance 通过高级弹幕申请
func (r *RPC) PassAdvance(c context.Context, a *model.ArgUpAdvance, res *struct{}) (err error) {
	err = r.s.PassAdvance(c, a.Mid, a.ID)
	return
}

// DenyAdvance 拒绝高级弹幕申请
func (r *RPC) DenyAdvance(c context.Context, a *model.ArgUpAdvance, res *struct{}) (err error) {
	err = r.s.DenyAdvance(c, a.Mid, a.ID)
	return
}

// CancelAdvance 取消高级弹幕申请
func (r *RPC) CancelAdvance(c context.Context, a *model.ArgUpAdvance, res *struct{}) (err error) {
	err = r.s.CancelAdvance(c, a.Mid, a.ID)
	return
}

// AddUserFilters add user filter.
func (r *RPC) AddUserFilters(c context.Context, a *model.ArgAddUserFilters, res *[]*model.UserFilter) (err error) {
	fltMap := make(map[string]string)
	for _, filter := range a.Filters {
		fltMap[filter] = a.Comment
	}
	*res, err = r.s.AddUserFilters(c, a.Mid, a.Type, fltMap)
	return
}

// UserFilters multi get user filters.
func (r *RPC) UserFilters(c context.Context, a *model.ArgMid, res *[]*model.UserFilter) (err error) {
	*res, err = r.s.UserFilters(c, a.Mid)
	return
}

// DelUserFilters delete user filters by filter id.
func (r *RPC) DelUserFilters(c context.Context, a *model.ArgDelUserFilters, affect *int64) (err error) {
	*affect, err = r.s.DelUserFilters(c, a.Mid, a.IDs)
	return
}

// AddUpFilters add up filters.
func (r *RPC) AddUpFilters(c context.Context, a *model.ArgAddUpFilters, res *struct{}) (err error) {
	fltMap := make(map[string]string)
	for _, filter := range a.Filters {
		fltMap[filter] = "" // NOTE here should be comment assignment
	}
	err = r.s.AddUpFilters(c, a.Mid, a.Type, fltMap)
	return
}

// UpFilters multi get up filters.
func (r *RPC) UpFilters(c context.Context, a *model.ArgUpFilters, res *[]*model.UpFilter) (err error) {
	*res, err = r.s.UpFilters(c, a.Mid)
	return
}

// BanUsers ban user by upper or assist.
func (r *RPC) BanUsers(c context.Context, a *model.ArgBanUsers, res *struct{}) (err error) {
	err = r.s.BanUsers(c, a.Mid, a.Oid, a.DMIDs)
	return
}

// CancelBanUsers cancel users by upper or assiat.
func (r *RPC) CancelBanUsers(c context.Context, a *model.ArgCancelBanUsers, res *struct{}) (err error) {
	err = r.s.CancelBanUsers(c, a.Mid, a.Aid, a.Filters)
	return
}

// EditUpFilters edit upper filters.
func (r *RPC) EditUpFilters(c context.Context, a *model.ArgEditUpFilters, res *int64) (err error) {
	*res, err = r.s.EditUpFilters(c, a.Mid, a.Type, a.Active, a.Filters)
	return
}

// AddGlobalFilter add global filters.
func (r *RPC) AddGlobalFilter(c context.Context, a *model.ArgAddGlobalFilter, res *model.GlobalFilter) (err error) {
	var v *model.GlobalFilter
	if v, err = r.s.AddGlobalFilter(c, a.Type, a.Filter); err == nil {
		*res = *v
	}
	return
}

// GlobalFilters multi get global filters.
func (r *RPC) GlobalFilters(c context.Context, a *model.ArgGlobalFilters, res *[]*model.GlobalFilter) (err error) {
	if a == nil {
		err = ecode.RequestErr
		return
	}
	*res, err = r.s.GlobalFilters(c)
	return
}

// DelGlobalFilters delete global filter.
func (r *RPC) DelGlobalFilters(c context.Context, a *model.ArgDelGlobalFilters, affect *int64) (err error) {
	*affect, err = r.s.DelGlobalFilters(c, a.IDs)
	return
}

// Mask get web mask.
func (r *RPC) Mask(c context.Context, a *model.ArgMask, res *model.Mask) (err error) {
	if a == nil {
		err = ecode.RequestErr
		return
	}
	var m *model.Mask
	if m, err = r.s.MaskList(c, a.Cid, a.Plat); err == nil && m != nil {
		*res = *m
	}
	return
}

// SubtitleGet .
func (r *RPC) SubtitleGet(c context.Context, arg *model.ArgSubtitleGet, res *model.VideoSubtitles) (err error) {
	var v *model.VideoSubtitles
	if v, err = r.s.GetWebVideoSubtitle(c, arg.Aid, arg.Oid, arg.Type); err == nil {
		*res = *v
	}
	return
}

// SubtitleSujectSubmit set archive allow submit
func (r *RPC) SubtitleSujectSubmit(c context.Context, arg *model.ArgSubtitleAllowSubmit, res *struct{}) (err error) {
	err = r.s.SubtitleSubjectSubmit(c, arg.Aid, arg.AllowSubmit, arg.Lan)
	return
}

// SubtitleSubjectSubmitGet get archive allow submit
func (r *RPC) SubtitleSubjectSubmitGet(c context.Context, arg *model.ArgArchiveID, res *model.SubtitleSubjectReply) (err error) {
	var (
		reply *model.SubtitleSubjectReply
	)
	reply, err = r.s.SubtitleSubject(c, arg.Aid)
	*res = *reply
	return
}
