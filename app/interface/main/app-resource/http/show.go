package http

import (
	"strconv"

	"go-common/app/interface/main/app-resource/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func tabs(c *bm.Context) {
	var (
		params = c.Request.Form
		header = c.Request.Header
		res    = map[string]interface{}{}
	)
	mobiApp := params.Get("mobi_app")
	buildStr := params.Get("build")
	buvid := header.Get(_headerBuvid)
	device := params.Get("device")
	verStr := params.Get("ver")
	language := params.Get("lang")
	plat := model.Plat(mobiApp, device)
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	data, version, ab, err := showSvc.Tabs(c, plat, build, buvid, verStr, mobiApp, language, mid)
	if ab != nil {
		res["abtest"] = ab
	}
	res["data"] = data
	res["ver"] = version
	c.JSONMap(res, err)
}
