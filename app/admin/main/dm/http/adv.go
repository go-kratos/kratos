package http

import (
	"strconv"

	"go-common/app/admin/main/dm/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// advList 高级弹幕列表
func advList(c *bm.Context) {
	var (
		p  = c.Request.Form
		pn = int64(1)
		ps = int64(50)
	)
	dmInid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil || dmInid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typ := p.Get("bType")
	if typ == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mode := p.Get("mode")
	if mode == "" {
		c.JSON(nil, ecode.RequestErr)
		return
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
	result, total, err := dmSvc.Advances(c, dmInid, typ, mode, pn, ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	pageInfo := &model.PageInfo{
		Num:   pn,
		Size:  ps,
		Total: total,
	}
	data := &model.AdvanceRes{
		Result: result,
		Page:   pageInfo,
	}
	c.JSON(data, nil)
}
