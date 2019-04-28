package rpc

import (
	"go-common/app/interface/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/net/rpc/context"
)

// ChannelResources get channel resources.
func (r *RPC) ChannelResources(c context.Context, arg *model.ArgChannelResource, res *model.ChannelResource) (err error) {
	if arg.Tid <= 0 && arg.Name == "" {
		return ecode.TagNotExist
	}
	if arg.From != model.ChannelFromH5 {
		arg.From = model.ChannelFromApp
	}
	var v *model.ChannelResource
	if v, err = r.srv.ChannelResources(c, arg); err != nil {
		return
	}
	*res = *v
	return
}

// ChannelCategory get all channel category.
func (r *RPC) ChannelCategory(c context.Context, arg *struct{}, res *[]*model.ChannelCategory) (err error) {
	*res, err = r.srv.ChannelCategory(c)
	return
}

// ChannelCategories get all channel category.
func (r *RPC) ChannelCategories(c context.Context, arg *model.ArgChannelCategories, res *[]*model.ChannelCategory) (err error) {
	*res, err = r.srv.ChannelCategories(c, arg)
	return
}

// ChanneList channel list.
func (r *RPC) ChanneList(c context.Context, arg *model.ArgChanneList, res *[]*model.Channel) (err error) {
	*res, err = r.srv.ChanneList(c, arg.Mid, arg.ID, arg.From)
	return
}

// DiscoverChannel discover channel.
func (r *RPC) DiscoverChannel(c context.Context, arg *model.ArgDiscoverChanneList, res *[]*model.Channel) (err error) {
	*res, err = r.srv.DiscoveryChannel(c, arg.Mid, arg.From)
	return
}

// ChannelSquare channel square.
func (r *RPC) ChannelSquare(c context.Context, arg *model.ReqChannelSquare, res *model.ChannelSquare) (err error) {
	if arg.TagNumber <= 0 {
		return ecode.TagNotExist
	}
	var v *model.ChannelSquare
	if v, err = r.srv.ChannelSquare(c, arg); err != nil {
		return
	}
	*res = *v
	return
}

// RecommandChannel RecommandChannel.
func (r *RPC) RecommandChannel(c context.Context, arg *model.ArgRecommandChannel, res *[]*model.Channel) (err error) {
	*res, err = r.srv.RecommandChannel(c, arg.Mid, arg.From)
	return
}

// ResChannelCheckBack .
func (r *RPC) ResChannelCheckBack(c context.Context, arg *model.ArgResChannel, res *map[int64]*model.ResChannelCheckBack) (err error) {
	*res, err = r.srv.ResChannelCheckBack(c, arg.Oids, arg.Type)
	return
}

// ChannelDetail channel detail.
func (r *RPC) ChannelDetail(c context.Context, arg *model.ReqChannelDetail, res *model.ChannelDetail) (err error) {
	if arg.Tid <= 0 && arg.TName == "" {
		return ecode.RequestErr
	}
	var v *model.ChannelDetail
	if v, err = r.srv.ChannelDetail(c, arg); err != nil {
		return
	}
	*res = *v
	return
}
