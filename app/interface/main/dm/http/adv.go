package http

import (
	"strconv"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// buyAdv 购买高级弹幕
func buyAdv(c *bm.Context) {
	p := c.Request.Form
	mid, _ := c.Get("mid")
	cid, err := strconv.ParseInt(p.Get("cid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mode := p.Get("mode")
	if mode == "" || (mode != model.AdvSpeMode && mode != model.AdvMode) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = dmSvc.BuyAdvance(c, mid.(int64), cid, mode); err != nil {
		c.JSON(nil, err)
		log.Error("dmSvc.BuyAdvance(mid:%v,cid:%d,mode:%s) error(%v)", mid, cid, mode, err)
		return
	}
	res := map[string]interface{}{}
	res["message"] = "已成功购买"
	c.JSONMap(res, err)
}

// advState 高级弹幕状态
func advState(c *bm.Context) {
	p := c.Request.Form
	mid, _ := c.Get("mid")
	cid, err := strconv.ParseInt(p.Get("cid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mode := p.Get("mode")
	if mode == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	state, err := dmSvc.AdvanceState(c, mid.(int64), cid, mode)
	if err != nil {
		c.JSON(nil, err)
		log.Error("dmSvc.AdvState(%v,%d,%s) error(%v)", mid, cid, mode, err)
		return
	}
	c.JSON(state, err)
}

// advList 高级弹幕列表
func advList(c *bm.Context) {
	mid, _ := c.Get("mid")
	list, err := dmSvc.Advances(c, mid.(int64))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(list, err)
}

// passAdv 通过高级弹幕
func passAdv(c *bm.Context) {
	var (
		err    error
		id     int64
		params = c.Request.Form
	)
	mid, _ := c.Get("mid")
	if id, err = strconv.ParseInt(params.Get("id"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.PassAdvance(c, mid.(int64), id)
	c.JSON(nil, err)
}

// denyAdv 拒绝高级弹幕
func denyAdv(c *bm.Context) {
	var (
		err    error
		id     int64
		params = c.Request.Form
	)
	mid, _ := c.Get("mid")
	if id, err = strconv.ParseInt(params.Get("id"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.DenyAdvance(c, mid.(int64), id)
	c.JSON(nil, err)
}

// cancelAdv 取消高级弹幕
func cancelAdv(c *bm.Context) {
	var (
		err    error
		id     int64
		params = c.Request.Form
	)
	mid, _ := c.Get("mid")
	if id, err = strconv.ParseInt(params.Get("id"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.CancelAdvance(c, mid.(int64), id)
	c.JSON(nil, err)
}
