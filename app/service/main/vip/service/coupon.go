package service

import (
	"context"
	"sort"

	"github.com/pkg/errors"

	colapi "go-common/app/service/main/coupon/api"
	col "go-common/app/service/main/coupon/model"
	v1 "go-common/app/service/main/vip/api"
	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// CouponBySuitID get coupon by mid and suit info.
func (s *Service) CouponBySuitID(c context.Context, arg *model.ArgCouponPanel) (cp *col.CouponAllowancePanelInfo, err error) {
	var p *model.VipPriceConfig
	if p = s.vipPriceMap[arg.Sid]; p == nil {
		err = ecode.VipSuitPirceNotFound
		return
	}
	if cp, err = s.couRPC.UsableAllowanceCoupon(c, &col.ArgAllowanceCoupon{
		Mid:            arg.Mid,
		Pirce:          p.DPrice,
		Platform:       int(p.Plat),
		ProdLimMonth:   int8(p.Month),
		ProdLimRenewal: model.MapProdLlimRenewal[p.SubType],
	}); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// CouponBySuitIDV2 get coupon by mid and suit info v2.
func (s *Service) CouponBySuitIDV2(c context.Context, arg *v1.CouponBySuitIDReq) (res *colapi.UsableAllowanceCouponV2Reply, err error) {
	var p *model.VipPriceConfig
	if p = s.vipPriceMap[arg.Sid]; p == nil {
		err = ecode.VipSuitPirceNotFound
		return
	}
	var vps []*model.VipPanelInfo
	if vps, err = s.VipUserPanel(c, arg.Mid, p.Plat, 0, arg.Build); err != nil {
		return
	}
	for _, v := range vps {
		if v.Id == arg.Sid {
			v.Selected = model.PanelSelected
		} else {
			v.Selected = model.PanelNotSelected
		}
	}
	return s.bestCoupon(c, arg.Mid, int64(p.Plat), vps)
}

func (s *Service) bestCoupon(c context.Context, mid int64, platID int64, vps []*model.VipPanelInfo) (res *colapi.UsableAllowanceCouponV2Reply, err error) {
	sort.Slice(vps, func(i int, j int) bool {
		return vps[i].Selected > vps[j].Selected
	})
	prices := []*colapi.ModelPriceInfo{}
	for _, v := range vps {
		prices = append(prices, &colapi.ModelPriceInfo{
			Price:          v.DPrice,
			Plat:           platID,
			ProdLimMonth:   v.Month,
			ProdLimRenewal: int32(model.MapProdLlimRenewal[int8(v.SubType)]),
		})
	}
	if res, err = s.coupongRPC.UsableAllowanceCouponV2(c, &colapi.UsableAllowanceCouponV2Req{
		Mid:       mid,
		PriceInfo: prices,
	}); err != nil {
		return
	}
	return
}

// CouponsForPanel get coupons for vip pirce panel.
func (s *Service) CouponsForPanel(c context.Context, arg *model.ArgCouponPanel) (res *col.CouponAllowancePanelResp, err error) {
	var p *model.VipPriceConfig
	if p = s.vipPriceMap[arg.Sid]; p == nil {
		err = ecode.VipSuitPirceNotFound
		return
	}
	if res, err = s.couRPC.AllowanceCouponPanel(c, &col.ArgAllowanceCoupon{
		Mid:            arg.Mid,
		Pirce:          p.DPrice,
		Platform:       int(p.Plat),
		ProdLimMonth:   int8(p.Month),
		ProdLimRenewal: model.MapProdLlimRenewal[p.SubType],
	}); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//AllowanceInfo get allowance info.
func (s *Service) AllowanceInfo(c context.Context, mid int64, token string) (cp *col.CouponAllowanceInfo, err error) {
	if cp, err = s.couRPC.AllowanceInfo(c, &col.ArgAllowance{
		Mid:         mid,
		CouponToken: token,
	}); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//CancelUseCoupon cancel use coupon.
func (s *Service) CancelUseCoupon(c context.Context, mid int64, token, ip string) (err error) {
	var (
		cp  *col.CouponAllowanceInfo
		o   *model.OrderInfo
		aff int64
	)
	if cp, err = s.couRPC.AllowanceInfo(c, &col.ArgAllowance{
		Mid:         mid,
		CouponToken: token,
	}); err != nil {
		return
	}
	if cp == nil {
		return
	}
	if cp.State != model.CouponInUse {
		err = ecode.CouPonStateCanNotCancelErr
		return
	}
	if o, err = s.dao.OrderInfo(c, cp.OrderNO); err != nil {
		err = errors.WithStack(err)
		return
	}
	if o == nil {
		err = ecode.VipOrderInfoNotFoundErr
		return
	}
	switch o.Status {
	case model.PAYING:
		if aff, err = s.dao.OldUpdateOrderCancel(c, &model.VipPayOrderOld{
			Status:  model.CANCEL,
			OrderNo: o.OrderNo,
		}); err != nil {
			return
		}
		if aff != 1 {
			log.Error("s.dao.OldUpdateOrderCancel aff!=1 order(%s)", o.OrderNo)
			err = ecode.VipOrderCancelFaildErr
			return
		}
	case model.CANCEL:
	case model.FAILED:
		err = s.couRPC.CouponNotify(c, &col.ArgNotify{
			Mid:     o.Mid,
			OrderNo: o.OrderNo,
			State:   col.AllowanceUseFaild,
		})
		return
	default:
		err = ecode.VipOrderStatusPayingErr
		return
	}
	if _, err = s.dao.PayClose(c, o.OrderNo, ip); err != nil {
		log.Error("s.dao.PayClose(%d,%s) err(%+v)", mid, o.OrderNo, err)
		err = ecode.VipOrderCancelFaildErr
		return
	}
	if err = s.couRPC.CancelUseCoupon(c, &col.ArgAllowance{
		Mid:         mid,
		CouponToken: token,
	}); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UseAllowance .
func (s *Service) UseAllowance(c context.Context, arg *col.ArgUseAllowance) (err error) {
	return s.couRPC.UseAllowance(c, arg)
}
