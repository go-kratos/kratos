package http

import (
	"go-common/app/interface/main/app-resource/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"strconv"
	"time"
)

// getStatic get static
func getStatic(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	ver := params.Get("ver")
	buildStr := params.Get("build")
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	data, version, err := staticSvc.Static(plat, build, ver, time.Now())
	res := map[string]interface{}{
		"data": data,
		"ver":  version,
	}
	c.JSONMap(res, err)
}
