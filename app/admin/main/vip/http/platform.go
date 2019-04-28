package http

import (
	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func platformList(c *bm.Context) {
	arg := new(struct {
		Order string `form:"order" default:"desc"`
		// PN    int    `form:"pn" default:"1"`
		// PS    int    `form:"ps" default:"20"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.PlatformAll(c, arg.Order))
}

func platformInfo(c *bm.Context) {
	arg := new(model.ArgID)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.PlatformByID(c, arg))
}

func platformSave(c *bm.Context) {
	arg := new(model.ConfPlatform)
	if err := c.Bind(arg); err != nil {
		return
	}
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	// platform：必填，可选项为：android、ios、web、public
	if _, ok := model.PlatformMap[arg.Platform]; !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	arg.Operator = username.(string)
	c.JSON(vipSvc.PlatformSave(c, arg))
}

func platformDel(c *bm.Context) {
	arg := new(model.ArgID)
	if err := c.Bind(arg); err != nil {
		return
	}
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	c.JSON(vipSvc.PlatformDel(c, arg, username.(string)))
}
