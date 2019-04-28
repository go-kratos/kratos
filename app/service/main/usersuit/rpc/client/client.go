package client

import (
	"context"

	"go-common/app/service/main/usersuit/model"
	"go-common/library/net/rpc"
)

const (
	_appid = "account.service.usersuit"
)

var (
	_noRes = &struct{}{}
)

// Service2 struct
type Service2 struct {
	client *rpc.Client2
}

// New Service2 init
func New(c *rpc.ClientConfig) (s *Service2) {
	s = &Service2{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

const (
	_buy             = "RPC.Buy"
	_apply           = "RPC.Apply"
	_stat            = "RPC.Stat"
	_generate        = "RPC.Generate"
	_list            = "RPC.List"
	_equip           = "RPC.Equip"
	_grantByMids     = "RPC.GrantByMids"
	_groupPendantMid = "RPC.GroupPendantMid"
)

// Buy buy invite
func (s *Service2) Buy(c context.Context, arg *model.ArgBuy) (res []*model.Invite, err error) {
	res = make([]*model.Invite, 0)
	err = s.client.Call(c, _buy, arg, &res)
	return
}

// Apply apply
func (s *Service2) Apply(c context.Context, arg *model.ArgApply) (err error) {
	err = s.client.Call(c, _apply, arg, _noRes)
	return
}

// Stat stat
func (s *Service2) Stat(c context.Context, arg *model.ArgStat) (res *model.InviteStat, err error) {
	res = new(model.InviteStat)
	err = s.client.Call(c, _stat, arg, res)
	return
}

// Generate generator
func (s *Service2) Generate(c context.Context, arg *model.ArgGenerate) (res []*model.Invite, err error) {
	res = make([]*model.Invite, 0)
	err = s.client.Call(c, _generate, arg, &res)
	return
}

// List list
func (s *Service2) List(c context.Context, arg *model.ArgList) (res []*model.Invite, err error) {
	res = make([]*model.Invite, 0)
	err = s.client.Call(c, _list, arg, &res)
	return
}

// Equip pendant equip.
func (s *Service2) Equip(c context.Context, arg *model.ArgEquip) (err error) {
	err = s.client.Call(c, _equip, arg, _noRes)
	return
}

// GrantByMids one pendant give to multiple users.
func (s *Service2) GrantByMids(c context.Context, arg *model.ArgGrantByMids) (err error) {
	err = s.client.Call(c, _grantByMids, arg, _noRes)
	return
}

// GroupPendantMid get share group pendant by mid
func (s *Service2) GroupPendantMid(c context.Context, arg *model.ArgGPMID) (res []*model.GroupPendantList, err error) {
	res = make([]*model.GroupPendantList, 0)
	err = s.client.Call(c, _groupPendantMid, arg, &res)
	return
}
