package client

import (
	"context"

	v1 "go-common/app/service/main/account/api"
	"go-common/app/service/main/account/model"
	"go-common/library/net/rpc"
)

const (
	_info3        = "RPC.Info3"
	_card3        = "RPC.Card3"
	_infos3       = "RPC.Infos3"
	_infosByName3 = "RPC.InfosByName3"
	_cards3       = "RPC.Cards3"
	_profile3     = "RPC.Profile3"
	_profileStat3 = "RPC.ProfileWithStat3"

	_addExp3   = "RPC.AddExp3"
	_addMoral3 = "RPC.AddMoral3"

	_relation3      = "RPC.Relation3"
	_relations3     = "RPC.Relations3"
	_attentions3    = "RPC.Attentions3"
	_blacks3        = "RPC.Blacks3"
	_richRelations3 = "RPC.RichRelations3"
)

const (
	_appid = "account.service"
)

var (
	_noArg = &struct{}{}
)

// Service3 for server client3
type Service3 struct {
	client *rpc.Client2
}

// New3 for new struct Service2 obj
func New3(c *rpc.ClientConfig) (s *Service3) {
	s = &Service3{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// Info3 receive ArgMid contains mid and real ip, then init user info.
func (s *Service3) Info3(c context.Context, arg *model.ArgMid) (res *v1.Info, err error) {
	res = new(v1.Info)
	err = s.client.Call(c, _info3, arg, res)
	return
}

// Card3 receive ArgMid contains mid and real ip, then init user card.
func (s *Service3) Card3(c context.Context, arg *model.ArgMid) (res *v1.Card, err error) {
	res = new(v1.Card)
	err = s.client.Call(c, _card3, arg, res)
	return
}

// Infos3 receive ArgMids contains mid and real ip, then init user info.
func (s *Service3) Infos3(c context.Context, arg *model.ArgMids) (res map[int64]*v1.Info, err error) {
	err = s.client.Call(c, _infos3, arg, &res)
	return
}

// InfosByName3 receive ArgMids contains mid and real ip, then init user info.
func (s *Service3) InfosByName3(c context.Context, arg *model.ArgNames) (res map[int64]*v1.Info, err error) {
	err = s.client.Call(c, _infosByName3, arg, &res)
	return
}

// Cards3 receive ArgMids contains mid and real ip, then init user card.
func (s *Service3) Cards3(c context.Context, arg *model.ArgMids) (res map[int64]*v1.Card, err error) {
	err = s.client.Call(c, _cards3, arg, &res)
	return
}

// Profile3 get user profile.
func (s *Service3) Profile3(c context.Context, arg *model.ArgMid) (res *v1.Profile, err error) {
	res = new(v1.Profile)
	err = s.client.Call(c, _profile3, arg, res)
	return
}

// ProfileWithStat3 get user profile.
func (s *Service3) ProfileWithStat3(c context.Context, arg *model.ArgMid) (res *model.ProfileStat, err error) {
	res = new(model.ProfileStat)
	err = s.client.Call(c, _profileStat3, arg, res)
	return
}

// AddExp3 receive ArgExp contains mid, money and reason, then add exp for user.
func (s *Service3) AddExp3(c context.Context, arg *model.ArgExp) (err error) {
	err = s.client.Call(c, _addExp3, arg, _noArg)
	return
}

// AddMoral3 receive ArgMoral contains mid, moral, oper, reason and remark, then add moral for user.
func (s *Service3) AddMoral3(c context.Context, arg *model.ArgMoral) (err error) {
	err = s.client.Call(c, _addMoral3, arg, _noArg)
	return
}

// Relation3 get user friend relation.
func (s *Service3) Relation3(c context.Context, arg *model.ArgRelation) (res *model.Relation, err error) {
	res = new(model.Relation)
	err = s.client.Call(c, _relation3, arg, res)
	return
}

// Relations3 batch get user friend relation.
func (s *Service3) Relations3(c context.Context, arg *model.ArgRelations) (res map[int64]*model.Relation, err error) {
	err = s.client.Call(c, _relations3, arg, &res)
	return
}

// Attentions3 get user attentions ,include followings and whispers.
func (s *Service3) Attentions3(c context.Context, arg *model.ArgMid) (res []int64, err error) {
	err = s.client.Call(c, _attentions3, arg, &res)
	return
}

// Blacks3 get user black list.
func (s *Service3) Blacks3(c context.Context, arg *model.ArgMid) (res map[int64]struct{}, err error) {
	err = s.client.Call(c, _blacks3, arg, &res)
	return
}

// RichRelations3 get relation between owner and mids.
func (s *Service3) RichRelations3(c context.Context, arg *model.ArgRichRelation) (res map[int64]int, err error) {
	err = s.client.Call(c, _richRelations3, arg, &res)
	return
}
