package http

import (
	"strconv"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/dynamic/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func regionTotal(c *bm.Context) {
	c.JSON(dySvc.RegionTotal(c), nil)
}

// regionTagArcs get new arcs of region and hot tag
func regionTagArcs(c *bm.Context) {
	var (
		count      int
		rid, tagID int64
		pn, ps     int64
		arcs       []*api.Arc
		err        error
	)
	query := c.Request.Form
	ridStr := query.Get("rid")
	tagIDStr := query.Get("tag_id")
	pnStr := query.Get("pn")
	psStr := query.Get("ps")
	if rid, err = strconv.ParseInt(ridStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tagID, err = strconv.ParseInt(tagIDStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, err = strconv.ParseInt(pnStr, 10, 64); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.ParseInt(psStr, 10, 64); err != nil || ps < 1 {
		ps = int64(conf.Conf.Rule.NumArcs)
	}
	if arcs, count, err = dySvc.RegionTagArcs3(c, int32(rid), tagID, int(pn), int(ps)); err != nil {
		c.JSON(nil, err)
		log.Error("dySvc.RegionTagArcs(%d,%d) error(%v)", rid, tagID, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"num":   int(pn),
		"size":  int(ps),
		"count": count,
	}
	data["page"] = page
	data["archives"] = arcs
	c.JSON(data, nil)
}

// regionArcs get new arcs of region.
func regionArcs(c *bm.Context) {
	var (
		count  int
		rid    int64
		pn, ps int64
		arcs   []*api.Arc
		err    error
	)
	query := c.Request.Form
	ridStr := query.Get("rid")
	pnStr := query.Get("pn")
	psStr := query.Get("ps")
	if rid, err = strconv.ParseInt(ridStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, err = strconv.ParseInt(pnStr, 10, 64); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.ParseInt(psStr, 10, 64); err != nil || ps < 1 {
		ps = int64(conf.Conf.Rule.NumArcs)
	}
	if arcs, count, err = dySvc.RegionArcs3(c, int32(rid), int(pn), int(ps)); err != nil {
		c.JSON(nil, err)
		log.Error("dySvc.RegionArcs(%d) error(%v)", rid, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"num":   int(pn),
		"size":  int(ps),
		"count": count,
	}
	data["page"] = page
	data["archives"] = arcs
	c.JSON(data, nil)
}

// regionsArcs get batch new arcs of region.
func regionsArcs(c *bm.Context) {
	var (
		ridStr  string
		count   int
		rids    []int32
		ridsTmp []int64
		err     error
	)
	query := c.Request.Form
	if count, err = strconv.Atoi(query.Get("count")); err != nil || count < 1 {
		count = conf.Conf.Rule.NumIndexArcs
	}
	if ridStr = query.Get("rids"); ridStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ridsTmp, err = xstr.SplitInts(ridStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, rid := range ridsTmp {
		rids = append(rids, int32(rid))
	}
	c.JSON(dySvc.RegionsArcs3(c, rids, count))
}
