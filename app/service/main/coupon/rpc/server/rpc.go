package server

import (
	"go-common/app/service/main/coupon/conf"
	"go-common/app/service/main/coupon/model"
	"go-common/app/service/main/coupon/service"
	"go-common/library/log"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC server
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

// SalaryCoupon salary coupon.
func (r *RPC) SalaryCoupon(c context.Context, a *model.ArgSalaryCoupon, res *struct{}) (err error) {
	return r.s.SalaryCoupon(c, a)
}

// SalaryCouponForThird salary coupon for third.
func (r *RPC) SalaryCouponForThird(c context.Context, a *model.ArgSalaryCoupon, res *model.SalaryCouponForThirdResp) (err error) {
	var ct *model.SalaryCouponForThirdResp
	if ct, err = r.s.SalaryCouponForThird(c, a); err == nil && ct != nil {
		*res = *ct
	}
	return
}

// CouponPage coupon page.
func (r *RPC) CouponPage(c context.Context, arg *model.ArgRPCPage, res *model.CouponPageRPCResp) (err error) {
	var (
		cr    *model.CouponPageRPCResp
		count int64
		list  []*model.CouponPageResp
	)
	if count, list, err = r.s.CouponPage(c, arg.Mid, arg.State, arg.Pn, arg.Ps); err != nil {
		log.Error("r.s.CouponPage(%d) err(%+v)", arg.Mid, err)
		return
	}
	cr = &model.CouponPageRPCResp{
		Count: count,
		Res:   list,
	}
	*res = *cr
	return
}

// CouponCartoonPage coupon cartoon page.
func (r *RPC) CouponCartoonPage(c context.Context, arg *model.ArgRPCPage, res *model.CouponCartoonPageResp) (err error) {
	var p *model.CouponCartoonPageResp
	if p, err = r.s.CouponCartoonPage(c, arg.Mid, arg.State, arg.Pn, arg.Ps); err != nil || p == nil {
		log.Error("r.s.CouponCartoonPage(%d) err(%+v)", arg.Mid, err)
		return
	}
	*res = *p
	return
}

// UsableAllowanceCoupon get usable allowance coupon.
func (r *RPC) UsableAllowanceCoupon(c context.Context, a *model.ArgAllowanceCoupon, res *model.CouponAllowancePanelInfo) (err error) {
	var cr *model.CouponAllowancePanelInfo
	if cr, err = r.s.UsableAllowanceCoupon(c, a.Mid, a.Pirce, a.Platform, a.ProdLimMonth, a.ProdLimRenewal); err == nil && cr != nil {
		*res = *cr
	}
	if err != nil {
		log.Error("rpc.UsableAllowanceCoupon(%+v) err(%+v)", a, err)
	}
	return
}

// AllowanceCouponPanel get allowance coupon info for pay panel.
func (r *RPC) AllowanceCouponPanel(c context.Context, a *model.ArgAllowanceCoupon, res *model.CouponAllowancePanelResp) (err error) {
	var (
		cr *model.CouponAllowancePanelResp
		us []*model.CouponAllowancePanelInfo
		ds []*model.CouponAllowancePanelInfo
		ui []*model.CouponAllowancePanelInfo
	)
	if us, ds, ui, err = r.s.AllowancePanelCoupons(c, a.Mid, a.Pirce, a.Platform, a.ProdLimMonth, a.ProdLimRenewal); err != nil {
		log.Error("rpc.AllowancePanelCoupons(%+v) err(%+v)", a, err)
		return
	}
	cr = &model.CouponAllowancePanelResp{
		Usables:  us,
		Disables: ds,
		Using:    ui,
	}
	*res = *cr
	return
}

// MultiUsableAllowanceCoupon get usable allowance coupon by muti pirce.
func (r *RPC) MultiUsableAllowanceCoupon(c context.Context, a *model.ArgUsablePirces, res *map[float64]*model.CouponAllowancePanelInfo) (err error) {
	if *res, err = r.s.MultiUsableAllowanceCoupon(c, a.Mid, a.Pirce, a.Platform, a.ProdLimMonth, a.ProdLimRenewal); err != nil {
		log.Error("rpc.MultiUsableAllowanceCoupon(%+v) err(%+v)", a, err)
		return
	}
	return
}

// JudgeCouponUsable judge coupon is usable.
func (r *RPC) JudgeCouponUsable(c context.Context, a *model.ArgJuageUsable, res *model.CouponAllowanceInfo) (err error) {
	var cp *model.CouponAllowanceInfo
	if cp, err = r.s.JudgeCouponUsable(c, a.Mid, a.Pirce, a.CouponToken, a.Platform, a.ProdLimMonth, a.ProdLimRenewal); err == nil && cp != nil {
		*res = *cp
		return
	}
	if err != nil {
		log.Error("rpc.JudgeCouponUsable(%+v) err(%+v)", a, err)
	}
	return
}

// AllowanceInfo allowance info.
func (r *RPC) AllowanceInfo(c context.Context, a *model.ArgAllowance, res *model.CouponAllowanceInfo) (err error) {
	var cp *model.CouponAllowanceInfo
	if cp, err = r.s.AllowanceInfo(c, a.Mid, a.CouponToken); err == nil && cp != nil {
		*res = *cp
		return
	}
	if err != nil {
		log.Error("rpc.AllowanceInfo(%+v) err(%+v)", a, err)
	}
	return
}

// CancelUseCoupon cancel use coupon .
func (r *RPC) CancelUseCoupon(c context.Context, a *model.ArgAllowance, res *struct{}) (err error) {
	if err = r.s.CancelUseCoupon(c, a.Mid, a.CouponToken); err != nil {
		log.Error("rpc.CancelUseCoupon(%+v) err(%+v)", a, err)
	}
	return
}

// CouponNotify notify coupon .
func (r *RPC) CouponNotify(c context.Context, a *model.ArgNotify, res *struct{}) (err error) {
	if err = r.s.CouponNotify(c, a.Mid, a.OrderNo, a.State); err != nil {
		log.Error("rpc.CouponNotify(%+v) err(%+v)", a, err)
	}
	return
}

// AllowanceList allowance list .
func (r *RPC) AllowanceList(c context.Context, a *model.ArgAllowanceList, res *[]*model.CouponAllowancePanelInfo) (err error) {
	if *res, err = r.s.AllowanceList(c, a.Mid, a.State); err != nil {
		log.Error("rpc.AllowanceList(%+v) err(%+v)", a, err)
	}
	return
}

// UseAllowance use allowance .
func (r *RPC) UseAllowance(c context.Context, a *model.ArgUseAllowance, res *struct{}) (err error) {
	if err = r.s.UseAllowanceCoupon(c, a); err != nil {
		log.Error("rpc.UseAllowanceCoupon(%+v) err(%+v)", a, err)
	}
	return
}

// AllowanceCount allowance count
func (r *RPC) AllowanceCount(c context.Context, a *model.ArgAllowanceMid, res *int) (err error) {
	var rs []*model.CouponAllowanceInfo
	if rs, err = r.s.AllowanceCoupon(c, &model.ArgAllowanceCoupons{
		Mid:   a.Mid,
		State: model.NotUsed,
	}); err == nil {
		*res = len(rs)
	}
	return
}

//ReceiveAllowance receive allowance.
func (r *RPC) ReceiveAllowance(c context.Context, arg *model.ArgReceiveAllowance, res *string) (err error) {
	var couponToken string
	if couponToken, err = r.s.ReceiveAllowance(c, arg); err != nil {
		log.Error("receive allowance(%+v) err(%+v)", arg, err)
		return
	}
	*res = couponToken
	return
}

//PrizeCards .
func (r *RPC) PrizeCards(c context.Context, arg *model.ArgCount, res *[]*model.PrizeCardRep) (err error) {
	if *res, err = r.s.PrizeCards(c, arg.Mid); err != nil {
		log.Error("r.s.PrizeCards(%+v) err(%+v)", arg, err)
		return
	}
	return
}

//PrizeDraw .
func (r *RPC) PrizeDraw(c context.Context, arg *model.ArgPrizeDraw, res *model.PrizeCardRep) (err error) {
	var pc = &model.PrizeCardRep{}
	if pc, err = r.s.PrizeDraw(c, arg.Mid, arg.CardType); err != nil {
		log.Error("r.s.PrizeDraw(%+v) err(%+v)", arg, err)
		return
	}
	*res = *pc
	return
}
