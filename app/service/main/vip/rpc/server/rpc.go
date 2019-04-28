package server

import (
	col "go-common/app/service/main/coupon/model"
	"go-common/app/service/main/vip/conf"
	"go-common/app/service/main/vip/model"
	"go-common/app/service/main/vip/service"
	"go-common/library/log"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC represent rpc server
type RPC struct {
	svc *service.Service
}

// New init rpc.
func New(c *conf.Config, s *service.Service) (svr *rpc.Server) {
	r := &RPC{svc: s}
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

// VipInfo vipinfo.
func (r *RPC) VipInfo(c context.Context, arg *model.ArgRPCMid, res *model.VipInfoResp) (err error) {
	var v *model.VipInfoResp
	if v, err = r.svc.ByMid(c, arg.Mid); err != nil {
		log.Error("RPC.VipInfo(%d) err(%+v)", arg.Mid, err)
		return
	}
	if v != nil {
		*res = *v
	}
	return
}

// VipInfos vipinfo list.
func (r *RPC) VipInfos(c context.Context, arg *model.ArgRPCMids, res *map[int64]*model.VipInfoResp) (err error) {
	if *res, err = r.svc.VipInfos(c, arg.Mids); err != nil {
		log.Error("RPC.VipInfos(%v) err(%+v)", arg.Mids, err)
	}
	return
}

// BcoinList bcoin list.
func (r *RPC) BcoinList(c context.Context, arg *model.ArgRPCMid, res *model.BcoinSalaryResp) (err error) {
	var b *model.BcoinSalaryResp
	if b, err = r.svc.BcoinGive(c, arg.Mid); err != nil {
		log.Error("RPC.BcoinGive(%d) err(%+v)", arg.Mid, err)
		return
	}
	if b != nil {
		*res = *b
	}
	return
}

// History user change history.
func (r *RPC) History(c context.Context, arg *model.ArgChangeHistory, res *[]*model.VipChangeHistoryVo) (err error) {

	if *res, err = r.svc.H5History(c, arg); err != nil {
		log.Error("RPC.History(%v) err(%+v)", arg, err)
	}
	return
}

// CreateOrder create order.
func (r *RPC) CreateOrder(c context.Context, arg *model.ArgCreateOrder, res *map[string]interface{}) (err error) {
	if *res, err = r.svc.CreateOrder(c, arg, arg.IP); err != nil {
		log.Error("RPC.CreateOrder(%v) err(%+v)", arg, err)
	}
	return
}

// OrderInfo vipinfo.
func (r *RPC) OrderInfo(c context.Context, arg *model.ArgRPCOrderNo, res *model.OrderInfo) (err error) {
	var v *model.OrderInfo
	if v, err = r.svc.OrderInfo(c, arg.OrderNo); err != nil {
		log.Error("RPC.OrderInfo(%s) err(%+v)", arg.OrderNo, err)
		return
	}
	if v != nil {
		*res = *v
	}
	return
}

// Tips info.
func (r *RPC) Tips(c context.Context, arg *model.ArgTips, res *[]*model.TipsResp) (err error) {
	if *res, err = r.svc.Tips(c, arg); err != nil {
		log.Error("RPC.Tips(%v) err(%+v)", arg, err)
	}
	return
}

// CouponBySuitID coupon by suit id.
func (r *RPC) CouponBySuitID(c context.Context, arg *model.ArgCouponPanel, res *col.CouponAllowancePanelInfo) (err error) {
	var v *col.CouponAllowancePanelInfo
	if v, err = r.svc.CouponBySuitID(c, arg); err == nil && v != nil {
		*res = *v
	}
	if err != nil {
		log.Error("rpc.CouponBySuitID(%+v) err(%+v)", arg, err)
	}
	return
}

// CouponBySuitIDV2 coupon by suit id.
func (r *RPC) CouponBySuitIDV2(c context.Context, arg *model.ArgCouponPanelV2, res *col.CouponAllowancePanelInfo) (err error) {
	var v *col.CouponAllowancePanelInfo
	arg1 := &model.ArgCouponPanel{
		Mid: arg.Mid,
		Sid: arg.Sid,
	}
	if v, err = r.svc.CouponBySuitID(c, arg1); err == nil && v != nil {
		*res = *v
	}
	if err != nil {
		log.Error("rpc.CouponBySuitID(%+v) err(%+v)", arg, err)
	}
	return
}

// CouponsForPanel coupon by suit id.
func (r *RPC) CouponsForPanel(c context.Context, arg *model.ArgCouponPanel, res *col.CouponAllowancePanelResp) (err error) {
	var v *col.CouponAllowancePanelResp
	if v, err = r.svc.CouponsForPanel(c, arg); err == nil && v != nil {
		*res = *v
	}
	if err != nil {
		log.Error("rpc.CouponsForPanel(%+v) err(%+v)", arg, err)
	}
	return
}

// CouponsForPanelV2 coupon by suit id.
func (r *RPC) CouponsForPanelV2(c context.Context, arg *model.ArgCouponPanelV2, res *col.CouponAllowancePanelResp) (err error) {
	var v *col.CouponAllowancePanelResp
	arg1 := &model.ArgCouponPanel{
		Mid: arg.Mid,
		Sid: arg.Sid,
	}
	if v, err = r.svc.CouponsForPanel(c, arg1); err == nil && v != nil {
		*res = *v
	}
	if err != nil {
		log.Error("rpc.CouponsForPanel(%+v) err(%+v)", arg, arg1)
	}
	return
}

// CancelUseCoupon cancel use coupon.
func (r *RPC) CancelUseCoupon(c context.Context, arg *model.ArgCouponCancel, res *struct{}) (err error) {
	if err = r.svc.CancelUseCoupon(c, arg.Mid, arg.CouponToken, arg.IP); err != nil {
		log.Error("rpc.CancelUseCoupon(%+v) err(%+v)", arg, err)
	}
	return
}

// PrivilegeBySid privilege by sid.
func (r *RPC) PrivilegeBySid(c context.Context, arg *model.ArgPrivilegeBySid, res *model.PrivilegesResp) (err error) {
	var v *model.PrivilegesResp
	if v, err = r.svc.PrivilegesBySid(c, arg); err == nil && v != nil {
		*res = *v
	}
	return
}

// PrivilegeByType privilege by type.
func (r *RPC) PrivilegeByType(c context.Context, arg *model.ArgPrivilegeDetail, res *[]*model.PrivilegeDetailResp) (err error) {
	*res, err = r.svc.PrivilegesByType(c, arg)
	return
}

// PanelExplain panel explain.
func (r *RPC) PanelExplain(c context.Context, arg *model.ArgPanelExplain, res *model.VipPanelExplain) (err error) {
	var v *model.VipPanelExplain
	if v, err = r.svc.VipPanelExplain(c, arg); err == nil && v != nil {
		*res = *v
	}
	return
}

// Jointly jointly info.
func (r *RPC) Jointly(c context.Context, arg *struct{}, res *[]*model.JointlyResp) (err error) {
	*res, err = r.svc.Jointly(c)
	return
}

//SurplusFrozenTime surplus frozen time.
func (r *RPC) SurplusFrozenTime(c context.Context, arg *model.ArgRPCMid, stime *int64) (err error) {
	*stime, err = r.svc.SurplusFrozenTime(c, arg.Mid)
	return
}

//Unfrozen unfrozen.
func (r *RPC) Unfrozen(c context.Context, arg *model.ArgRPCMid, res *struct{}) (err error) {
	err = r.svc.Unfrozen(c, arg.Mid)
	return
}

//ResourceBatchOpenVip resource batch open vip.
func (r *RPC) ResourceBatchOpenVip(c context.Context, arg *model.ArgUseBatch, res *struct{}) (err error) {
	err = r.svc.ResourceBatchOpenVip(c, arg)
	return
}

//OrderPayResult .
func (r *RPC) OrderPayResult(c context.Context, arg *model.ArgDialog, res *model.OrderResult) (err error) {
	var o *model.OrderResult
	if o, err = r.svc.OrderPayResult(c, arg.OrderNo, arg.Mid, arg.AppID, arg.Platform, arg.Device, arg.MobiApp, arg.Build, arg.PanelType); err == nil && o != nil {
		*res = *o
	}
	return
}

// AssociateVips get all associate vip infos
func (r *RPC) AssociateVips(c context.Context, arg *model.ArgAssociateVip, res *[]*model.AssociateVipResp) (err error) {
	*res = r.svc.AssociateVips(c, arg)
	return
}
