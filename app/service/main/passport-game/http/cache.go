package http

import (
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func infoPBCache(c *bm.Context) {
	var (
		err error
		key = c.Request.Form.Get("key")
	)
	if key == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	info, err := srv.InfoPBCache(c, key)
	if err != nil {
		log.Error("service.InfoPBCache(%s) error(%v)", key, err)
		res := map[string]interface{}{}
		res["code"] = err
		res["message"] = err.Error()
		c.JSONMap(res, nil)
		return
	}
	c.JSON(info, nil)
}

func tokenPBCache(c *bm.Context) {
	var (
		err error
		key = c.Request.Form.Get("key")
	)
	if key == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	info, err := srv.TokenPBCache(c, key)
	if err != nil {
		log.Error("service.TokenPBCache(%s) error(%v)", key, err)
		res := map[string]interface{}{}
		res["code"] = err
		res["message"] = err.Error()
		c.JSONMap(res, nil)
		return
	}
	c.JSON(info, nil)
}
