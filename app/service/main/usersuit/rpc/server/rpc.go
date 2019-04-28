package server

import (
	"go-common/app/service/main/usersuit/conf"
	"go-common/app/service/main/usersuit/model"
	"go-common/app/service/main/usersuit/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC server struct
type RPC struct {
	s *service.Service
}

// New new rpc server.
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

// Buy buy invite code
func (r *RPC) Buy(c context.Context, arg *model.ArgBuy, res *[]*model.Invite) (err error) {
	*res, err = r.s.BuyInvite(c, arg.Mid, arg.Num, arg.IP)
	return
}

// Apply apply invite code
func (r *RPC) Apply(c context.Context, arg *model.ArgApply, res *struct{}) (err error) {
	err = r.s.ApplyInvite(c, arg.Mid, arg.Code, arg.Cookie, arg.IP)
	return
}

// Stat stat code
func (r *RPC) Stat(c context.Context, arg *model.ArgStat, res *model.InviteStat) (err error) {
	var stat *model.InviteStat
	if stat, err = r.s.Stat(c, arg.Mid, arg.IP); err == nil && stat != nil {
		*res = *stat
	}
	return
}

// Equip pendant equip.
func (r *RPC) Equip(c context.Context, arg *model.ArgEquip, res *struct{}) (err error) {
	err = r.s.EquipPendant(c, arg.Mid, arg.Pid, arg.Status, arg.Source)
	return
}

// GrantByMids one pendant give to multiple users.
func (r *RPC) GrantByMids(c context.Context, arg *model.ArgGrantByMids, res *struct{}) (err error) {
	err = r.s.BatchGrantPendantByMid(c, arg.Pid, arg.Expire, arg.Mids)
	return
}

// GroupPendantMid get  group pendant by mid and gid
func (r *RPC) GroupPendantMid(c context.Context, arg *model.ArgGPMID, res *[]*model.GroupPendantList) (err error) {
	*res, err = r.s.GroupPendantMid(c, arg)
	return
}
