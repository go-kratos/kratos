package http

import (
	"go-common/app/interface/main/app-resource/model"
	bm "go-common/library/net/http/blademaster"
)

func serverList(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	c.JSON(broadcastSvc.ServerList(c, plat))
}
