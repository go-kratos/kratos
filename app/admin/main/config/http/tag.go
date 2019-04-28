package http

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/admin/main/config/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func createTag(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.CreateTagReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	user := user(c)
	if _, err = svr.AuthApp(c, user, c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	tag := &model.Tag{}
	tag.Mark = v.Mark
	tag.ConfigIDs = v.ConfigIDs
	tag.Operator = user
	c.JSON(nil, svr.CreateTag(tag, v.TreeID, v.Env, v.Zone))
}

func lastTags(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.LastTagsReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.LastTags(v.TreeID, v.Env, v.Zone, v.Build))
}

func tagsByBuild(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.TagsByBuildReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.TagsByBuild(v.AppName, v.Env, v.Zone, v.Build, v.Ps, v.Pn, v.TreeID))
}

func tag(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.TagReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApps(c, user(c), c.Request.Header.Get("Cookie")); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.Tag(v.TagID))
}

func updateTag(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.UpdatetagReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if v.Force != 1 {
		v.Force = 0
	}
	user := user(c)
	if _, err = svr.AuthApp(c, user, c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	if len(strings.TrimSpace(v.Build)) == 0 {
		res["message"] = "build不能为空"
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	tag := &model.Tag{}
	tag.Operator = user
	tag.Mark = v.Mark
	tag.ConfigIDs = v.ConfigIDs
	tag.Force = v.Force
	c.JSON(nil, svr.UpdateTag(c, v.TreeID, v.Env, v.Zone, v.Build, tag))
}

func updateTagID(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.UpdateTagIDReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	user := user(c)
	if _, err = svr.AuthApp(c, user, c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	if len(strings.TrimSpace(v.Build)) == 0 {
		res["message"] = "build不能为空"
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	tag := &model.Tag{}
	if tag, err = svr.RollBackTag(v.TagID); err != nil {
		res["message"] = "tag_id有误"
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	tag.Operator = user
	tag.Mark = fmt.Sprintf("回滚操作：原tag_id：%d", v.TagID)
	c.JSON(nil, svr.UpdateTag(c, v.TreeID, v.Env, v.Zone, v.Build, tag))
}

func hostsForce(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.CreateForceReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	user := user(c)
	if _, err = svr.AuthApp(c, user, c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	var mHosts model.MapHosts
	err = json.Unmarshal([]byte(v.Hosts), &mHosts)
	if err != nil {
		res["message"] = "解析hosts失败"
		c.JSONMap(res, err)
		return
	}
	c.JSON(nil, svr.UpdateForce(c, v.TreeID, v.Version, v.Env, v.Zone, v.Build, user, mHosts))
}

func clearForce(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ClearForceReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	user := user(c)
	if _, err = svr.AuthApp(c, user, c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	var mHosts model.MapHosts
	err = json.Unmarshal([]byte(v.Hosts), &mHosts)
	if err != nil {
		res["message"] = "解析hosts失败"
		c.JSONMap(res, err)
		return
	}
	c.JSON(nil, svr.ClearForce(c, v.TreeID, v.Env, v.Zone, v.Build, mHosts))
}

func tagConfigDiff(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.TagConfigDiff)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	tag, err := svr.LastTasConfigDiff(v.TagID, v.AppID, v.BuildID)
	if err != nil {
		res["message"] = "版本未找到"
		c.JSONMap(res, err)
		return
	}
	var config *model.Config
	var configIDS []int64
	if len(tag.ConfigIDs) > 0 {
		tmpIDs := strings.Split(tag.ConfigIDs, ",")
		for _, vv := range tmpIDs {
			id, err := strconv.ParseInt(vv, 10, 64)
			if err != nil {
				log.Error("strconv.ParseInt() error(%v)", err)
				return
			}
			configIDS = append(configIDS, id)
		}
		config, err = svr.GetConfig(configIDS, v.Name)
		if err != nil {
			res["message"] = "配置文件未找到"
			c.JSONMap(res, err)
			return
		}
	}
	c.JSON(config, err)
}
