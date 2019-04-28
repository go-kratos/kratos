package http

import (
	"go-common/app/admin/main/config/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/time"
	"strconv"
	"strings"
)

func createComConfig(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.CreateComConfigReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	name := user(c)
	if _, err = svr.AuthApps(c, name, c.Request.Header.Get("Cookie")); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	conf := &model.CommonConf{}
	conf.Operator = name
	conf.State = v.State
	conf.Comment = v.Comment
	conf.Mark = v.Mark
	conf.Name = v.Name
	c.JSON(nil, svr.CreateComConf(conf, v.Team, v.Env, v.Zone, v.SkipLint))
}

func comValue(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ComValueReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApps(c, user(c), c.Request.Header.Get("Cookie")); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.ComConfig(v.ConfigID))
}

func configsByTeam(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ConfigsByTeamReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApps(c, user(c), c.Request.Header.Get("Cookie")); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.ComConfigsByTeam(v.Team, v.Env, v.Zone, v.Ps, v.Pn))
}

func comConfigsByName(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ComConfigsByNameReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApps(c, user(c), c.Request.Header.Get("Cookie")); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.ComConfigsByName(v.Team, v.Env, v.Zone, v.Name))
}

func updateComConfValue(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.UpdateComConfValueReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	user := user(c)
	if _, err = svr.AuthApps(c, user, c.Request.Header.Get("Cookie")); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	conf := &model.CommonConf{}
	conf.Mtime = time.Time(v.Mtime)
	conf.Mark = v.Mark
	conf.ID = v.ID
	conf.State = v.State
	conf.Comment = v.Comment
	conf.Operator = user
	conf.Name = v.Name
	c.JSON(nil, svr.UpdateComConfValue(conf, v.SkipLint))
}

func namesByTeam(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.NamesByTeamReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApps(c, user(c), c.Request.Header.Get("Cookie")); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.NamesByTeam(v.Team, v.Env, v.Zone))
}

func appByTeam(c *bm.Context) {
	v := new(struct {
		CommonConfigID int64 `form:"common_config_id" validate:"required"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	tagMap, err := svr.AppByTeam(v.CommonConfigID)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(tagMap, nil)
}

func tagPush(c *bm.Context) {
	v := new(struct {
		CommonConfigID int64  `form:"common_config_id" validate:"required"`
		Tags           string `form:"tags" validate:"required"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	user := user(c)
	tagMap, err := svr.AppByTeam(v.CommonConfigID)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	res := make(map[int64]interface{})
	tagIDS := strings.Split(v.Tags, ",")
	for _, val := range tagIDS {
		val, _ := strconv.ParseInt(val, 10, 64)
		if _, ok := tagMap[val]; ok {
			err = svr.CommonPush(c, val, v.CommonConfigID, user)
			if err == nil {
				res[val] = "success"
			} else {
				res[val] = "fail"
			}
		} else {
			res[val] = "data error"
		}
	}
	c.JSON(res, nil)
}
