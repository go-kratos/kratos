package http

import (
	"go-common/app/admin/main/config/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"strconv"
	"strings"
)

func updateToken(c *bm.Context) {
	v := new(model.UpdateTokenReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		c.JSON(nil, err)
		return
	}
	if err = svr.UpdateToken(c, v.Env, v.Zone, v.TreeID); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)
}

func create(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.CreateReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	creates := []string{"dev", "fat1", "uat", "pre", "prod"}
	for _, val := range creates {
		if err = svr.CreateApp(v.AppName, val, model.DefaultZone, v.TreeID); err != nil {
			res["message"] = "创建app失败"
			c.JSONMap(res, err)
			return
		}
	}
	c.JSON(nil, err)
}

func appList(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.AppListReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	nodes, err := svr.AuthApps(c, user(c), c.Request.Header.Get("Cookie"))
	if err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	app, err := svr.AppList(c, v.Bu, v.Team, v.AppName, model.DefaultEnv, model.DefaultZone, v.Ps, v.Pn, nodes, v.Status)
	if err != nil {
		res["message"] = "数据获取失败"
		c.JSONMap(res, err)
		return
	}
	result := app
	c.JSON(result, nil)
}

func envsByTeam(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.EnvsByTeamReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	nodes, err := svr.AuthApps(c, user(c), c.Request.Header.Get("Cookie"))
	if err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	data, err := svr.EnvsByTeam(c, v.AppName, v.Zone, nodes)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	result := data
	c.JSON(result, nil)
}

func envs(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.EnvsReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	user := user(c)
	nodes, err := svr.AuthApps(c, user, c.Request.Header.Get("Cookie"))
	if err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.Envs(c, user, v.AppName, v.Zone, v.TreeID, nodes))
}

func nodeTree(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.NodeTreeReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	cookie := c.Request.Header.Get("Cookie")
	user := user(c)
	nodes, err := svr.AuthApps(c, user, cookie)
	if err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.Node(c, user, v.Node, v.Team, cookie, nodes))
}

func zoneCopy(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ZoneCopyReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if v.From == v.To {
		res["message"] = "来源机房和目标机房不能是同一个"
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	if err = svr.ZoneCopy(c, v.AppName, v.From, v.To, v.TreeID); err != nil {
		res["message"] = "拷贝失败"
		c.JSONMap(res, err)
		return
	}
	c.JSON(nil, err)
}

func casterEnvs(c *bm.Context) {
	v := new(model.CasterEnvsReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if v.Auth != "caster_envs_all" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.CasterEnvs(v.Zone, v.TreeID))
}

func rename(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(struct {
		TreeID int64 `form:"tree_id" validate:"required"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(nil, svr.AppRename(v.TreeID, user(c), c.Request.Header.Get("Cookie")))
}

func getApps(c *bm.Context) {
	v := new(struct {
		Name string `form:"name" validate:"required"`
		Env  string `form:"env" validate:"required"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	apps, err := svr.GetApps(v.Env)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	var appIDS []int64
	for _, val := range apps {
		appIDS = append(appIDS, val.ID)
	}
	if len(appIDS) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	builds, err := svr.AllBuilds(appIDS)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	var tagIDS []int64
	for _, val := range builds {
		tagIDS = append(tagIDS, val.TagID)
	}
	tags, err := svr.GetConfigIDS(tagIDS)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	var configIDS []int64
	for _, val := range tags {
		tmpIDs := strings.Split(val.ConfigIDs, ",")
		for _, vv := range tmpIDs {
			id, err := strconv.ParseInt(vv, 10, 64)
			if err != nil {
				log.Error("strconv.ParseInt() error(%v)", err)
				return
			}
			configIDS = append(configIDS, id)
		}
	}
	var appids []int64
	var appslist []*model.App
	var names []string
	if len(configIDS) > 0 {
		configs, err := svr.GetConfigs(configIDS, v.Name)
		if err != nil {
			c.JSON(nil, err)
			return
		}
		for _, val := range configs {
			appids = append(appids, val.AppID)
		}
		appslist, err = svr.IdsGetApps(appids)
		if err != nil {
			c.JSON(nil, err)
			return
		}
		for _, val := range appslist {
			names = append(names, val.Name)
		}
	}
	c.JSON(names, nil)
}

func upAppStatus(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.AppStatusReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if !(v.Status == model.StatusShow || v.Status == model.StatusHidden) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	_, err = svr.AuthApps(c, user(c), c.Request.Header.Get("Cookie"))
	if err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(nil, svr.UpAppStatus(c, v.Status, v.TreeID))
}
