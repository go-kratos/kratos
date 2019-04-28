package http

import (
	"go-common/app/interface/main/passport-login/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func proxyAddCookie(c *bm.Context) {
	mid := parseInt(c.Request.Form.Get("mid"))
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(srv.ProxyAddCookie(c, mid))
}

func proxyDeleteCookie(c *bm.Context) {
	param := new(model.ParamModifyAuth)
	c.Bind(param)
	if param.Mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err := srv.ProxyDeleteCookie(c, param.Mid, param.Session)
	c.JSON(nil, err)
}

func proxyDeleteCookies(c *bm.Context) {
	mid := parseInt(c.Request.Form.Get("mid"))
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err := srv.ProxyDeleteCookies(c, mid)
	c.JSON(nil, err)
}
