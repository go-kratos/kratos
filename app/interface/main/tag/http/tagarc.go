package http

import (
	"strconv"

	"go-common/app/interface/main/tag/conf"
	"go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// newArcs  get new arcs of tag
func newArcs(c *bm.Context) {
	var (
		rid   int64
		tid   int64
		ps    int
		pn    int
		count int
		tp    int
		arcs  []*api.Arc
		err   error
	)
	params := c.Request.Form
	ridStr := params.Get("rid")
	tidStr := params.Get("tag_id")
	tpStr := params.Get("type")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 || ps > conf.Conf.Tag.MaxArcsPageSize {
		ps = conf.Conf.Tag.MaxArcsPageSize
	}
	if tid, err = strconv.ParseInt(tidStr, 10, 64); err != nil || tid < 1 {
		log.Error("strconv.ParseInt(%s) err(%v)", tidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ridStr == "" {
		// get newest arcs of tag of all region
		if arcs, count, err = svr.NewArcs(c, tid, ps, pn); err != nil {
			c.JSON(nil, err)
			return
		}
	} else {
		// get newest arcs of region'tag
		if rid, err = strconv.ParseInt(ridStr, 10, 64); err != nil || rid < 1 {
			log.Error("strconv.ParseInt(%s) err(%v)", ridStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if tpStr != "" {
			if tp, err = strconv.Atoi(tpStr); err != nil {
				log.Error("strconv.Atoi(%v) err(%v)", tpStr, err)
				c.JSON(nil, ecode.RequestErr)
				return
			}
		}
		if arcs, count, err = svr.RegionNewArcs(c, int32(rid), tid, int8(tp), ps, pn); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	data := make(map[string]interface{}, 2)
	if len(arcs) == 0 {
		arcs = []*api.Arc{}
	}
	data["archives"] = arcs
	data["page"] = map[string]int{
		"num":   pn,
		"size":  ps,
		"count": count,
	}
	c.JSON(data, nil)
}

func detailRankArc(c *bm.Context) {
	var (
		err   error
		count int
		arcs  []*api.Arc
		param = new(struct {
			Tid  int64 `form:"tag_id" validate:"required,gt=0"`
			Prid int64 `form:"prid" validate:"required,gt=0"`
			Pn   int   `form:"pn"`
			Ps   int   `form:"ps"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.Pn < 1 {
		param.Pn = 1
	}
	if param.Ps < 1 || param.Ps > conf.Conf.Tag.MaxArcsPageSize {
		param.Ps = conf.Conf.Tag.MaxArcsPageSize
	}
	if arcs, count, err = svr.DetailRankArc(c, param.Tid, param.Prid, param.Pn, param.Ps); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	data["archives"] = arcs
	data["page"] = map[string]int{
		"num":   param.Pn,
		"size":  param.Ps,
		"count": count,
	}
	c.JSON(data, nil)
}
