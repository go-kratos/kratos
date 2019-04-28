package http

import (
	usmdl "go-common/app/service/main/usersuit/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const _mobilePendant = "http://account.bilibili.com/mobile/pendant/#/my"

func pointFlagMobile(c *bm.Context) {
	mid, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	var (
		err       error
		pointFlag struct {
			Link struct {
				PendantLink string `json:"pendant_link"`
				MedalLink   string `json:"medal_link"`
			} `json:"link"`
			Flag *usmdl.PointFlag `json:"flag"`
		}
	)
	if pointFlag.Flag, err = usSvc.PointFlag(c, &usmdl.ArgMID{MID: mid.(int64)}); err != nil {
		c.JSON(nil, err)
		return
	}
	pointFlag.Link.PendantLink = _mobilePendant
	c.JSON(pointFlag, nil)
}

func pointFlag(c *bm.Context) {
	mid, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(usSvc.PointFlag(c, &usmdl.ArgMID{MID: mid.(int64)}))
}
