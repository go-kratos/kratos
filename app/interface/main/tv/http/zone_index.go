package http

import (
	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func zoneIdx(c *bm.Context) {
	var (
		req     = c.Request.Form
		typeStr string
		typeV   int
		pageStr string
		pageV   int
	)
	takeBuild(req) // take build number
	if typeStr = req.Get("season_type"); typeStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typeV = atoi(typeStr)
	if pageStr = req.Get("page"); pageStr == "" {
		pageV = 1
	} else {
		pageV = atoi(pageStr)
	}
	seasons, pager, err := tvSvc.LoadZoneIdx(pageV, typeV)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(model.IdxData{
		List:  seasons,
		Pager: pager,
	}, nil)
}
