package http

import (
	"strconv"

	"go-common/app/interface/main/dm/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// uptPaSwitch 申请保护弹幕开关
func uptPaSwitch(c *bm.Context) {
	var (
		err    error
		uid    int64
		status int
		params = c.Request.Form
	)
	// uid
	uid, err = strconv.ParseInt(params.Get("uid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// uid
	status, err = strconv.Atoi(params.Get("status"))
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.UptPaSwitch(c, uid, status)
	c.JSON(nil, err)
}

// UptPaStatus 处理保护弹幕申请
func UptPaStatus(c *bm.Context) {
	var (
		err    error
		uid    int64
		status int
		ids    []int64
		params = c.Request.Form
	)

	// uid
	uid, err = strconv.ParseInt(params.Get("uid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if uid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// status
	status, err = strconv.Atoi(params.Get("status"))
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ids, err = xstr.SplitInts(params.Get("ids"))
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.UptPaStatus(c, uid, ids, status)
	c.JSON(nil, err)
}

// paLs 保护弹幕申请列表
func paLs(c *bm.Context) {
	var (
		err      error
		uid, aid int64
		page     int
		data     *model.ApplyListResult
		params   = c.Request.Form
	)
	// uid
	uid, err = strconv.ParseInt(params.Get("uid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aid, err = strconv.ParseInt(params.Get("aid"), 10, 64)
	if err != nil {
		aid = 0
	}
	page, err = strconv.Atoi(params.Get("page"))
	if err != nil {
		page = 1
	}
	data, err = dmSvc.ProtectApplies(c, uid, aid, page, params.Get("sort"))
	if err != nil {
		c.JSON(nil, err)
		log.Error("dmSvc.PaLs(%v,%v,%v,%v) error(%v)", uid, aid, page, params.Get("sort"), err)
		return
	}
	c.JSON(data, nil)
}

// paVideoLs 保护弹幕申请的视频列表
func paVideoLs(c *bm.Context) {
	var (
		err    error
		uid    int64
		params = c.Request.Form
	)
	// uid
	uid, err = strconv.ParseInt(params.Get("uid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := dmSvc.PaVideoLs(c, uid)
	if err != nil {
		log.Error("dmSvc.PaVideoLs(%v) error(%v)", uid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
