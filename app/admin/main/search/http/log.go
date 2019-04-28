package http

import (
	"context"

	"go-common/app/admin/main/search/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func logSearch(c *bm.Context) {
	form := c.Request.Form
	appidStr := form.Get("appid")
	switch appidStr {
	case "log_audit":
		logAudit(c)
	case "log_audit_group":
		logAuditGroupBy(c)
	case "log_user_action":
		logUserAction(c)
	default:
		c.JSON(nil, ecode.RequestErr)
	}
}

func bAuth(c *bm.Context, appID string, businessID int) bool {
	if business, ok := svr.Check(appID, businessID); ok && business.PermissionPoint != "" {
		authSrv.Permit(business.PermissionPoint)(c)
		return !c.IsAborted()
	}
	c.JSON(nil, ecode.AccessDenied)
	c.Abort()
	return false
}

func logAudit(c *bm.Context) {
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
	res, err := svr.LogAudit(c, c.Request.Form, sp, business)
	if err != nil {
		log.Error("srv.LogAudit(%v) error(%v)", sp, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, err)
}

func logAuditGroupBy(c *bm.Context) {
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
	res, err := svr.LogAuditGroupBy(c, c.Request.Form, sp, business)
	if err != nil {
		log.Error("srv.LogAuditGroupBy(%v) error(%v)", sp, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, err)
}

func logUserAction(c *bm.Context) {
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
	res, err := svr.LogUserAction(c, c.Request.Form, sp, business)
	if err != nil {
		log.Error("srv.LogUserAction(%v) error(%v)", sp, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, err)
}

func bMlogAudit(c *bm.Context) {
	var (
		err error
		sp  = &model.LogParams{
			Bsp: &model.BasicSearchParams{},
		}
	)
	if err = c.Bind(sp); err != nil {
		return
	}
	if err = c.Bind(sp.Bsp); err != nil {
		return
	}
	if !bAuth(c, "log_audit", sp.Business) {
		return
	}
	business, ok := svr.Check("log_audit", sp.Business)
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if uid, ok := c.Get("uid"); ok {
		go svr.LogCount(context.Background(), "log_audit", sp.Business, uid)
	}
	res, err := svr.LogAudit(c, c.Request.Form, sp, business)
	if err != nil {
		log.Error("srv.bMlogAudit(%v) error(%v)", sp, err)
		return
	}
	c.JSON(res, err)
}

func bMlogAuditGroupBy(c *bm.Context) {
	var (
		err error
		sp  = &model.LogParams{
			Bsp: &model.BasicSearchParams{},
		}
	)
	if err = c.Bind(sp); err != nil {
		return
	}
	if err = c.Bind(sp.Bsp); err != nil {
		return
	}
	if !bAuth(c, "log_audit", sp.Business) {
		return
	}
	business, ok := svr.Check("log_audit", sp.Business)
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res, err := svr.LogAuditGroupBy(c, c.Request.Form, sp, business)
	if err != nil {
		log.Error("srv.bMlogAuditGroupBy(%v) error(%v)", sp, err)
		return
	}
	c.JSON(res, err)
}

func bMlogUserAction(c *bm.Context) {
	var (
		err error
		sp  = &model.LogParams{
			Bsp: &model.BasicSearchParams{},
		}
	)
	if err = c.Bind(sp); err != nil {
		return
	}
	if err = c.Bind(sp.Bsp); err != nil {
		return
	}
	if !bAuth(c, "log_user_action", sp.Business) {
		return
	}
	business, ok := svr.Check("log_user_action", sp.Business)
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if uid, ok := c.Get("uid"); ok {
		go svr.LogCount(context.Background(), "log_user_action", sp.Business, uid)
	}
	res, err := svr.LogUserAction(c, c.Request.Form, sp, business)
	if err != nil {
		log.Error("srv.bMlogUserAction(%v) error(%v)", sp, err)
		return
	}
	c.JSON(res, err)
}
