package http

import (
	"strconv"
	"time"

	"go-common/app/interface/main/app-intl/model"
	"go-common/app/interface/main/app-intl/model/view"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

const (
	_viewPath     = "/x/intl/view"
	_viewPagePath = "/x/intl/view/page"
)

var (
	// _dislike is.
	_dislike = []*view.Dislike{
		{
			ID:   5,
			Name: "标题党/封面党",
		},
		{
			ID:   6,
			Name: "内容质量差",
		},
		{
			ID:   7,
			Name: "内容/封面令人不适",
		},
		{
			ID:   8,
			Name: "营销广告",
		},
	}
)

// viewIndex view handler
func viewIndex(c *bm.Context) {
	var (
		mid, aid, movieID int64
		err               error
	)
	params := c.Request.Form
	header := c.Request.Header
	// get params
	aidStr := params.Get("aid")
	movieidStr := params.Get("movie_id")
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	ak := params.Get("access_key")
	buildStr := params.Get("build")
	from := params.Get("from")
	trackid := params.Get("trackid")
	network := params.Get("network")
	adExtra := params.Get("ad_extra")
	locale := params.Get("locale")
	// check params
	if aidStr == "" && movieidStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if aidStr != "" && aidStr != "0" {
		if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else if movieidStr != "" && movieidStr != "0" {
		if movieID, err = strconv.ParseInt(movieidStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if aid < 1 && movieID < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	buvid := header.Get("Buvid")
	disid := header.Get("Display-ID")
	cdnIP := header.Get("X-Cache-Server-Addr")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	autoplay, _ := strconv.Atoi(params.Get("autoplay"))
	now := time.Now()
	// view
	ip := metadata.String(c, metadata.RemoteIP)
	viewSvc.ViewInfoc(mid, int(plat), trackid, aidStr, ip, _viewPath, buildStr, buvid, disid, from, now, err, autoplay)
	data, err := viewSvc.View(c, mid, aid, movieID, plat, build, ak, mobiApp, device, buvid, cdnIP, network, adExtra, from, now, locale)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data.Dislikes = _dislike
	c.JSON(data, nil)
	viewSvc.RelateInfoc(mid, aid, int(plat), trackid, buildStr, buvid, disid, ip, _viewPath, data.ReturnCode, data.UserFeature, from, data.Relates, now, data.IsRec)
}

// viewPage view page handler.
func viewPage(c *bm.Context) {
	var (
		mid, aid int64
		build    int
		err      error
	)
	params := c.Request.Form
	header := c.Request.Header
	// get params
	aidStr := params.Get("aid")
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	ak := params.Get("access_key")
	buildStr := params.Get("build")
	from := params.Get("from")
	trackid := params.Get("trackid")
	locale := params.Get("locale")
	// check params
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if aid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if build, err = strconv.Atoi(buildStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	buvid := header.Get("Buvid")
	disid := header.Get("Display-ID")
	cdnIP := header.Get("X-Cache-Server-Addr")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	autoplay, _ := strconv.Atoi(params.Get("autoplay"))
	ip := metadata.String(c, metadata.RemoteIP)
	now := time.Now()
	// view page
	viewSvc.ViewInfoc(mid, int(plat), trackid, aidStr, ip, _viewPagePath, buildStr, buvid, disid, from, now, err, autoplay)
	data, err := viewSvc.ViewPage(c, mid, aid, 0, plat, build, ak, mobiApp, device, cdnIP, false, now, locale)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
