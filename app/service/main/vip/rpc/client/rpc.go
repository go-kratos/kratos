package client

import (
	"context"

	col "go-common/app/service/main/coupon/model"
	"go-common/app/service/main/vip/model"
	"go-common/library/net/rpc"
)

const (
	_vipInfo              = "RPC.VipInfo"
	_vipInfos             = "RPC.VipInfos"
	_bcoinList            = "RPC.BcoinList"
	_history              = "RPC.History"
	_createOrder          = "RPC.CreateOrder"
	_orderInfo            = "RPC.OrderInfo"
	_tips                 = "RPC.Tips"
	_couponBySuitID       = "RPC.CouponBySuitID"
	_couponBySuitIDV2     = "RPC.CouponBySuitIDV2"
	_couponsForPanel      = "RPC.CouponsForPanel"
	_couponsForPanelV2    = "RPC.CouponsForPanelV2"
	_cancelUseCoupon      = "RPC.CancelUseCoupon"
	_privilegeBySid       = "RPC.PrivilegeBySid"
	_privilegeByType      = "RPC.PrivilegeByType"
	_panelExplain         = "RPC.PanelExplain"
	_jointly              = "RPC.Jointly"
	_resourceBatchOpenVip = "RPC.ResourceBatchOpenVip"
	_orderPayResult       = "RPC.OrderPayResult"
	_surplusFrozenTime    = "RPC.SurplusFrozenTime"
	_unfrozen             = "RPC.Unfrozen"
	_associateVips        = "RPC.AssociateVips"
)

const (
	_appid = "account.service.vip"
)

var (
	_noArg = &struct{}{}
)

// Service is a question service.
type Service struct {
	client *rpc.Client2
}

// New new rpc service.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// Ping def.
func (s *Service) Ping(c context.Context) (res *int, err error) {
	err = s.client.Call(c, "RPC.Ping", nil, res)
	return
}

// VipInfo def.
func (s *Service) VipInfo(c context.Context, arg *model.ArgRPCMid) (res *model.VipInfoResp, err error) {
	res = new(model.VipInfoResp)
	err = s.client.Call(c, _vipInfo, arg, &res)
	return
}

// VipInfos vipinfo list.
func (s *Service) VipInfos(c context.Context, arg *model.ArgRPCMids) (res map[int64]*model.VipInfoResp, err error) {
	err = s.client.Call(c, _vipInfos, arg, &res)
	return
}

// BcoinList bcoin list.
func (s *Service) BcoinList(c context.Context, arg *model.ArgRPCMid) (res *model.BcoinSalaryResp, err error) {
	res = new(model.BcoinSalaryResp)
	err = s.client.Call(c, _bcoinList, arg, &res)
	return
}

// History user change history.
func (s *Service) History(c context.Context, arg *model.ArgRPCHistory) (res []*model.VipChangeHistoryVo, err error) {
	err = s.client.Call(c, _history, arg, &res)
	return
}

// CreateOrder create order.
func (s *Service) CreateOrder(c context.Context, arg *model.ArgRPCCreateOrder) (res map[string]interface{}, err error) {
	err = s.client.Call(c, _createOrder, arg, &res)
	return
}

// OrderInfo order info.
func (s *Service) OrderInfo(c context.Context, arg *model.ArgRPCOrderNo) (res *model.OrderInfo, err error) {
	res = new(model.OrderInfo)
	err = s.client.Call(c, _orderInfo, arg, &res)
	return
}

// Tips info.
func (s *Service) Tips(c context.Context, arg *model.ArgTips) (res []*model.TipsResp, err error) {
	err = s.client.Call(c, _tips, arg, &res)
	return
}

//CouponBySuitID by suit info.
func (s *Service) CouponBySuitID(c context.Context, arg *model.ArgCouponPanel) (res *col.CouponAllowancePanelInfo, err error) {
	res = new(col.CouponAllowancePanelInfo)
	err = s.client.Call(c, _couponBySuitID, arg, &res)
	return
}

//CouponBySuitIDV2 by suit info.
func (s *Service) CouponBySuitIDV2(c context.Context, arg *model.ArgCouponPanelV2) (res *col.CouponAllowancePanelInfo, err error) {
	res = new(col.CouponAllowancePanelInfo)
	err = s.client.Call(c, _couponBySuitIDV2, arg, &res)
	return
}

//CouponsForPanel by suit info.
func (s *Service) CouponsForPanel(c context.Context, arg *model.ArgCouponPanel) (res *col.CouponAllowancePanelResp, err error) {
	res = new(col.CouponAllowancePanelResp)
	err = s.client.Call(c, _couponsForPanel, arg, &res)
	return
}

//CouponsForPanelV2 by suit info.
func (s *Service) CouponsForPanelV2(c context.Context, arg *model.ArgCouponPanelV2) (res *col.CouponAllowancePanelResp, err error) {
	res = new(col.CouponAllowancePanelResp)
	err = s.client.Call(c, _couponsForPanelV2, arg, &res)
	return
}

//CancelUseCoupon cancel use coupon.
func (s *Service) CancelUseCoupon(c context.Context, arg *model.ArgCouponCancel) (err error) {
	err = s.client.Call(c, _cancelUseCoupon, arg, _noArg)
	return
}

// PrivilegeBySid privileges by sid.
func (s *Service) PrivilegeBySid(c context.Context, arg *model.ArgPrivilegeBySid) (res *model.PrivilegesResp, err error) {
	res = new(model.PrivilegesResp)
	err = s.client.Call(c, _privilegeBySid, arg, &res)
	return
}

// PrivilegeByType privileges by type.
func (s *Service) PrivilegeByType(c context.Context, arg *model.ArgPrivilegeDetail) (res []*model.PrivilegeDetailResp, err error) {
	err = s.client.Call(c, _privilegeByType, arg, &res)
	return
}

// PanelExplain panel explain.
func (s *Service) PanelExplain(c context.Context, arg *model.ArgPanelExplain) (res *model.VipPanelExplain, err error) {
	res = new(model.VipPanelExplain)
	err = s.client.Call(c, _panelExplain, arg, &res)
	return
}

//Jointly jointly info.
func (s *Service) Jointly(c context.Context) (res []*model.JointlyResp, err error) {
	err = s.client.Call(c, _jointly, _noArg, &res)
	return
}

//SurplusFrozenTime surplus frozen.
func (s *Service) SurplusFrozenTime(c context.Context, mid int64) (stime int64, err error) {
	err = s.client.Call(c, _surplusFrozenTime, &model.ArgRPCMid{Mid: mid}, &stime)
	return
}

//Unfrozen unfrozen
func (s *Service) Unfrozen(c context.Context, mid int64) (err error) {
	err = s.client.Call(c, _unfrozen, &model.ArgRPCMid{Mid: mid}, _noArg)
	return
}

//ResourceBatchOpenVip resource batch open.
func (s *Service) ResourceBatchOpenVip(c context.Context, arg *model.ArgUseBatch) (err error) {
	err = s.client.Call(c, _resourceBatchOpenVip, arg, _noArg)
	return
}

//OrderPayResult .
func (s *Service) OrderPayResult(c context.Context, arg *model.ArgDialog) (res *model.OrderResult, err error) {
	res = new(model.OrderResult)
	err = s.client.Call(c, _orderPayResult, arg, &res)
	return
}

// AssociateVips associate vips.
func (s *Service) AssociateVips(c context.Context, arg *model.ArgAssociateVip) (res []*model.AssociateVipResp, err error) {
	err = s.client.Call(c, _associateVips, arg, &res)
	return
}
