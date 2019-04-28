package http

import (
	"strconv"
	"time"

	"go-common/app/interface/main/app-resource/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// splashs splash handler
func splashs(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	widthStr := params.Get("width")
	heightStr := params.Get("height")
	buildStr := params.Get("build")
	channel := params.Get("channel")
	ver := params.Get("ver")
	// check params
	width, err := strconv.Atoi(widthStr)
	if err != nil {
		log.Error("width(%s) is invalid", widthStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	height, err := strconv.Atoi(heightStr)
	if err != nil {
		log.Error("height(%s) is invalid", heightStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	build, _ := strconv.Atoi(buildStr)
	device := params.Get("device")
	var plat int8
	if mobiApp != "" {
		plat = model.Plat(mobiApp, device)
	} else {
		plat = model.PlatAndroid
	}
	result, version, err := splashSvc.Display(c, plat, width, height, int(build), channel, ver, time.Now())
	res := map[string]interface{}{
		"data": result,
		"ver":  version,
	}
	c.JSONMap(res, err)
}

func birthSplash(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	widthStr := params.Get("width")
	heightStr := params.Get("height")
	year := params.Get("year")
	birth := params.Get("birth")
	if year == "1930" && birth == "0101" {
		log.Error("birth day is empty(%s,%s)", year, birth)
		c.JSON(nil, ecode.NothingFound)
		return
	}
	// check params
	width, err := strconv.Atoi(widthStr)
	if err != nil {
		log.Error("width(%s) is invalid", widthStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	height, err := strconv.Atoi(heightStr)
	if err != nil {
		log.Error("height(%s) is invalid", heightStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	result, err := splashSvc.Birthday(c, plat, width, height, birth)
	c.JSON(result, err)
}

// splashList ad splash handler
func splashList(c *bm.Context) {
	var (
		header = c.Request.Header
		params = c.Request.Form
	)
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	widthStr := params.Get("width")
	heightStr := params.Get("height")
	buildStr := params.Get("build")
	birth := params.Get("birth")
	adExtra := params.Get("ad_extra")
	device := params.Get("device")
	// check params
	width, err := strconv.Atoi(widthStr)
	if err != nil {
		log.Error("width(%s) is invalid", widthStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	height, err := strconv.Atoi(heightStr)
	if err != nil {
		log.Error("height(%s) is invalid", heightStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	build, _ := strconv.Atoi(buildStr)
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	buvid := header.Get(_headerBuvid)
	plat := model.Plat(mobiApp, device)
	result, err := splashSvc.AdList(c, plat, mobiApp, device, buvid, birth, adExtra, height, width, build, mid)
	c.JSON(result, err)
}
