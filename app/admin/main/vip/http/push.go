package http

import (
	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func pushs(c *bm.Context) {
	arg := new(model.ArgPushData)

	if err := c.Bind(arg); err != nil {
		return
	}
	res, count, err := vipSvc.PushDatas(c, arg)

	result := make(map[string]interface{})
	result["data"] = res
	result["total"] = count
	c.JSON(result, err)
}

func push(c *bm.Context) {
	arg := new(model.ArgID)

	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(vipSvc.GetPushData(c, arg.ID))
}

func disablePush(c *bm.Context) {
	arg := new(model.ArgID)

	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(nil, vipSvc.DisablePushData(c, arg.ID))
}

func delPush(c *bm.Context) {
	arg := new(model.ArgID)

	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(nil, vipSvc.DelPushData(c, arg.ID))
}

func savePush(c *bm.Context) {
	arg := new(model.VipPushData)
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	if err := c.Bind(arg); err != nil {
		return
	}
	arg.Operator = username.(string)
	c.JSON(nil, vipSvc.SavePushData(c, arg))
}
