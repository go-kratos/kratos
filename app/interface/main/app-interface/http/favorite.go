package http

import (
	"strconv"

	"go-common/app/interface/main/app-interface/model/favorite"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// folder get folder.
func folder(c *bm.Context) {
	var (
		aid, vmid, mid int64
		build          int
		err            error
	)
	params := c.Request.Form
	accessKey := params.Get("access_key")
	actionKey := params.Get("actionKey")
	device := params.Get("device")
	mobiApp := params.Get("mobi_app")
	platform := params.Get("platform")
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aid, _ = strconv.ParseInt(params.Get("aid"), 10, 64)
	vmid, _ = strconv.ParseInt(params.Get("vmid"), 10, 64)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(favSvr.Folder(c, accessKey, actionKey, device, mobiApp, platform, build, aid, vmid, mid))
}

func favoriteVideo(c *bm.Context) {
	var (
		mid, vmid, fid int64
		build, tid     int
		pn, ps         int
		err            error
	)
	params := c.Request.Form
	accessKey := params.Get("access_key")
	actionKey := params.Get("actionKey")
	device := params.Get("device")
	mobiApp := params.Get("mobi_app")
	platform := params.Get("platform")
	keyword := params.Get("keyword")
	order := params.Get("order")
	if build, _ = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	tid, _ = strconv.Atoi(params.Get("tid"))
	fid, _ = strconv.ParseInt(params.Get("fid"), 10, 64)
	vmid, _ = strconv.ParseInt(params.Get("vmid"), 10, 64)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(favSvr.FolderVideo(c, accessKey, actionKey, device, mobiApp, platform, keyword, order, build, tid, pn, ps, mid, fid, vmid), nil)
}

func topic(c *bm.Context) {
	var (
		mid    int64
		build  int
		pn, ps int
		err    error
	)
	params := c.Request.Form
	accessKey := params.Get("access_key")
	actionKey := params.Get("actionKey")
	device := params.Get("device")
	mobiApp := params.Get("mobi_app")
	platform := params.Get("platform")
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(favSvr.Topic(c, accessKey, actionKey, device, mobiApp, platform, build, ps, pn, mid), nil)
}

func article(c *bm.Context) {
	var (
		mid    int64
		pn, ps int
	)
	params := c.Request.Form
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(favSvr.Article(c, mid, pn, ps), nil)
}

func favClips(c *bm.Context) {
	var (
		mid    int64
		build  int
		pn, ps int
		err    error
	)
	params := c.Request.Form
	accessKey := params.Get("access_key")
	actionKey := params.Get("actionKey")
	device := params.Get("device")
	mobiApp := params.Get("mobi_app")
	platform := params.Get("platform")
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(favSvr.Clips(c, mid, accessKey, actionKey, device, mobiApp, platform, build, pn, ps), nil)
}

func favAlbums(c *bm.Context) {
	var (
		mid    int64
		build  int
		pn, ps int
		err    error
	)
	params := c.Request.Form
	accessKey := params.Get("access_key")
	actionKey := params.Get("actionKey")
	device := params.Get("device")
	mobiApp := params.Get("mobi_app")
	platform := params.Get("platform")
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(favSvr.Albums(c, mid, accessKey, actionKey, device, mobiApp, platform, build, pn, ps), nil)
}

func specil(c *bm.Context) {
	var (
		build  int
		pn, ps int
		err    error
	)
	params := c.Request.Form
	accessKey := params.Get("access_key")
	actionKey := params.Get("actionKey")
	device := params.Get("device")
	mobiApp := params.Get("mobi_app")
	platform := params.Get("platform")
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	c.JSON(favSvr.Specil(c, accessKey, actionKey, device, mobiApp, platform, build, pn, ps), nil)
}

func audio(c *bm.Context) {
	var (
		mid    int64
		pn, ps int
	)
	params := c.Request.Form
	accessKey := params.Get("access_key")
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(favSvr.Audio(c, accessKey, mid, pn, ps), nil)
}

func tab(c *bm.Context) {
	param := &favorite.TabParam{}
	if err := c.Bind(param); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		param.Mid = midInter.(int64)
	}
	c.JSON(favSvr.Tab(c, param.AccessKey, param.ActionKey, param.Device, param.MobiApp, param.Platform, param.Filtered, param.Build, param.Mid))
}
