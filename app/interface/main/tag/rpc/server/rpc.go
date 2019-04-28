package rpc

import (
	"go-common/app/interface/main/tag/conf"
	"go-common/app/interface/main/tag/model"
	"go-common/app/interface/main/tag/service"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
	"go-common/library/xstr"
)

// RPC represent rpc server
type RPC struct {
	c   *conf.Config
	srv *service.Service
}

// Init init rpc.
func Init(c *conf.Config, s *service.Service) (svr *rpc.Server) {
	r := &RPC{
		c:   c,
		srv: s,
	}
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

// InfoByID Info return a tag info by id
func (r *RPC) InfoByID(c context.Context, arg *model.ArgID, res *model.Tag) (err error) {
	var v *model.Tag
	if v, err = r.srv.InfoByID(c, arg.Mid, arg.ID); err == nil {
		*res = *v
	}
	return
}

// InfoByName Info return a tag info by name.
func (r *RPC) InfoByName(c context.Context, arg *model.ArgName, res *model.Tag) (err error) {
	var v *model.Tag
	if v, err = r.srv.InfoByName(c, arg.Mid, arg.Name); err == nil {
		*res = *v
	}
	return
}

// InfoByIDs return tags by ids
func (r *RPC) InfoByIDs(c context.Context, arg *model.ArgIDs, res *[]*model.Tag) (err error) {
	*res, err = r.srv.MinfoByIDs(c, arg.Mid, arg.IDs)
	return
}

// InfoByNames return tags by names
func (r *RPC) InfoByNames(c context.Context, arg *model.ArgNames, res *[]*model.Tag) (err error) {
	*res, err = r.srv.MinfoByNames(c, arg.Mid, arg.Names)
	return
}

// ArcTags return archive tags by aid and mid.
func (r *RPC) ArcTags(c context.Context, arg *model.ArgAid, res *[]*model.Tag) (err error) {
	*res, err = r.srv.ArcTags(c, arg.Aid, arg.Mid)
	return
}

// SubTags return the user subscribe tags.
func (r *RPC) SubTags(c context.Context, arg *model.ArgSub, res *model.Sub) (err error) {
	tags, total, err := r.srv.SubTags(c, arg.Mid, arg.Vmid, arg.Pn, arg.Ps, arg.Order)
	if err != nil {
		return
	}
	res.Tags = tags
	res.Total = total
	return
}

// AddSub AddSub.
func (r *RPC) AddSub(c context.Context, arg *model.ArgAddSub, res *struct{}) (err error) {
	return r.srv.AddSub(c, arg.Mid, arg.Tids, arg.Now)
}

// CancelSub CancelSub.
func (r *RPC) CancelSub(c context.Context, arg *model.ArgCancelSub, res *struct{}) (err error) {
	return r.srv.CancelSub(c, arg.Tid, arg.Mid, arg.Now)
}

// UpdateCustomSort update custome sort tags.
func (r *RPC) UpdateCustomSort(c context.Context, arg *model.ArgUpdateCustomSort, res *struct{}) (err error) {
	tids, err := xstr.SplitInts(arg.Tids)
	if err != nil {
		log.Error("xstr.SplitInts(%s) error(%v)", arg.Tids, err)
		return
	}
	err = r.srv.UpCustomSubChannels(c, arg.Mid, tids, arg.Type)
	return
}

// CustomSortChannel custome sort channel.
func (r *RPC) CustomSortChannel(c context.Context, arg *model.ArgCustomSort, res *model.CustomSortChannel) (err error) {
	v := new(model.CustomSortChannel)
	v.Custom, v.Standard, v.Total, err = r.srv.CustomSubTags(c, arg.Mid, arg.Order, arg.Type, arg.Ps, arg.Pn)
	if err == nil {
		*res = *v
	}
	return
}

// TagTop web-interface tag top, include tag info, and similar tags.
func (r *RPC) TagTop(c context.Context, arg *model.ReqTagTop, res *model.TagTop) (err error) {
	if arg.TName != "" {
		if arg.TName, err = r.srv.CheckName(arg.TName); err != nil {
			return
		}
	}
	if arg.Tid <= 0 && arg.TName == "" {
		return ecode.RequestErr
	}
	var v *model.TagTop
	if v, err = r.srv.TagTop(c, arg); err == nil {
		*res = *v
	}
	return
}
