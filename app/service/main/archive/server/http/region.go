package http

import (
	"strconv"

	"go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// regionArcs.
func regionArcs(c *bm.Context) {
	params := c.Request.Form
	ridStr := params.Get("rid")
	tpStr := params.Get("type")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	// check params
	rid, err := strconv.ParseInt(ridStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var tp int64
	if tpStr != "" {
		if tp, err = strconv.ParseInt(tpStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps < 1 || ps > 100 {
		ps = 20
	}
	// service
	var (
		as    []*api.Arc
		count int
	)
	if tp == 0 {
		as, count, err = arcSvc.RegionArcs3(c, int16(rid), pn, ps)
	} else {
		as, count, err = arcSvc.RegionOriginArcs3(c, int16(rid), pn, ps)
	}
	if err != nil {
		c.JSON(nil, err)
		return
	}
	var res struct {
		Archives []*api.Arc `json:"archives"`
		Page     struct {
			Count int `json:"count"`
			Num   int `json:"num"`
			Size  int `json:"size"`
		} `json:"page"`
	}
	res.Archives = as
	res.Page.Num = pn
	res.Page.Size = ps
	res.Page.Count = count
	c.JSON(res, nil)
}

// addRegionArc.
func addRegionArc(c *bm.Context) {
	params := c.Request.Form
	ridStr := params.Get("rid")
	// check params
	rid, _ := strconv.ParseInt(ridStr, 10, 64)
	c.JSON(nil, arcSvc.AddRegionArcs(c, int16(rid)))
}
