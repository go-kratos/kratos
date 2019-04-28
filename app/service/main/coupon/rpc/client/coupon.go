package client

import (
	"context"

	"go-common/app/service/main/coupon/model"
	"go-common/library/net/rpc"
)

var (
	_noRes = &struct{}{}
)

const (
	_salaryCoupon               = "RPC.SalaryCoupon"
	_salaryCouponForThird       = "RPC.SalaryCouponForThird"
	_couponPage                 = "RPC.CouponPage"
	_couponCartoonPage          = "RPC.CouponCartoonPage"
	_usableAllowanceCoupon      = "RPC.UsableAllowanceCoupon"
	_allowanceCouponPanel       = "RPC.AllowanceCouponPanel"
	_multiUsableAllowanceCoupon = "RPC.MultiUsableAllowanceCoupon"
	_judgeCouponUsable          = "RPC.JudgeCouponUsable"
	_allowanceInfo              = "RPC.AllowanceInfo"
	_cancelUseCoupon            = "RPC.CancelUseCoupon"
	_couponNotify               = "RPC.CouponNotify"
	_allowanceList              = "RPC.AllowanceList"
	_useAllowance               = "RPC.UseAllowance"
	_allowanceCount             = "RPC.AllowanceCount"
	_receiveAllowance           = "RPC.ReceiveAllowance"
	_prizeCards                 = "RPC.PrizeCards"
	_prizeDraw                  = "RPC.PrizeDraw"
)

const (
	_appid = "account.service.coupon"
)

var (
	_noArg = &struct{}{}
)

// Service struct info.
type Service struct {
	client *rpc.Client2
}

// New create instance of service and return.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return s
}

// SalaryCoupon salary coupon.
func (s *Service) SalaryCoupon(c context.Context, arg *model.ArgSalaryCoupon) (err error) {
	err = s.client.Call(c, _salaryCoupon, arg, _noRes)
	return
}

// SalaryCouponForThird salary coupon.
func (s *Service) SalaryCouponForThird(c context.Context, arg *model.ArgSalaryCoupon) (res *model.SalaryCouponForThirdResp, err error) {
	res = new(model.SalaryCouponForThirdResp)
	err = s.client.Call(c, _salaryCouponForThird, arg, res)
	return
}

// CouponPage coupon page.
func (s *Service) CouponPage(c context.Context, arg *model.ArgRPCPage) (res *model.CouponPageRPCResp, err error) {
	res = new(model.CouponPageRPCResp)
	err = s.client.Call(c, _couponPage, arg, res)
	return
}

// CouponCartoonPage coupon cartoon page.
func (s *Service) CouponCartoonPage(c context.Context, arg *model.ArgRPCPage) (res *model.CouponCartoonPageResp, err error) {
	res = new(model.CouponCartoonPageResp)
	err = s.client.Call(c, _couponCartoonPage, arg, res)
	return
}

// UsableAllowanceCoupon get usable allowance coupon.
func (s *Service) UsableAllowanceCoupon(c context.Context, arg *model.ArgAllowanceCoupon) (res *model.CouponAllowancePanelInfo, err error) {
	res = new(model.CouponAllowancePanelInfo)
	err = s.client.Call(c, _usableAllowanceCoupon, arg, res)
	return
}

// AllowanceCouponPanel get allowance coupon.
func (s *Service) AllowanceCouponPanel(c context.Context, arg *model.ArgAllowanceCoupon) (res *model.CouponAllowancePanelResp, err error) {
	res = new(model.CouponAllowancePanelResp)
	err = s.client.Call(c, _allowanceCouponPanel, arg, res)
	return
}

// MultiUsableAllowanceCoupon get usable allowance coupon by muti pirce.
func (s *Service) MultiUsableAllowanceCoupon(c context.Context, arg *model.ArgUsablePirces) (res map[float64]*model.CouponAllowancePanelInfo, err error) {
	err = s.client.Call(c, _multiUsableAllowanceCoupon, arg, &res)
	return
}

// JudgeCouponUsable judge coupon is usable.
func (s *Service) JudgeCouponUsable(c context.Context, arg *model.ArgJuageUsable) (res *model.CouponAllowanceInfo, err error) {
	res = new(model.CouponAllowanceInfo)
	err = s.client.Call(c, _judgeCouponUsable, arg, res)
	return
}

// AllowanceInfo allowance info.
func (s *Service) AllowanceInfo(c context.Context, arg *model.ArgAllowance) (res *model.CouponAllowanceInfo, err error) {
	res = new(model.CouponAllowanceInfo)
	err = s.client.Call(c, _allowanceInfo, arg, res)
	return
}

// CancelUseCoupon cancel use coupon.
func (s *Service) CancelUseCoupon(c context.Context, arg *model.ArgAllowance) (err error) {
	err = s.client.Call(c, _cancelUseCoupon, arg, _noArg)
	return
}

// CouponNotify notify coupon.
func (s *Service) CouponNotify(c context.Context, arg *model.ArgNotify) (err error) {
	err = s.client.Call(c, _couponNotify, arg, _noArg)
	return
}

// AllowanceList allowance list.
func (s *Service) AllowanceList(c context.Context, arg *model.ArgAllowanceList) (res []*model.CouponAllowancePanelInfo, err error) {
	err = s.client.Call(c, _allowanceList, arg, &res)
	return
}

// UseAllowance use allowance.
func (s *Service) UseAllowance(c context.Context, arg *model.ArgUseAllowance) (err error) {
	err = s.client.Call(c, _useAllowance, arg, _noArg)
	return
}

// AllowanceCount allowance count.
func (s *Service) AllowanceCount(c context.Context, arg *model.ArgAllowanceMid) (res int, err error) {
	err = s.client.Call(c, _allowanceCount, arg, &res)
	return
}

//ReceiveAllowance receive allowance
func (s *Service) ReceiveAllowance(c context.Context, arg *model.ArgReceiveAllowance) (res string, err error) {
	err = s.client.Call(c, _receiveAllowance, arg, &res)
	return
}

// PrizeCards .
func (s *Service) PrizeCards(c context.Context, arg *model.ArgCount) (res []*model.PrizeCardRep, err error) {
	err = s.client.Call(c, _prizeCards, arg, &res)
	return
}

// PrizeDraw .
func (s *Service) PrizeDraw(c context.Context, arg *model.ArgPrizeDraw) (res *model.PrizeCardRep, err error) {
	err = s.client.Call(c, _prizeDraw, arg, &res)
	return
}
