package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func cookieInfo(c *bm.Context) {
	session := c.Request.Form.Get("session")
	if session == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res, err := svr.CookieInfo(c, session)
	if err == nil {
		c.Set("mid", res.Mid)
	}
	c.JSON(res, err)
}
