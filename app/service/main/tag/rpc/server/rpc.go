package server

import (
	"go-common/app/service/main/tag/conf"
	"go-common/app/service/main/tag/model"
	"go-common/app/service/main/tag/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC .
type RPC struct {
	conf *conf.Config
	svr  *service.Service
}

// Init .
func Init(c *conf.Config, s *service.Service) (svr *rpc.Server) {
	r := &RPC{
		conf: c,
		svr:  s,
	}
	svr = rpc.NewServer(c.RPCServer)
	if err := svr.Register(r); err != nil {
		panic(err)
	}
	return
}

// Ping .
func (r *RPC) Ping(c context.Context, arg *struct{}, res *struct{}) (err error) {
	return
}

// InfoByID .
func (r *RPC) InfoByID(c context.Context, arg *model.ArgID, res *model.Tag) (err error) {
	var v *model.Tag
	if v, err = r.svr.Info(c, arg.Mid, arg.ID); err == nil {
		*res = *v
	}
	return
}

// InfoByName .
func (r *RPC) InfoByName(c context.Context, arg *model.ArgName, res *model.Tag) (err error) {
	var v *model.Tag
	if v, err = r.svr.InfoByName(c, arg.Mid, arg.Name); err == nil {
		*res = *v
	}
	return
}

// CheckName .
func (r *RPC) CheckName(c context.Context, arg *model.ArgCheckName, res *model.Tag) (err error) {
	var v *model.Tag
	if v, err = r.svr.CheckTag(c, arg.Name, int32(arg.Type), arg.Now, arg.RealIP); err == nil {
		*res = *v
	}
	return
}

// Count .
func (r *RPC) Count(c context.Context, arg *model.ArgID, res *model.Count) (err error) {
	if v, err := r.svr.Count(c, arg.ID); err == nil {
		*res = *v
	}
	return
}

// Counts . res *map[int64]*model.Count
func (r *RPC) Counts(c context.Context, arg *model.ArgIDs, res *map[int64]*model.Count) (err error) {
	if v, err := r.svr.Counts(c, arg.IDs); err == nil {
		*res = v
	}
	return
}

// InfoByIDs .
func (r *RPC) InfoByIDs(c context.Context, arg *model.ArgIDs, res *[]*model.Tag) (err error) {
	*res, err = r.svr.Infos(c, arg.Mid, arg.IDs)
	return
}

// InfoByNames .
func (r *RPC) InfoByNames(c context.Context, arg *model.ArgNames, res *[]*model.Tag) (err error) {
	*res, err = r.svr.InfosByNames(c, arg.Mid, arg.Names)
	return
}

// ResTags .
func (r *RPC) ResTags(c context.Context, arg *model.ArgResTags, res *[]*model.Resource) (err error) {
	*res, err = r.svr.ResTags(c, arg.Oid, arg.Type, arg.Mid)
	return
}

// ResTagLog .
func (r *RPC) ResTagLog(c context.Context, arg *model.ArgResTagLog, res *[]*model.ResourceLog) (err error) {
	*res, err = r.svr.ResTagLog(c, arg.Oid, arg.Type, arg.Mid, arg.Pn, arg.Ps, arg.RealIP)
	return
}

// AddSub .
func (r *RPC) AddSub(c context.Context, arg *model.ArgAddSub, res *struct{}) (err error) {
	err = r.svr.AddSub(c, arg.Mid, arg.Tids, arg.RealIP)
	return
}

// CancelSub .
func (r *RPC) CancelSub(c context.Context, arg *model.ArgCancelSub, res *struct{}) (err error) {
	err = r.svr.CancelSub(c, arg.Mid, arg.Tid, arg.RealIP)
	return
}

// SubTags .
func (r *RPC) SubTags(c context.Context, arg *model.ArgSub, res *model.ResSub) (err error) {
	v := new(model.ResSub)
	v.Tags, v.Total, err = r.svr.SubTags(c, arg.Mid, arg.Pn, arg.Ps, arg.Order)
	if err == nil {
		*res = *v
	}
	return
}

// AddCustomSubTag .
func (r *RPC) AddCustomSubTag(c context.Context, arg *model.ArgCustomSub, res *struct{}) (err error) {
	err = r.svr.AddCustomSubTags(c, arg.Mid, arg.Type, arg.Tids, arg.RealIP)
	return
}

// CustomSubTag .
func (r *RPC) CustomSubTag(c context.Context, arg *model.ArgSub, res *model.ResSubSort) (err error) {
	v := new(model.ResSubSort)
	v.Sort, v.Tags, v.Total, err = r.svr.CustomSubTags(c, arg.Mid, arg.Type, arg.Pn, arg.Ps, arg.Order)
	if err == nil {
		*res = *v
	}
	return
}

// AddCustomSubChannel .
func (r *RPC) AddCustomSubChannel(c context.Context, arg *model.ArgCustomSub, res *struct{}) (err error) {
	err = r.svr.AddCustomSubChannels(c, arg.Mid, arg.Type, arg.Tids, arg.RealIP)
	return
}

// CustomSubChannel .
func (r *RPC) CustomSubChannel(c context.Context, arg *model.ArgSub, res *model.ResSubSort) (err error) {
	v := new(model.ResSubSort)
	v.Sort, v.Tags, v.Total, err = r.svr.CustomSubTags(c, arg.Mid, arg.Type, arg.Pn, arg.Ps, arg.Order)
	if err == nil {
		*res = *v
	}
	return
}

// Like .
func (r *RPC) Like(c context.Context, arg *model.ArgResAction, res *struct{}) (err error) {
	err = r.svr.Like(c, arg.Mid, arg.Oid, arg.Tid, arg.Type, arg.RealIP)
	return
}

// Hate .
func (r *RPC) Hate(c context.Context, arg *model.ArgResAction, res *struct{}) (err error) {
	err = r.svr.Hate(c, arg.Mid, arg.Oid, arg.Tid, arg.Type, arg.RealIP)
	return
}

// ResAction .
func (r *RPC) ResAction(c context.Context, arg *model.ArgResAction, res *int32) (err error) {
	*res, err = r.svr.Action(c, arg.Mid, arg.Oid, arg.Tid, arg.Type)
	return
}

// ResActionMap .
func (r *RPC) ResActionMap(c context.Context, arg *model.ArgResActions, res *map[int64]int32) (err error) {
	v, err := r.svr.ActionMap(c, arg.Mid, arg.Oid, arg.Type, arg.Tids)
	if err == nil {
		*res = v
	}
	return
}

// HideTag .
func (r *RPC) HideTag(c context.Context, arg *model.ArgHide, res *struct{}) (err error) {
	err = r.svr.HideTag(c, arg.Tid, arg.State)
	return
}

// CreateTag .
func (r *RPC) CreateTag(c context.Context, arg *model.ArgCreate, res *struct{}) (err error) {
	err = r.svr.CreateTag(c, arg.Tag)
	return
}

// CreateTags .
func (r *RPC) CreateTags(c context.Context, arg *model.ArgCreate, res *struct{}) (err error) {
	err = r.svr.CreateTags(c, arg.Tags)
	return
}

// PlatformUpBind .
func (r *RPC) PlatformUpBind(c context.Context, arg model.ArgUPBind, res *struct{}) (err error) {
	err = r.svr.PlatformUpBind(c, arg.Oid, arg.Mid, arg.Tids, arg.Type, arg.RealIP)
	return
}

// PlatformAdminBind .
func (r *RPC) PlatformAdminBind(c context.Context, arg model.ArgUPBind, res *struct{}) (err error) {
	err = r.svr.PlatformAdminBind(c, arg.Oid, arg.Mid, arg.Tids, arg.Type, arg.RealIP)
	return
}

// PlatformUserBind .
func (r *RPC) PlatformUserBind(c context.Context, arg model.ArgUserBind, res *struct{}) (err error) {
	err = r.svr.PlatformUserBind(c, arg.Oid, arg.Mid, arg.Tid, arg.Type, arg.Role, arg.Action, arg.RealIP)
	return
}

// ReportAction .
func (r *RPC) ReportAction(c context.Context, arg model.ArgReportAction, res *struct{}) (err error) {
	err = r.svr.ReportAction(c, arg.Oid, arg.LogID, arg.Mid, arg.Type, arg.PartID, arg.Reason, arg.Score, arg.Content, arg.RealIP)
	return
}

// LimitResource .
func (r *RPC) LimitResource(c context.Context, arg *struct{}, res *[]*model.ResourceLimit) (err error) {
	*res, err = r.svr.LimitResource(c)
	return
}

// WhiteUser .
func (r *RPC) WhiteUser(c context.Context, arg *struct{}, res *map[int64]struct{}) (err error) {
	if v, err := r.svr.WhiteUser(c); err == nil {
		*res = v
	}
	return
}

// TagGroup .
func (r *RPC) TagGroup(c context.Context, arg *struct{}, res *[]*model.Synonym) (err error) {
	if v, err := r.svr.TagGroup(c); err == nil {
		*res = v
	}
	return
}

// ResOidsByTid .
func (r *RPC) ResOidsByTid(c context.Context, arg model.ArgRes, res *[]int64) (err error) {
	if v, err := r.svr.ResOidsByTid(c, arg.Tid, arg.Limit, arg.Type, arg.RealIP); err == nil {
		*res = v
	}
	return
}

// RecommandTag .
func (r *RPC) RecommandTag(c context.Context, arg *struct{}, res *map[int64]map[string][]*model.UploadTag) (err error) {
	if v, err := r.svr.RecommandTag(c); err == nil {
		*res = v
	}
	return
}

// RankingHot .
func (r *RPC) RankingHot(c context.Context, arg *struct{}, res *[]*model.Tag) (err error) {
	*res, err = r.svr.RankingHot(c)
	return
}

// RankingBangumi .
func (r *RPC) RankingBangumi(c context.Context, arg *struct{}, res *model.ResBangumi) (err error) {
	v := new(model.ResBangumi)
	v.Sids, v.Bangumi, err = r.svr.RankingBangumi(c)
	if err == nil {
		*res = *v
	}
	return
}

// RankingRegion .
func (r *RPC) RankingRegion(c context.Context, arg *model.ArgRankingRegion, res *[]*model.RankingRegion) (err error) {
	*res, err = r.svr.RankingRegion(c, arg.Rid)
	return
}

// Hots .
func (r *RPC) Hots(c context.Context, arg *model.ArgHots, res *[]*model.HotTag) (err error) {
	*res, err = r.svr.Hots(c, arg.Rid, arg.Type)
	return
}

// HotMap .
func (r *RPC) HotMap(c context.Context, arg *struct{}, res *map[int16][]int64) (err error) {
	if v, err := r.svr.HotMap(c); err == nil {
		*res = v
	}
	return
}

// Prids .
func (r *RPC) Prids(c context.Context, arg *struct{}, res *[]int16) (err error) {
	if v, err := r.svr.Prids(c); err == nil {
		*res = v
	}
	return
}

// Rids .
func (r *RPC) Rids(c context.Context, arg *struct{}, res *map[int64]int64) (err error) {
	if v, err := r.svr.Rids(c); err == nil {
		*res = v
	}
	return
}

// DefaultUpBind .
func (r *RPC) DefaultUpBind(c context.Context, arg model.ArgDefaultBind, res *struct{}) (err error) {
	err = r.svr.DefaultUpBind(c, arg.Oid, arg.Mid, arg.Tids, arg.Type, arg.RealIP)
	return
}

// DefaultAdminBind .
func (r *RPC) DefaultAdminBind(c context.Context, arg model.ArgDefaultBind, res *struct{}) (err error) {
	err = r.svr.DefaultAdminBind(c, arg.Oid, arg.Mid, arg.Tids, arg.Type, arg.RealIP)
	return
}
