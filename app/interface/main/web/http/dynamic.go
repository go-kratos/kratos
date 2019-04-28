package http

import (
	"strconv"

	"go-common/app/interface/main/web/conf"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func dynamicRegion(c *bm.Context) {
	var (
		rid, pn, ps int
		err         error
	)
	params := c.Request.Form
	ridStr := params.Get("rid")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	if rid, err = strconv.Atoi(ridStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 {
		ps = conf.Conf.Rule.DynamicNumArcs
	} else if ps > conf.Conf.Rule.MaxArcsPageSize {
		ps = conf.Conf.Rule.MaxArcsPageSize
	}
	c.JSON(webSvc.DynamicRegion(c, int32(rid), pn, ps))
}

func dynamicRegionTag(c *bm.Context) {
	var (
		tagID       int64
		rid, pn, ps int
		err         error
	)
	params := c.Request.Form
	ridStr := params.Get("rid")
	tagIDStr := params.Get("tag_id")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	if rid, err = strconv.Atoi(ridStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tagID, err = strconv.ParseInt(tagIDStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 {
		ps = conf.Conf.Rule.DynamicNumArcs
	} else if ps > conf.Conf.Rule.MaxArcsPageSize {
		ps = conf.Conf.Rule.MaxArcsPageSize
	}
	c.JSON(webSvc.DynamicRegionTag(c, tagID, int32(rid), pn, ps))
}

func dynamicRegionTotal(c *bm.Context) {
	c.JSON(webSvc.DynamicRegionTotal(c))
}

func dynamicRegions(c *bm.Context) {
	c.JSON(webSvc.DynamicRegions(c))
}
