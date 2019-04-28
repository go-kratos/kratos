package http

import (
	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func jointlys(c *bm.Context) {
	arg := new(model.ArgQueryJointly)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.JointlysByState(c, arg.State))
}

func addJointly(c *bm.Context) {
	arg := new(model.ArgAddJointly)
	if err := c.Bind(arg); err != nil {
		return
	}
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	arg.Operator = username.(string)
	c.JSON(nil, vipSvc.AddJointly(c, arg))
}

func modifyJointly(c *bm.Context) {
	arg := new(model.ArgModifyJointly)
	if err := c.Bind(arg); err != nil {
		return
	}
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	arg.Operator = username.(string)
	c.JSON(nil, vipSvc.ModifyJointly(c, arg))
}

func deleteJointly(c *bm.Context) {
	arg := new(model.ArgJointlyID)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, vipSvc.DeleteJointly(c, arg.ID))
}
