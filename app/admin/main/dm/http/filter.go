package http

import (
	"strconv"

	"go-common/app/admin/main/dm/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func upFilters(c *bm.Context) {
	var (
		p     = c.Request.Form
		pn    = int64(1)
		ps    = int64(50)
		fType = int64(-1)
	)
	mid, err := strconv.ParseInt(p.Get("mid"), 10, 64)
	if err != nil || mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if p.Get("type") != "" {
		if fType, err = strconv.ParseInt(p.Get("type"), 10, 64); err != nil || fType < -1 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}

	if p.Get("pn") != "" {
		if pn, err = strconv.ParseInt(p.Get("pn"), 10, 64); err != nil || pn <= 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if p.Get("ps") != "" {
		if ps, err = strconv.ParseInt(p.Get("ps"), 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	list, total, err := dmSvc.UpFilters(c, mid, fType, pn, ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	pageInfo := &model.PageInfo{
		Num:   pn,
		Size:  ps,
		Total: total,
	}
	data := &model.UpFilterRes{
		Result: list,
		Page:   pageInfo,
	}
	c.JSON(data, nil)
}

func editUpFilters(c *bm.Context) {
	p := c.Request.Form
	mid, err := strconv.ParseInt(p.Get("mid"), 10, 64)
	if err != nil || mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	id, err := strconv.ParseInt(p.Get("id"), 10, 64)
	if err != nil || id < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fType, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil || fType < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	active, err := strconv.ParseInt(p.Get("active"), 10, 64)
	if err != nil || (int8(active) != model.FilterActive && int8(active) != model.FilterUnActive) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if _, err = dmSvc.EditUpFilters(c, id, mid, int8(fType), int8(active)); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
