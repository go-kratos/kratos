package http

import (
	"go-common/app/service/main/vip/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func buyVipWithPoint(c *bm.Context) {
	var (
		err error
	)
	arg := new(model.ArgBuyPoint)
	if err = c.Bind(arg); err != nil {
		log.Error("buyVipWithPoint Bind err(%+v)", err)
		return
	}
	if err = vipSvc.BuyVipWithPoint(c, arg.Mid, arg.Month); err != nil {
		log.Error("BuyVipWithPoint(%d)  err(%+v)", arg.Mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func rule(c *bm.Context) {
	c.JSON(vipSvc.PointRule(c))
}
