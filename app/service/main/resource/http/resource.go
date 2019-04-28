package http

import (
	"strconv"
	"strings"

	"go-common/app/service/main/resource/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func resource(c *bm.Context) {
	var (
		params = c.Request.Form
		rid    int
		err    error
	)
	ridStr := params.Get("rid")
	if rid, err = strconv.Atoi(ridStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(resSvc.Resource(c, rid), nil)
}

func resources(c *bm.Context) {
	var (
		params = c.Request.Form
		rid    int
		rids   []int
		err    error
	)
	ridsStr := params.Get("rids")
	sArr := strings.Split(ridsStr, ",")
	for _, s := range sArr {
		if rid, err = strconv.Atoi(s); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		rids = append(rids, rid)
	}
	c.JSON(resSvc.Resources(c, rids), nil)
}

func indexIcon(c *bm.Context) {
	c.JSON(resSvc.IndexIcon(c), nil)
}

func playerIcon(c *bm.Context) {
	c.JSON(resSvc.PlayerIcon(c))
}

func cmtbox(c *bm.Context) {
	var (
		params = c.Request.Form
		id     int64
		err    error
	)
	if id, err = strconv.ParseInt(params.Get("id"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(resSvc.Cmtbox(c, id))
}

func regionCard(c *bm.Context) {
	var (
		params = c.Request.Form
		err    error
	)
	mobiApp := params.Get("mobi_app")
	buildStr := params.Get("build")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	c.JSON(resSvc.RegionCard(c, plat, build))
}

func audit(c *bm.Context) {
	c.JSON(resSvc.Audit(c), nil)
}
