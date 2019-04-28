package http

import (
	"strings"

	"go-common/app/admin/main/config/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func createBuild(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.CreateBuildReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	name := user(c)
	if _, err = svr.AuthApp(c, name, c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	if len(strings.TrimSpace(v.Name)) == 0 {
		res["message"] = "name不能为空"
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	build := &model.Build{}
	build.TagID = v.TagID
	build.Operator = name
	build.Name = v.Name
	if err = svr.CreateBuild(build, v.TreeID, v.Env, v.Zone); err != nil {
		res["message"] = "创建build失败"
		c.JSONMap(res, err)
		return
	}
	c.JSON(nil, err)
}

func builds(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.BuildsReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.Builds(v.TreeID, v.AppName, v.Env, v.Zone))
}

func build(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.BuildReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApps(c, user(c), c.Request.Header.Get("Cookie")); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.Build(v.BuildID))
}

func buildDel(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.BuildReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApps(c, user(c), c.Request.Header.Get("Cookie")); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	if err = svr.GetDelInfos(c, v.BuildID); err != nil {
		res["message"] = "主机列表中有正在使用该build的机器，请让主机离线3小时自动清除后再删除"
		c.JSONMap(res, err)
		return
	}
	if err = svr.Delete(v.BuildID); err != nil {
		res["message"] = "删除build失败"
		c.JSONMap(res, err)
		return
	}
	c.JSON(nil, err)
}
