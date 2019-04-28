package http

import (
	"strconv"

	"go-common/app/admin/main/dm/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func addTrJob(c *bm.Context) {
	var (
		state int8
		p     = c.Request.Form
	)
	from, err := strconv.ParseInt(p.Get("from"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	to, err := strconv.ParseInt(p.Get("to"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, err := strconv.ParseInt(p.Get("mid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	offset, err := strconv.ParseFloat(p.Get("offset"), 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.AddTransferJob(c, from, to, mid, offset, state)
	c.JSON(nil, err)
}

func transferList(c *bm.Context) {
	var (
		p     = c.Request.URL.Query()
		pn    = int64(1)
		ps    = int64(50)
		state = int64(-1)
	)
	cid, err := strconv.ParseInt(p.Get("cid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if p.Get("state") != "" {
		if state, err = strconv.ParseInt(p.Get("state"), 10, 64); err != nil || state < -1 || state > 3 {
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
	list, total, err := dmSvc.TransferList(c, cid, state, pn, ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	pageInfo := &model.PageInfo{
		Num:   pn,
		Size:  ps,
		Total: total,
	}
	data := &model.TransListRes{
		Result: list,
		Page:   pageInfo,
	}
	c.JSON(data, nil)
}

// reTransferJob retransfer job
func reTransferJob(c *bm.Context) {
	p := c.Request.Form
	id, err := strconv.ParseInt(p.Get("id"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, err := strconv.ParseInt(p.Get("mid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.ReTransferJob(c, id, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
