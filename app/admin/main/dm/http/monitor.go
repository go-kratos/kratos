package http

import (
	"strconv"

	"go-common/app/admin/main/dm/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func monitorList(c *bm.Context) {
	var (
		err                  error
		tp                   = int64(model.SubTypeVideo)
		aid, cid, mid, state int64
		page, size           int64 = 1, 50
		params                     = c.Request.Form
	)
	if params.Get("type") != "" {
		tp, err = strconv.ParseInt(params.Get("type"), 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	aidStr := params.Get("aid")
	if len(aidStr) > 0 {
		aid, err = strconv.ParseInt(aidStr, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	cidStr := params.Get("cid")
	if len(cidStr) > 0 {
		cid, err = strconv.ParseInt(cidStr, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	midStr := params.Get("mid")
	if len(midStr) > 0 {
		mid, err = strconv.ParseInt(midStr, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	kw := params.Get("keyword")
	if params.Get("state") != "" {
		state, err = strconv.ParseInt(params.Get("state"), 10, 64)
		if err != nil || (int32(state) != model.MonitorBefore && int32(state) != model.MonitorAfter) {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	pageStr := params.Get("page")
	if len(pageStr) > 0 {
		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	psStr := params.Get("page_size")
	if len(psStr) > 0 {
		size, err = strconv.ParseInt(psStr, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	sort := params.Get("sort")
	order := params.Get("order")
	data, err := dmSvc.MonitorList(c, int32(tp), aid, cid, mid, int32(state), kw, sort, order, page, size)
	res := map[string]interface{}{}
	res["data"] = data
	c.JSONMap(res, err)
}

func editMonitor(c *bm.Context) {
	p := c.Request.Form
	tp, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	state, err := strconv.ParseInt(p.Get("state"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if int32(state) != model.MonitorClosed &&
		int32(state) != model.MonitorAfter &&
		int32(state) != model.MonitorBefore {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oids, err := xstr.SplitInts(p.Get("oids"))
	if err != nil || len(oids) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if _, err = dmSvc.UpdateMonitor(c, int32(tp), oids, int32(state)); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
