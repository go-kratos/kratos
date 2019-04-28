package tag

import (
	"context"

	"go-common/app/interface/main/tag/model"
)

const (
	_channelResources    = "RPC.ChannelResources"
	_channelCategory     = "RPC.ChannelCategory"
	_channelCategories   = "RPC.ChannelCategories"
	_channeList          = "RPC.ChanneList"
	_discoverChannel     = "RPC.DiscoverChannel"
	_channelSquare       = "RPC.ChannelSquare"
	_recommandChannel    = "RPC.RecommandChannel"
	_resChannelCheckBack = "RPC.ResChannelCheckBack"
	_channelDetail       = "RPC.ChannelDetail"
)

var (
	_noArg = &struct{}{}
)

// ChannelResources .
func (s *Service) ChannelResources(c context.Context, arg *model.ArgChannelResource) (res *model.ChannelResource, err error) {
	err = s.client.Call(c, _channelResources, arg, &res)
	return
}

// ChannelCategory .
func (s *Service) ChannelCategory(c context.Context) (res []*model.ChannelCategory, err error) {
	err = s.client.Call(c, _channelCategory, _noArg, &res)
	return
}

// ChannelCategories .
func (s *Service) ChannelCategories(c context.Context, arg *model.ArgChannelCategories) (res []*model.ChannelCategory, err error) {
	err = s.client.Call(c, _channelCategories, arg, &res)
	return
}

// ChanneList .
func (s *Service) ChanneList(c context.Context, arg *model.ArgChanneList) (res []*model.Channel, err error) {
	err = s.client.Call(c, _channeList, arg, &res)
	return
}

// DiscoverChannel .
func (s *Service) DiscoverChannel(c context.Context, arg *model.ArgDiscoverChanneList) (res []*model.Channel, err error) {
	err = s.client.Call(c, _discoverChannel, arg, &res)
	return
}

// ChannelSquare .
func (s *Service) ChannelSquare(c context.Context, arg *model.ReqChannelSquare) (res *model.ChannelSquare, err error) {
	err = s.client.Call(c, _channelSquare, arg, &res)
	return
}

// RecommandChannel .
func (s *Service) RecommandChannel(c context.Context, arg *model.ArgRecommandChannel) (res []*model.Channel, err error) {
	err = s.client.Call(c, _recommandChannel, arg, &res)
	return
}

// ResChannelCheckBack .
func (s *Service) ResChannelCheckBack(c context.Context, arg *model.ArgResChannel) (res map[int64]*model.ResChannelCheckBack, err error) {
	err = s.client.Call(c, _resChannelCheckBack, arg, &res)
	return
}

// ChannelDetail .
func (s *Service) ChannelDetail(c context.Context, arg *model.ReqChannelDetail) (res *model.ChannelDetail, err error) {
	err = s.client.Call(c, _channelDetail, arg, &res)
	return
}
