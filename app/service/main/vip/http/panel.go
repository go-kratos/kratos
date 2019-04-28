package http

import (
	"go-common/app/service/main/vip/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func vipUserMonthPanel(c *bm.Context) {
	arg := new(model.ArgVipConfigMonth)
	if err := c.Bind(arg); err != nil {
		return
	}
	plat := vipSvc.GetPlatID(c, arg.Platform, arg.PanelType, arg.MobiApp, arg.Device, 0)
	c.JSON(vipSvc.VipUserPrice(c, arg.Mid, arg.Month, plat, arg.SubType, arg.IgnoreAutoRenewStatus, arg.Build))
}

func vipPirce(c *bm.Context) {
	arg := new(model.ArgVipConfigMonth)
	if err := c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	plat := vipSvc.GetPlatID(c, arg.Platform, arg.PanelType, arg.MobiApp, arg.Device, 0)
	c.JSON(vipSvc.VipPrice(c, arg.Mid, arg.Month, plat, arg.SubType, arg.CouponToken, arg.Platform, arg.Build))
}

func priceceByProductID(c *bm.Context) {
	arg := new(model.ArgPriceByProduct)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.PriceByProductID(c, arg.ProductID))
}

func priceceByID(c *bm.Context) {
	arg := new(model.ArgVipPriceByID)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.VipPriceByID(c, arg))
}
