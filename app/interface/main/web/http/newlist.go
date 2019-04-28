package http

import (
	"strconv"

	"go-common/app/interface/main/web/conf"
	"go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func newList(c *bm.Context) {
	var (
		rid, pn, ps, tp, count int
		rs                     []*api.Arc
		err                    error
	)
	params := c.Request.Form
	ridStr := params.Get("rid")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	tpStr := params.Get("type")
	if ridStr != "" {
		if rid, err = strconv.Atoi(ridStr); err != nil || rid < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 || ps > conf.Conf.Rule.MaxArcsPageSize {
		ps = conf.Conf.Rule.MaxArcsPageSize
	}
	if tpStr != "" {
		if tp, err = strconv.Atoi(tpStr); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if rs, count, err = webSvc.NewList(c, int32(rid), int8(tp), pn, ps); err != nil {
		c.JSON(nil, err)
		log.Error("webSvc.Newlist(%d,%d,%d) error(%v)", rid, pn, ps, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"num":   pn,
		"size":  ps,
		"count": count,
	}
	data["page"] = page
	data["archives"] = rs
	c.JSON(data, nil)
}
