package http

import (
	"strconv"

	"go-common/app/interface/main/app-resource/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// getParam get param data.
func getParam(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	platStr := params.Get("plat")
	ver := params.Get("ver")
	buildStr := params.Get("build")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("stronv.ParseInt(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var plat int8
	if mobiApp != "" {
		device := params.Get("device")
		plat = model.Plat(mobiApp, device)
	} else if platStr != "" { // android have not mobi_app when 4.18
		var platInt int64
		platInt, err = strconv.ParseInt(platStr, 10, 64)
		if err != nil {
			log.Error("stronv.ParseInt(%s) error(%v)", platStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		plat = int8(platInt)
	}
	// service
	data, version, err := paramSvc.Param(plat, build, ver)
	res := map[string]interface{}{
		"data": data,
		"ver":  version,
	}
	c.JSONMap(res, err)
}
