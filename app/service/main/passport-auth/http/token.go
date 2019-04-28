package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func tokenInfo(c *bm.Context) {
	ak := c.Request.Form.Get("access_token")
	if len(ak) != 32 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res, err := svr.TokenInfo(c, ak)
	if err == nil {
		c.Set("mid", res.Mid)
	}
	c.JSON(res, err)
}

func refreshInfo(c *bm.Context) {
	rk := c.Request.Form.Get("refresh_token")
	if len(rk) != 32 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res, err := svr.RefreshInfo(c, rk)
	if res != nil {
		c.Set("mid", res.Mid)
	}
	c.JSON(res, err)
}

func oldTokenInfo(c *bm.Context) {
	ak := c.Request.Form.Get("access_token")
	if len(ak) != 32 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res, err := svr.OldTokenInfo(c, ak)
	if res != nil {
		c.Set("mid", res.Mid)
	}
	c.JSON(res, err)
}
