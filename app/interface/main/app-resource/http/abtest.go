package http

import (
	"strconv"

	"go-common/app/interface/main/app-resource/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func abTest(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	buildStr := params.Get("build")
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(abSvc.Experiment(c, plat, build), nil)
}

func abTestV2(c *bm.Context) {
	params := c.Request.Form
	buvid := params.Get("buvid")
	if buvid == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(abSvc.TemporaryABTests(c, buvid), nil)
}

func abserver(c *bm.Context) {
	params := c.Request.Form
	buvid := params.Get("buvid")
	device := params.Get("device")
	mobiAPP := params.Get("mobi_app")
	buildStr := params.Get("build")
	filteredStr := params.Get("filtered")
	if buvid == "" || device == "" || mobiAPP == "" || buildStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(abSvc.AbServer(c, buvid, device, mobiAPP, filteredStr, build, mid))
}
