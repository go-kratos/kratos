package http

import (
	"strconv"
	"time"

	"go-common/app/interface/main/app-feed/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func tags(c *bm.Context) {
	var mid int64
	params := c.Request.Form
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	// get params
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	ridStr := params.Get("rid")
	ver := params.Get("ver")
	// check params
	rid, err := strconv.Atoi(ridStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	data, version, err := regionSvc.HotTags(c, mid, int16(rid), ver, plat, time.Now())
	c.JSONMap(map[string]interface{}{"data": data, "ver": version}, err)
}

func subTags(c *bm.Context) {
	var (
		mid    int64
		pn, ps int
		err    error
	)
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	params := c.Request.Form
	// check params
	if pn, err = strconv.Atoi(params.Get("pn")); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.Atoi(params.Get("ps")); err != nil || ps < 1 {
		ps = 20
	}
	c.JSON(regionSvc.SubTags(c, mid, pn, ps), nil)
}

func addTag(c *bm.Context) {
	var mid int64
	params := c.Request.Form
	header := c.Request.Header
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	// check params
	tid, _ := strconv.ParseInt(params.Get("tag_id"), 10, 64)
	if tid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rid, _ := strconv.Atoi(params.Get("rid"))
	// get params
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	buildStr := params.Get("build")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	ip := metadata.String(c, metadata.RemoteIP)
	buvid := header.Get(_headerBuvid)
	disid := header.Get(_headerDisplayID)
	now := time.Now()
	err = regionSvc.AddTag(c, mid, tid, now)
	c.JSON(nil, err)
	if err != nil {
		return
	}
	regionSvc.AddTagInfoc(mid, plat, build, buvid, disid, ip, "/subscribe/tags/add", rid, tid, now)
}

func cancelTag(c *bm.Context) {
	var mid int64
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	params := c.Request.Form
	header := c.Request.Header
	tidStr := params.Get("tag_id")
	ridStr := params.Get("rid")
	// check params
	tid, err := strconv.ParseInt(tidStr, 10, 64)
	if err != nil || tid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rid, err := strconv.Atoi(ridStr)
	if err != nil {
		rid = 0
	}
	// get params
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	buildStr := params.Get("build")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	ip := metadata.String(c, metadata.RemoteIP)
	buvid := header.Get(_headerBuvid)
	disid := header.Get(_headerDisplayID)
	now := time.Now()
	err = regionSvc.CancelTag(c, mid, tid, now)
	c.JSON(nil, err)
	if err != nil {
		return
	}
	regionSvc.CancelTagInfoc(mid, plat, build, buvid, disid, ip, "/subscribe/tags/cancel", rid, tid, now)
}
