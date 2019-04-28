package http

import (
	"go-common/app/service/main/coupon/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func userCoupon(c *bm.Context) {
	var (
		err   error
		count int
	)
	arg := new(model.ArgMid)
	if err = c.Bind(arg); err != nil {
		log.Error("user coupon bind %+v", err)
		return
	}
	switch arg.Type {
	case model.CouponVideo:
		if count, err = svc.VideoCouponCount(c, arg.Mid, arg.Type); err != nil {
			log.Error("user coupon(%d) %+v", arg.Mid, err)
		}
	case model.CouponCartoon:
		if count, err = svc.CarToonCouponCount(c, arg.Mid, arg.Type); err != nil {
			log.Error("user coupon(%d) %+v", arg.Mid, err)
		}
	}
	c.JSON(map[string]interface{}{
		"count": count,
		"mid":   arg.Mid,
	}, nil)
}

func useCoupon(c *bm.Context) {
	var (
		err   error
		state int8
		token string
	)
	arg := new(model.ArgUseCoupon)
	if err = c.Bind(arg); err != nil {
		log.Error("use coupon bind %+v", err)
		return
	}
	if state, token, err = svc.UseCoupon(c, arg.Mid, arg.Oid, arg.Remark, arg.OrderNO, arg.Type, arg.Ver); err != nil {
		log.Error("use coupon(%d) %+v", arg.Mid, err)
	}
	log.Info("use coupon ret(%d,%s)", state, token)
	c.JSON(map[string]interface{}{
		"state":        state,
		"coupon_token": token,
	}, nil)
}

func useCartoonCoupon(c *bm.Context) {
	var (
		err   error
		state int8
		token string
	)
	arg := new(model.ArgUseCartoonCoupon)
	if err = c.Bind(arg); err != nil {
		log.Error("use coupon bind %+v", err)
		return
	}
	if state, token, err = svc.CartoonUse(c, arg.Mid, arg.OrderNO, arg.Type, arg.Ver, arg.Remark, arg.Tips, arg.Count); err != nil {
		log.Error("use coupon(%d) %+v", arg.Mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"state":        state,
		"coupon_token": token,
	}, nil)
}

func couponInfo(c *bm.Context) {
	var err error
	arg := new(model.ArgCoupon)
	if err = c.Bind(arg); err != nil {
		log.Error("coupon info bind %+v", err)
		return
	}
	c.JSON(svc.CouponInfo(c, arg.Mid, arg.CouponToken))
}

func addCoupon(c *bm.Context) {
	var err error
	arg := new(model.ArgAdd)
	if err = c.Bind(arg); err != nil {
		log.Error("add coupon bind %+v", err)
		return
	}
	c.JSON(nil, svc.AddCoupon(c, arg.Mid, arg.StartTime, arg.ExpireTime, arg.Type, arg.Origin))
}

func changeCoupon(c *bm.Context) {
	var err error
	arg := new(model.ChangeCoupon)
	if err = c.Bind(arg); err != nil {
		log.Error("coupon info bind %+v", err)
		return
	}
	c.JSON(nil, svc.ChangeState(c, arg.Mid, arg.UseVer, arg.Ver, arg.CouponToken))
}

func salaryCoupon(c *bm.Context) {
	var err error
	arg := new(model.ArgSalaryCoupon)
	if err = c.Bind(arg); err != nil {
		log.Error("coupon info bind %+v", err)
		return
	}
	arg.Origin = int64(model.VipSalary)
	if err = svc.SalaryCoupon(c, arg); err != nil {
		log.Error("svc.SalaryCoupon(%d) %+v", arg.Mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
