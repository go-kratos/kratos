package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"strconv"
)

func upPorder(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pd, err := arcSvc.Porder(c, aid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if pd == nil || pd.ID == 0 {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(pd, nil)
}

func arcOrderGameInfo(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	platformStr := params.Get("platform")
	platform, err := strconv.Atoi(platformStr)
	if err != nil || platform <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	gameInfo, err := arcSvc.ArcOrderGameInfo(c, aid, platform, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if gameInfo == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(gameInfo, nil)
}
