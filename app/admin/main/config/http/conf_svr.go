package http

import (
	"go-common/app/admin/main/config/model"
	bm "go-common/library/net/http/blademaster"
)

func setToken(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.SetTokenReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	// update & write cache
	if err = svr.SetToken(c, v.TreeID, v.Env, v.Zone, v.Token); err != nil {
		res["message"] = "重置token失败"
		c.JSONMap(res, err)
		return
	}
	c.JSON(nil, err)
}

// hosts client hosts
func hosts(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.HostsReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.Hosts(c, v.TreeID, v.App, v.Env, v.Zone))
}

//clear  host in redis
func clearhost(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.HostsReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(nil, svr.ClearHost(c, v.TreeID, v.Env, v.Zone))
}
