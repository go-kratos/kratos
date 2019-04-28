package http

import (
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// banner get banner
func banner(c *bm.Context) {
	params := c.Request.Form
	aid, _ := strconv.ParseInt(params.Get("aid"), 10, 64)
	platStr := params.Get("plat")
	plat, err := strconv.Atoi(platStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mobiApp := params.Get("mobi_app")
	resIds := params.Get("resource_ids")
	buildStr := params.Get("build")
	channel := params.Get("channel")
	network := params.Get("network")
	isAd, _ := strconv.ParseBool(params.Get("is_ad"))
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	buvid := params.Get("buvid")
	ip := params.Get("ip")
	device := params.Get("device")
	openEvent := params.Get("open_event")
	version := params.Get("version")
	adExtra := params.Get("ad_extra")
	banner := resSvc.Banners(c, int8(plat), build, aid, mid, resIds, channel, ip, buvid, network, mobiApp, device, openEvent, adExtra, version, isAd)
	if banner != nil {
		c.JSON(banner.Banner, nil)
	}
}
