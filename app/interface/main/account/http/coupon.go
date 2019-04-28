package http

import (
	"go-common/app/interface/main/account/model"
	v1 "go-common/app/service/main/coupon/api"
	col "go-common/app/service/main/coupon/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func allowanceList(c *bm.Context) {
	var (
		err error
	)
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	arg := new(model.ArgAllowanceList)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(couponSvc.AllowanceList(c, mid.(int64), arg.State))
}

func couponPage(c *bm.Context) {
	var err error
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	arg := new(model.ArgCouponPage)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(couponSvc.CouponPage(c, &col.ArgRPCPage{
		Mid:   mid.(int64),
		State: arg.State,
		Pn:    arg.Pn,
		Ps:    arg.Ps,
	}))
}

// func couponCartoonPage(c *bm.Context) {
// 	var err error
// 	mid, exists := c.Get("mid")
// 	if !exists {
// 		c.JSON(nil, ecode.AccountNotLogin)
// 		return
// 	}
// 	arg := new(model.ArgCouponPage)
// 	if err = c.Bind(arg); err != nil {
// 		return
// 	}
// 	c.JSON(couponSvc.CouponCartoonPage(c, &col.ArgRPCPage{
// 		Mid:   mid.(int64),
// 		State: arg.State,
// 		Pn:    arg.Pn,
// 		Ps:    arg.Ps,
// 	}))
// }

func prizeCards(c *bm.Context) {
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	c.JSON(couponSvc.PrizeCards(c, &col.ArgCount{Mid: mid.(int64)}))
}

func prizeDraw(c *bm.Context) {
	var err error
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	arg := new(model.ArgPrizeDraw)
	if err = c.Bind(arg); err != nil {
		return
	}
	if arg.CardType == -1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(couponSvc.PrizeDraw(c, &col.ArgPrizeDraw{Mid: mid.(int64), CardType: arg.CardType}))
}

func captchaToken(c *bm.Context) {
	c.JSON(couponSvc.CaptchaToken(c, &v1.CaptchaTokenReq{Ip: metadata.String(c, metadata.RemoteIP)}))
}

func useCouponCode(c *bm.Context) {
	var err error
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	arg := new(col.ArgUseCouponCode)
	if err = c.Bind(arg); err != nil {
		return
	}
	arg.IP = metadata.String(c, metadata.RemoteIP)
	arg.Mid = mid.(int64)
	c.JSON(couponSvc.UseCouponCode(c, arg))
}
