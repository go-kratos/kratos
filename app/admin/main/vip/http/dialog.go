package http

import (
	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func dialogList(c *bm.Context) {
	arg := new(struct {
		AppID    int64  `form:"app_id"`
		Platform int64  `form:"platform"`
		Status   string `form:"status"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.DialogAll(c, arg.AppID, arg.Platform, arg.Status))
}

func dialogInfo(c *bm.Context) {
	arg := new(model.ArgID)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.DialogByID(c, arg))
}

func dialogSave(c *bm.Context) {
	arg := new(model.ConfDialog)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.StartTime != 0 && arg.EndTime != 0 && arg.StartTime >= arg.EndTime {
		c.JSON(nil, ecode.VipDialogTimeErr)
		return
	}
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	arg.Operator = username.(string)
	c.JSON(vipSvc.DialogSave(c, arg))
}

func dialogEnable(c *bm.Context) {
	arg := new(struct {
		ID    int64 `form:"id" validate:"required"`
		Stage bool  `form:"stage"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	c.JSON(vipSvc.DialogEnable(c, &model.ConfDialog{ID: arg.ID, Stage: arg.Stage, Operator: username.(string)}))
}

func dialogDel(c *bm.Context) {
	arg := new(model.ArgID)
	if err := c.Bind(arg); err != nil {
		return
	}
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	c.JSON(vipSvc.DialogDel(c, arg, username.(string)))
}
