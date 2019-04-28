package http

import (
	"go-common/app/interface/main/passport-login/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func proxyAddToken(c *bm.Context) {
	mid := parseInt(c.Request.Form.Get("mid"))
	appID := parseInt(c.Request.Form.Get("appid"))
	if mid <= 0 || appID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(srv.ProxyAddToken(c, appID, mid))
}

func proxyDeleteToken(c *bm.Context) {
	param := new(model.ParamModifyAuth)
	c.Bind(param)
	if param.Token == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err := srv.ProxyDeleteToken(c, param.Token)
	c.JSON(nil, err)
}

func proxyDeleteTokens(c *bm.Context) {
	mid := parseInt(c.Request.Form.Get("mid"))
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err := srv.ProxyDeleteTokens(c, mid)
	c.JSON(nil, err)
}

func proxyDeleteGameTokens(c *bm.Context) {
	mid := parseInt(c.Request.Form.Get("mid"))
	appID := parseInt(c.Request.Form.Get("appid"))
	if mid <= 0 || appID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err := srv.ProxyDeleteGameTokens(c, mid, appID)
	c.JSON(nil, err)
}

func proxyRenewToken(c *bm.Context) {
	token := c.Request.Form.Get("token")
	if token == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(srv.ProxyRenewToken(c, token))
}
