package http

import (
	"go-common/app/service/main/coupon/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func useAllowance(c *bm.Context) {
	var err error
	arg := new(model.ArgUseAllowance)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, svc.UseAllowanceCoupon(c, arg))
}

func allowanceCount(c *bm.Context) {
	var (
		err error
		res []*model.CouponAllowanceInfo
	)
	arg := new(model.ArgCount)
	if err = c.Bind(arg); err != nil {
		return
	}
	res, err = svc.AllowanceCoupon(c, &model.ArgAllowanceCoupons{
		Mid:   arg.Mid,
		State: model.NotUsed,
	})
	c.JSON(len(res), err)
}

func receiveAllowance(c *bm.Context) {
	var err error
	arg := new(model.ArgReceiveAllowance)
	if err = c.Bind(arg); err != nil {
		log.Error("receive allowance bind %+v", err)
		return
	}
	c.JSON(svc.ReceiveAllowance(c, arg))
}

func useNotify(c *bm.Context) {
	var err error
	arg := new(model.ArgAllowanceCheck)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(svc.UseNotify(c, arg))
}

func prizeCards(c *bm.Context) {
	var err error
	arg := new(model.ArgCount)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(svc.PrizeCards(c, arg.Mid))
}

func prizeDraw(c *bm.Context) {
	var err error
	arg := new(model.ArgPrizeDraw)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(svc.PrizeDraw(c, arg.Mid, arg.CardType))
}
