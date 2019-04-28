package http

import (
	"encoding/json"
	"strconv"

	"go-common/app/service/main/spy/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func info(c *bm.Context) {
	var (
		err    error
		mid    int64
		midStr = c.Request.Form.Get("mid")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ui, err := spySvc.Info(c, mid)
	if err != nil {
		log.Error("spySvc.UserInfo(%d), ui(%v), err(%v)", mid, ui, err)
		c.JSON(nil, ecode.ServerErr)
		return
	}
	c.JSON(&model.UserScore{Mid: ui.Mid, Score: ui.Score}, nil)
}

func purgeUser(c *bm.Context) {
	var (
		err    error
		mid    int64
		params = c.Request.Form
		midStr = params.Get("mid")
		action = params.Get("modifiedAttr")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if action == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = spySvc.PurgeUser(c, mid, action)
	c.JSON(nil, err)
}

func purgeUser2(c *bm.Context) {
	var (
		err    error
		params = c.Request.Form
		msg    = params.Get("msg")
		info   = new(model.NotifyInfo)
	)
	if err = json.Unmarshal([]byte(msg), &info); err != nil {
		log.Error("purgeUser2 (%s) json.Unmarshal error(%v)", msg, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if info == nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if info.Action == "" || info.Mid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = spySvc.PurgeUser(c, info.Mid, info.Action); err != nil {
		log.Error("purgeUser2 (%+v) error(%v)", info, err)
	}
	c.JSON(nil, err)
}

func stat(c *bm.Context) {
	var (
		err     error
		mid, id int64
		params  = c.Request.Form
		midStr  = params.Get("mid")
		idStr   = params.Get("id")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); midStr != "" && err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if id, err = strconv.ParseInt(idStr, 10, 64); idStr != "" && err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if id == 0 && mid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	stat, err := spySvc.StatByID(c, mid, id)
	if err != nil {
		log.Error("spySvc.StatByID err(%v)", err)
		c.JSON(nil, ecode.ServerErr)
		return
	}
	c.JSON(stat, nil)
}
