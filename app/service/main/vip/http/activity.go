package http

import (
	"go-common/app/service/main/vip/model"
	bm "go-common/library/net/http/blademaster"
)

func thirdPrizeGrant(c *bm.Context) {
	a := new(model.ArgThirdPrizeGrant)
	if err := c.Bind(a); err != nil {
		return
	}
	a.AppID = model.EleAppID
	c.JSON(nil, vipSvc.ThirdPrizeGrant(c, a))
}
