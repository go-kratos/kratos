package http

import (
	"net/url"
	"strconv"

	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	timeZoneMin = -12
	timeZoneMax = 14
)

func report(c *bm.Context) {
	var (
		mid    int64
		params = c.Request.Form
	)
	midItf, _ := c.Get("mid")
	if midItf != nil {
		mid = midItf.(int64)
	}
	platform := params.Get("mobi_app")
	if platform == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Warn("mobi_app is empty")
		return
	}
	buvid := params.Get("buvid")
	if buvid == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Warn("buvid is empty")
		return
	}
	dt := params.Get("device_token")
	build, _ := strconv.Atoi(params.Get("build"))
	if build < 1 {
		c.JSON(nil, ecode.RequestErr)
		log.Warn("build is wrong: %s", params.Get("build"))
		return
	}
	pushSDK, _ := strconv.Atoi(params.Get("push_sdk"))
	platformID := pushmdl.Platform(platform, pushSDK)
	if platformID == pushmdl.PlatformUnknown {
		c.JSON(nil, ecode.RequestErr)
		log.Warn("push_sdk is wrong: %s", params.Get("push_sdk"))
		return
	}
	tm, _ := strconv.Atoi(params.Get("time_zone"))
	if tm < timeZoneMin || tm > timeZoneMax {
		c.JSON(nil, ecode.RequestErr)
		log.Warn("time_zone is wrong: %s", params.Get("time_zone"))
		return
	}
	ns, _ := strconv.Atoi(params.Get("notify_switch"))
	if ns != 0 && ns != 1 {
		c.JSON(nil, ecode.RequestErr)
		log.Warn("notify_switch is wrong: %s", params.Get("notify_switch"))
		return
	}
	// tp, _ := strconv.Atoi(params.Get("type"))
	appID, _ := strconv.ParseInt(params.Get("app_id"), 10, 64)
	if appID == 0 {
		appID = pushmdl.APPIDBBPhone
	}
	deviceBrand, _ := url.QueryUnescape(params.Get("mobile_brand"))
	deviceModel, _ := url.QueryUnescape(params.Get("mobile_model"))
	osVersion, _ := url.QueryUnescape(params.Get("mobile_version"))
	r := &pushmdl.Report{
		APPID:        appID,
		PlatformID:   platformID,
		Mid:          mid,
		Buvid:        buvid,
		DeviceToken:  dt,
		Build:        build,
		TimeZone:     tm,
		NotifySwitch: ns,
		DeviceBrand:  deviceBrand,
		DeviceModel:  deviceModel,
		OSVersion:    osVersion,
	}
	if dt == "" {
		log.Warn("device_token is empty(%+v)", r)
		c.JSON(nil, nil) // iOS可能取不到device token
		return
	}
	c.JSON(nil, pushSrv.PubReport(c, r))
}

func reportOld(ctx *bm.Context) {
	var params = ctx.Request.Form
	buvid := params.Get("buvid")
	if buvid == "" {
		log.Warn("buvid is empty")
		ctx.JSON(nil, nil)
		return
	}
	token := params.Get("device_token")
	if token == "" {
		log.Warn("device_token is empty")
		ctx.JSON(nil, nil)
		return
	}
	ma := params.Get("mobi_app")
	if ma != "" && ma != "android" && ma != "iphone" && ma != "ipad" {
		log.Warn("invalid mobi_app(%s)", ma)
		ctx.JSON(nil, nil)
		return
	}
	mid, _ := strconv.ParseInt(params.Get("mid"), 10, 64)
	pid, _ := strconv.Atoi(params.Get("pid"))
	timezone, _ := strconv.Atoi(params.Get("time_zone"))
	version := params.Get("ver")
	ctx.JSON(nil, pushSrv.ReportOld(ctx, token, buvid, version, mid, pid, timezone))
}
