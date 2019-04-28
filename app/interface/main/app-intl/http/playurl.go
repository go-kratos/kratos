package http

import (
	"go-common/app/interface/main/app-intl/model"
	"go-common/app/interface/main/app-intl/model/player"
	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

var errCount = prom.BusinessErrCount

func playurl(c *bm.Context) {
	params := &player.Param{}
	if err := c.Bind(params); err != nil {
		return
	}
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	header := c.Request.Header
	buvid := header.Get("Buvid")
	fp := header.Get("X-BVC-FINGERPRINT")
	if params.AID <= 0 {
		errCount.Incr("no_aid")
		log.Warn("juranmeichuan aid %s", c.Request.URL.Path+"?"+c.Request.Form.Encode())
		if env.DeployEnv != env.DeployEnvProd {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if params.CID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if params.Qn < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if params.Npcybs < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if params.Otype != "json" && params.Otype != "xml" {
		params.Otype = "json"
	}
	plat := model.Plat(params.MobiApp, params.Device)
	c.JSON(playerSvc.Playurl(c, mid, params, plat, buvid, fp))
}
