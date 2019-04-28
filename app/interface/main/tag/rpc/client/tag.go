package tag

import (
	"context"

	"go-common/app/interface/main/tag/model"
	"go-common/library/net/rpc"
)

const (
	_infoByID          = "RPC.InfoByID"
	_infoByName        = "RPC.InfoByName"
	_infoByIDs         = "RPC.InfoByIDs"
	_infoByNames       = "RPC.InfoByNames"
	_arcTags           = "RPC.ArcTags"
	_subTags           = "RPC.SubTags"
	_updateCustomSort  = "RPC.UpdateCustomSort"
	_customSortChannel = "RPC.CustomSortChannel"
	_addSub            = "RPC.AddSub"
	_cancelSub         = "RPC.CancelSub"
	_tagTop            = "RPC.TagTop"
)

// Service .
type Service struct {
	client *rpc.Client2
}

const (
	_appid = "main.community.tag"
)

// New2 .
func New2(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// InfoByID .
func (s *Service) InfoByID(c context.Context, arg *model.ArgID) (res *model.Tag, err error) {
	res = new(model.Tag)
	err = s.client.Call(c, _infoByID, arg, res)
	return
}

// InfoByName .
func (s *Service) InfoByName(c context.Context, arg *model.ArgName) (res *model.Tag, err error) {
	res = new(model.Tag)
	err = s.client.Call(c, _infoByName, arg, res)
	return
}

// InfoByIDs .
func (s *Service) InfoByIDs(c context.Context, arg *model.ArgIDs) (res []*model.Tag, err error) {
	err = s.client.Call(c, _infoByIDs, arg, &res)
	return
}

// InfoByNames .
func (s *Service) InfoByNames(c context.Context, arg *model.ArgNames) (res []*model.Tag, err error) {
	err = s.client.Call(c, _infoByNames, arg, &res)
	return
}

// ArcTags .
func (s *Service) ArcTags(c context.Context, arg *model.ArgAid) (res []*model.Tag, err error) {
	err = s.client.Call(c, _arcTags, arg, &res)
	return
}

// SubTags .
func (s *Service) SubTags(c context.Context, arg *model.ArgSub) (res *model.Sub, err error) {
	res = new(model.Sub)
	err = s.client.Call(c, _subTags, arg, res)
	return
}

// UpdateCustomSortTags custome sort tags.
func (s *Service) UpdateCustomSortTags(c context.Context, arg *model.ArgUpdateCustomSort) (err error) {
	err = s.client.Call(c, _updateCustomSort, arg, _noRes)
	return
}

// CustomSortChannel custome sort channel.
func (s *Service) CustomSortChannel(c context.Context, arg *model.ArgCustomSort) (res *model.CustomSortChannel, err error) {
	err = s.client.Call(c, _customSortChannel, arg, &res)
	return
}

// AddSub .
func (s *Service) AddSub(c context.Context, arg *model.ArgAddSub) (err error) {
	return s.client.Call(c, _addSub, arg, _noRes)
}

// CancelSub .
func (s *Service) CancelSub(c context.Context, arg *model.ArgCancelSub) (err error) {
	return s.client.Call(c, _cancelSub, arg, _noRes)
}

// TagTop web-interface tag top, include tag info, and similar tags.
func (s *Service) TagTop(c context.Context, arg *model.ReqTagTop) (res *model.TagTop, err error) {
	err = s.client.Call(c, _tagTop, arg, &res)
	return
}
