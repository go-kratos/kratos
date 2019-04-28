package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func status(c *bm.Context) {
	var (
		params  = c.Request.Form
		mid, ok = c.Get("mid")
		uuid    = params.Get("uuid")
	)
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	msg, err := memberSvc.Status(c, mid.(int64), uuid)
	if err != nil {
		c.JSON(nil, ecode.RemoteLoginStatusQueryError)
		return
	}
	c.JSON(msg, nil)
}

func closeNotify(c *bm.Context) {
	c.JSON(struct{}{}, nil)
}

func feedback(c *bm.Context) {
	c.JSON(struct{}{}, nil)
}
