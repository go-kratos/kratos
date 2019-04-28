package http

import (
	"strconv"

	"go-common/app/admin/main/open/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// addApp.
func addApp(c *bm.Context) {
	appname := c.Request.Form.Get("app_name")
	if appname == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("appname is Empty")
		return
	}
	c.JSON(nil, mngSvc.AddApp(c, appname))
}

// deleteApp.
func delApp(c *bm.Context) {
	form := c.Request.Form
	appid, _ := strconv.ParseInt(form.Get("appid"), 10, 64)
	if appid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("appid (%d) is not exit", appid)
		return
	}
	c.JSON(nil, mngSvc.DelApp(c, appid))
}

// updateApp.
func updateApp(c *bm.Context) {
	arg := new(model.AppParams)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, mngSvc.UpdateApp(c, arg))
}

// listApp.
func listApp(c *bm.Context) {
	t := &model.AppListParams{}
	if err := c.Bind(t); err != nil {
		return
	}
	data, total, err := mngSvc.ListApp(c, t)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	page := map[string]int64{
		"num":   t.PN,
		"size":  t.PS,
		"total": total,
	}
	c.JSON(map[string]interface{}{
		"page": page,
		"data": data,
	}, err)
}
