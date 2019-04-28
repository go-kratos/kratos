package http

import (
	"go-common/app/service/main/passport-game/model"
	"go-common/app/service/main/passport-game/service"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func oauth(c *bm.Context) {
	var (
		err       error
		params    = c.Request.Form
		accessKey = params.Get("access_key")
		from      = params.Get("from")
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
	var token *model.Token
	if token, err = srv.Oauth(c, app.(*model.App), accessKey, from); err != nil {
		log.Error("service.Oauth(%s) error(%v)", accessKey, err)
		res := map[string]interface{}{}
		if err == service.ErrDispatcherError {
			res["message"] = err.Error()
		}
		c.JSONMap(res, err)
		return
	}
	c.JSON(token, nil)
}

func getKeyOrigin(c *bm.Context) {
	var err error
	var t *model.RSAKey
	if t, err = srv.RSAKeyOrigin(c); err != nil {
		log.Error("service.RSAKeyOrigin() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(t, nil)
}

func getKey(c *bm.Context) {
	c.JSON(srv.RSAKey(c), nil)
}

func getKeyProxy(c *bm.Context) {
	if srv.Proxy(c) {
		getKeyOrigin(c)
		return
	}
	getKey(c)
}

func renewToken(c *bm.Context) {
	var (
		err       error
		params    = c.Request.Form
		accessKey = params.Get("access_key")
		from      = params.Get("from")
	)
	var r *model.RenewToken
	if r, err = srv.RenewToken(c, accessKey, from); err != nil {
		log.Error("service.RenewToken() error(%v)", err)
		res := map[string]interface{}{}
		if err == service.ErrDispatcherError {
			res["message"] = err.Error()
		}
		c.JSONMap(res, err)
		return
	}
	c.JSON(r, nil)
}
