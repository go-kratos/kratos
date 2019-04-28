package http

import (
	"fmt"

	"go-common/app/interface/main/app-player/model"
	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func playurl(c *bm.Context) {
	params := &model.Param{}
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
		errCount.Incr(fmt.Sprintf("%s_%d", params.MobiApp, params.Build))
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
	c.JSON(svr.Playurl(c, mid, params, plat, buvid, fp))
}
