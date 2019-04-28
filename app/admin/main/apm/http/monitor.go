package http

import (
	bm "go-common/library/net/http/blademaster"
)

func appNameList(c *bm.Context) {
	c.JSON(apmSvc.AppNameList(c), nil)
}

func prometheusList(c *bm.Context) {
	v := new(struct {
		AppName string `form:"app_name" validate:"required"`
		Method  string `form:"method" validate:"required"`
		MType   string `form:"mtype" default:"count"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	pts, err := apmSvc.PrometheusList(c, v.AppName, v.Method, v.MType)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(pts, nil)
}

func onlineList(c *bm.Context) {
	ols, err := apmSvc.OnlineList(c)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(ols, nil)
}

func broadcastList(c *bm.Context) {
	bcs, err := apmSvc.BroadCastList(c)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(bcs, nil)
}

func databusList(c *bm.Context) {
	dbs, err := apmSvc.DataBusList(c)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(dbs, nil)
}
