package http

import (
	"go-common/app/service/main/vip/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func cancelUseCoupon(c *bm.Context) {
	var (
		err error
		arg = new(model.ArgCancelUseCoupon)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, vipSvc.CancelUseCoupon(c, arg.Mid, arg.CouponToken, metadata.String(c, metadata.RemoteIP)))
}

func allowanceInfo(c *bm.Context) {
	var err error
	arg := new(model.ArgCancelUseCoupon)
	if err = c.Bind(arg); err != nil {
		log.Error("use allowance coupon bind %+v", err)
		return
	}
	c.JSON(vipSvc.AllowanceInfo(c, arg.Mid, arg.CouponToken))
}
