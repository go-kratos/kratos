package http

import (
	"strconv"

	"go-common/app/service/main/passport-game/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func loginOrigin(c *bm.Context) {
	var (
		err    error
		t      *model.LoginToken
		params = c.Request.Form
		cookie = c.Request.Header.Get("Cookie")
	)
	if t, err = srv.LoginOrigin(c, params.Encode(), cookie); err != nil {
		log.Error("service.LoginOrigin(%s, %s) error(%v)", params.Encode(), cookie, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(t, nil)
}

func login(c *bm.Context) {
	var (
		err      error
		subid    = int64(0)
		params   = c.Request.Form
		subidStr = params.Get("subid")
		userid   = params.Get("userid")
		rsaPwd   = params.Get("pwd")
	)
	if subidStr != "" {
		if subid, err = strconv.ParseInt(subidStr, 10, 32); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	app, ok := c.Get("app")
	if !ok {
		c.JSON(nil, ecode.AppKeyInvalid)
		return
	}
	t, err := srv.Login(c, app.(*model.App), int32(subid), userid, rsaPwd)
	if err != nil {
		log.Error("service.Login() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(t, err)
}

func loginProxy(c *bm.Context) {
	if srv.Proxy(c) {
		loginOrigin(c)
		return
	}
	login(c)
}
