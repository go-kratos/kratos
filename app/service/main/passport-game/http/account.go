package http

import (
	"strconv"

	"go-common/app/service/main/passport-game/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func myInfo(c *bm.Context) {
	var (
		err       error
		params    = c.Request.Form
		accessKey = params.Get("access_key")
	)
	if accessKey == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	app, ok := c.Get("app")
	if !ok {
		c.JSON(nil, ecode.AppKeyInvalid)
		return
	}
	info, err := srv.MyInfo(c, app.(*model.App), accessKey)
	if err != nil {
		log.Error("service.MyInfo(%s) error(%v)", accessKey, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(info, nil)
}

func info(c *bm.Context) {
	var (
		err    error
		midStr = c.Request.Form.Get("mid")
	)
	if midStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(srv.Info(c, mid), nil)
}
