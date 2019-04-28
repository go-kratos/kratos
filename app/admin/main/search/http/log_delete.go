package http

import (
	"go-common/app/admin/main/search/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func logDelete(c *bm.Context) {
	form := c.Request.Form
	appidStr := form.Get("appid")
	switch appidStr {
	case "log_audit":
		logAuditDelete(c)
	case "log_user_action":
		logUserActionDelete(c)
	default:
		c.JSON(nil, ecode.RequestErr)
	}
}

func logAuditDelete(c *bm.Context) {
	var (
		err error
		sp  = &model.LogParams{
			Bsp: &model.BasicSearchParams{},
		}
	)
	if err = c.Bind(sp); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = c.Bind(sp.Bsp); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	business, ok := svr.Check("log_audit", sp.Business)
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res, err := svr.LogAuditDelete(c, c.Request.Form, sp, business)
	if err != nil {
		log.Error("srv.LogAuditDelete(%v) error(%v)", sp, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, err)
}

func logUserActionDelete(c *bm.Context) {
	var (
		err error
		sp  = &model.LogParams{
			Bsp: &model.BasicSearchParams{},
		}
	)
	if err = c.Bind(sp); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = c.Bind(sp.Bsp); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	business, ok := svr.Check("log_user_action", sp.Business)
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res, err := svr.LogUserActionDelete(c, c.Request.Form, sp, business)
	if err != nil {
		log.Error("srv.logUserActionDelete(%v) error(%v)", sp, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, err)
}
