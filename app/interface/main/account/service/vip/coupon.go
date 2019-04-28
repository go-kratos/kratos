package vip

import (
	"context"

	"go-common/app/interface/main/account/model"
	col "go-common/app/service/main/coupon/model"
	v1 "go-common/app/service/main/vip/api"
	vipml "go-common/app/service/main/vip/model"
)

// CouponBySuitID coupon by suit id.
func (s *Service) CouponBySuitID(c context.Context, mid int64, sid int64) (res *col.CouponAllowancePanelInfo, err error) {
	res, err = s.vipRPC.CouponBySuitIDV2(c, &vipml.ArgCouponPanelV2{Mid: mid, Sid: sid})
	return
}

// CouponBySuitIDV2 get coupon by mid and suit info.
func (s *Service) CouponBySuitIDV2(c context.Context, a *model.ArgCouponBySuitID) (res *v1.CouponBySuitIDReply, err error) {
	return s.vipgRPC.CouponBySuitID(c, &v1.CouponBySuitIDReq{
		Mid:       a.Mid,
		Sid:       a.Sid,
		MobiApp:   a.MobiApp,
		Device:    a.Device,
		Platform:  a.Platform,
		PanelType: a.PanelType,
		Build:     a.Build,
	})
}

// CouponsForPanel coupon for panel.
func (s *Service) CouponsForPanel(c context.Context, mid int64, sid int64, platform string) (res *col.CouponAllowancePanelResp, err error) {
	res, err = s.vipRPC.CouponsForPanel(c, &vipml.ArgCouponPanel{Mid: mid, Sid: sid, Platform: vipml.PlatformByName[platform]})
	return
}

// CouponsForPanelV2 coupon for panel.
func (s *Service) CouponsForPanelV2(c context.Context, mid int64, sid int64) (res *col.CouponAllowancePanelResp, err error) {
	res, err = s.vipRPC.CouponsForPanelV2(c, &vipml.ArgCouponPanelV2{Mid: mid, Sid: sid})
	return
}

// CancelUseCoupon coupon cancel use.
func (s *Service) CancelUseCoupon(c context.Context, arg *vipml.ArgCancelUseCoupon) (err error) {
	err = s.vipDao.CancelUseCoupon(c, arg)
	return
}
