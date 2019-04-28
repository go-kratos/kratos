package http

import (
	"go-common/app/interface/main/passport-login/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func proxyCheckUserData(c *bm.Context) {
	param := new(model.ParamLogin)
	c.Bind(param)
	if param.UserName == "" || param.Pwd == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(srv.ProxyCheckUser(c, param))
}
