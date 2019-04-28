package http

import (
	"strconv"

	"go-common/app/interface/main/app-show/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// banners get banners.
func banners(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	buildStr := params.Get("build")
	channel := params.Get("channel")
	module := params.Get("module")
	position := params.Get("position")
	// check param
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data := bannerSvc.Display(c, plat, build, channel, module, position, mobiApp)
	returnJSON(c, data, nil)
}
