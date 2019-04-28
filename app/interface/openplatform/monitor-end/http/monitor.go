package http

import (
	"strconv"

	"go-common/app/interface/openplatform/monitor-end/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func report(c *bm.Context) {
	var (
		params           = &model.LogParams{}
		mid              int64
		ip               = metadata.String(c, metadata.RemoteIP)
		err              error
		buvid, userAgent string
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if err := c.Bind(params); err != nil {
		c.JSON(nil, err)
		return
	}
	if cookie, _ := c.Request.Cookie("buvid3"); cookie != nil {
		buvid = cookie.Value
	}
	if buvid == "" {
		buvid = c.Request.Header.Get("Buvid")
	}
	userAgent = c.Request.Header.Get("User-Agent")
	if err = mfSvc.Report(c, params, mid, ip, buvid, userAgent); err != nil {
		err = ecode.RequestErr
	}
	c.JSON(nil, err)
}

func startConsume(c *bm.Context) {
	var err error
	if err = mfSvc.StartConsume(); err != nil {
		c.JSON(err.Error(), nil)
	}
	c.JSON("success", nil)
}

func stopConsume(c *bm.Context) {
	var err error
	if err = mfSvc.StopConsume(); err != nil {
		c.JSON(err.Error(), nil)
	}
	c.JSON("success", nil)
}

func pauseConsume(c *bm.Context) {
	var (
		t   int64
		err error
	)
	duration := c.Request.Form.Get("duration")
	if t, err = strconv.ParseInt(duration, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
	}
	if err = mfSvc.PauseConsume(t); err != nil {
		c.JSON(err.Error(), nil)
	}
	c.JSON("success", nil)
}
