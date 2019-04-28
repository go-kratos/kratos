package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"

	"go-common/app/service/main/resource/model"
)

// pasterAPP get paster for APP
func pasterAPP(c *bm.Context) {
	var (
		params             = c.Request.Form
		aid, typeID, buvid string
		platform, adType   int
		err                error
	)
	aid = params.Get("aid")
	typeID = params.Get("type_id")
	if aid == "" && typeID == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if buvid = params.Get("buvid"); buvid == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if platform, err = strconv.Atoi(params.Get("platform")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if adType, err = strconv.Atoi(params.Get("type")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(resSvc.PasterAPP(c, int8(platform), int8(adType), aid, typeID, buvid))
}

// pasterPGC get paster for PGC
func pasterPGC(c *bm.Context) {
	var (
		params                = c.Request.Form
		sid, platform, device string
		adType                int
		plat                  int8
		err                   error
	)
	sid = params.Get("season_id")
	if _, err = strconv.ParseInt(sid, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if adType, err = strconv.Atoi(params.Get("type")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if platform = params.Get("platform"); platform == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	device = params.Get("device")
	if platform == "web" {
		plat = model.PlatWEB
	} else {
		plat = model.Plat(platform, device)
	}
	c.JSON(resSvc.PasterPGC(c, plat, int8(adType), sid))
}
