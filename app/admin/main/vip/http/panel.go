package http

import (
	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func vipPriceConfigs(c *bm.Context) {
	arg := new(model.ArgVipPrice)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.VipPriceConfigs(c, arg))
}

func vipPriceConfigID(c *bm.Context) {
	arg := new(model.ArgVipPriceID)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.VipPriceConfigID(c, arg))
}

func addVipPriceConfig(c *bm.Context) {
	arg := new(model.ArgAddOrUpVipPrice)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, vipSvc.AddVipPriceConfig(c, arg))
}

func upVipPriceConfig(c *bm.Context) {
	arg := new(model.ArgAddOrUpVipPrice)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, vipSvc.UpVipPriceConfig(c, arg))
}

func delVipPriceConfig(c *bm.Context) {
	arg := new(model.ArgVipPriceID)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, vipSvc.DelVipPriceConfig(c, arg))
}

func vipDPriceConfigs(c *bm.Context) {
	arg := new(model.ArgVipPriceID)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.VipDPriceConfigs(c, arg))
}

func vipDPriceConfigID(c *bm.Context) {
	arg := new(model.ArgVipDPriceID)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.VipDPriceConfigID(c, arg))
}

func addVipDPriceConfig(c *bm.Context) {
	arg := new(model.ArgAddOrUpVipDPrice)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, vipSvc.AddVipDPriceConfig(c, arg))
}

func upVipDPriceConfig(c *bm.Context) {
	arg := new(model.ArgAddOrUpVipDPrice)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, vipSvc.UpVipDPriceConfig(c, arg))
}

func delVipDPriceConfig(c *bm.Context) {
	arg := new(model.ArgVipDPriceID)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, vipSvc.DelVipDPriceConfig(c, arg))
}

func vipPanelTypes(c *bm.Context) {
	c.JSON(vipSvc.PanelPlatFormTypes(c))
}
