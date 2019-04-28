package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func blacklistIndex(c *bm.Context) {
	param := &struct {
		Mid int64 `form:"mid"`
		Pn  int   `form:"pn" default:"1"`
		Ps  int   `form:"ps" default:"20"`
	}{}
	if err := c.Bind(param); err != nil {
		return
	}
	c.JSON(spcSvc.BlacklistIndex(param.Mid, param.Pn, param.Ps))
}

func blacklistAdd(c *bm.Context) {
	var (
		uid  int64
		name string
	)
	res := map[string]interface{}{}
	param := &struct {
		Mids []int64 `form:"mids,split" validate:"required"`
	}{}
	if err := c.Bind(param); err != nil {
		return
	}
	if uidInter, ok := c.Get("uid"); ok {
		uid = uidInter.(int64)
	}
	if usernameCtx, ok := c.Get("username"); ok {
		name = usernameCtx.(string)
	}
	if err := spcSvc.BlacklistAdd(param.Mids, name, uid); err != nil {
		res["message"] = "添加失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func blacklistUp(c *bm.Context) {
	var (
		uid  int64
		name string
	)
	res := map[string]interface{}{}
	param := &struct {
		ID     int64 `form:"id" validate:"required"`
		Status int   `form:"status" validate:"min=0,gte=0"`
	}{}
	if err := c.Bind(param); err != nil {
		return
	}
	if uidInter, ok := c.Get("uid"); ok {
		uid = uidInter.(int64)
	}
	if usernameCtx, ok := c.Get("username"); ok {
		name = usernameCtx.(string)
	}
	if err := spcSvc.BlacklistUp(param.ID, param.Status, name, uid); err != nil {
		res["message"] = "更新失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}
