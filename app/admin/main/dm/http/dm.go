package http

import (
	"strconv"

	"go-common/app/admin/main/dm/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func contentList(c *bm.Context) {
	var (
		v = new(model.SearchDMParams)
	)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(dmSvc.DMSearch(c, v))
}

// xmlCacheFlush flush danmu xml cache.
func xmlCacheFlush(c *bm.Context) {
	var (
		err error
		tp  = int64(model.SubTypeVideo)
		p   = c.Request.Form
	)
	if p.Get("type") != "" {
		if tp, err = strconv.ParseInt(p.Get("type"), 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	dmSvc.XMLCacheFlush(c, int32(tp), oid)
	c.JSON(nil, nil)
}

// dmSearch danmu content List by cid
func dmSearch(c *bm.Context) {
	p := c.Request.Form
	params := &model.SearchDMParams{
		Mid:          model.CondIntNil,
		State:        p.Get("state"),
		Pool:         p.Get("pool"),
		ProgressFrom: model.CondIntNil,
		ProgressTo:   model.CondIntNil,
		CtimeFrom:    model.CondIntNil,
		CtimeTo:      model.CondIntNil,
		Page:         1,
		Size:         100,
		Sort:         p.Get("sort"),
		Order:        p.Get("order"),
		Keyword:      p.Get("keyword"),
		IP:           p.Get("ip"),
		Attrs:        p.Get("attrs"),
	}
	tp, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	params.Type = int32(tp)
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil {
		log.Error("param err oid %s %v", p.Get("oid"), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	params.Oid = oid
	if p.Get("page") != "" {
		if params.Page, err = strconv.ParseInt(p.Get("page"), 10, 64); err != nil {
			log.Error("param err page %s %v", p.Get("page"), err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if p.Get("page_size") != "" {
		if params.Size, err = strconv.ParseInt(p.Get("page_size"), 10, 64); err != nil {
			log.Error("param err page_size %s %v", p.Get("page_size"), err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if p.Get("mid") != "" {
		if params.Mid, err = strconv.ParseInt(p.Get("mid"), 10, 64); err != nil {
			log.Error("param err mid %s %v", p.Get("mid"), err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if p.Get("progress_from") != "" {
		if params.ProgressFrom, err = strconv.ParseInt(p.Get("progress_from"), 10, 64); err != nil {
			log.Error("param err progress_from %s %v", p.Get("progress_from"), err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if p.Get("progress_to") != "" {
		if params.ProgressTo, err = strconv.ParseInt(p.Get("progress_to"), 10, 64); err != nil {
			log.Error("param err progress_to %s %v", p.Get("progress_to"), err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if p.Get("ctime_from") != "" {
		if params.CtimeFrom, err = strconv.ParseInt(p.Get("ctime_from"), 10, 64); err != nil {
			log.Error("param err ctime_from %s %v", p.Get("ctime_from"), err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if p.Get("ctime_to") != "" {
		if params.CtimeTo, err = strconv.ParseInt(p.Get("ctime_to"), 10, 64); err != nil {
			log.Error("param err ctime_to %s %v", p.Get("ctime_to"), err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	data, err := dmSvc.DMSearch(c, params)
	c.JSON(data, err)
}

// editDMState batch operation by danmu content id
func editDMState(c *bm.Context) {
	var (
		moral, reason int64
		p             = c.Request.Form
	)
	tp, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if p.Get("reason_id") != "" {
		reason, err = strconv.ParseInt(p.Get("reason_id"), 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	state, err := strconv.ParseInt(p.Get("state"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if p.Get("moral") != "" {
		moral, err = strconv.ParseInt(p.Get("moral"), 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	adminID, err := strconv.ParseInt(p.Get("adminId"), 10, 64)
	if err != nil || adminID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	uname := p.Get("uname")
	if uname == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("empty uname is not allow")
		return
	}
	remark := p.Get("remark")
	dmids, err := xstr.SplitInts(p.Get("dmids"))
	if err != nil || len(dmids) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.EditDMState(c, int32(tp), int32(state), oid, int8(reason), dmids, float64(moral), adminID, uname, remark)
	c.JSON(nil, err)
}

// editDMPool batch operation by danmu content id
func editDMPool(c *bm.Context) {
	p := c.Request.Form
	tp, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pool, err := strconv.ParseInt(p.Get("pool"), 10, 64)
	if err != nil || (pool != 0 && pool != 1) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	dmids, err := xstr.SplitInts(p.Get("dmids"))
	if err != nil || len(dmids) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	adminID, err := strconv.ParseInt(p.Get("adminId"), 10, 64)
	if err != nil || adminID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.EditDMPool(c, int32(tp), oid, int32(pool), dmids, adminID)
	c.JSON(nil, err)
}

// editDMAttr change attr
func editDMAttr(c *bm.Context) {
	var (
		p     = c.Request.Form
		bit   uint
		value int32
	)
	tp, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	attr, err := strconv.ParseInt(p.Get("attr"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	switch attr {
	case 0: // unprotect dm
		bit = model.AttrProtect
		value = model.AttrNo
	case 1: // protect dm
		bit = model.AttrProtect
		value = model.AttrYes
	default:
		c.JSON(nil, ecode.RequestErr)
		return
	}
	dmids, err := xstr.SplitInts(p.Get("dmids"))
	if err != nil || len(dmids) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	adminID, err := strconv.ParseInt(p.Get("adminId"), 10, 64)
	if err != nil || adminID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.EditDMAttr(c, int32(tp), oid, dmids, bit, value, adminID)
	c.JSON(nil, err)
}

// dmIndexInfo get dm_index info
func dmIndexInfo(c *bm.Context) {
	p := c.Request.Form
	cid, err := strconv.ParseInt(p.Get("cid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	info, err := dmSvc.DMIndexInfo(c, cid)
	c.JSON(info, err)
}

func fixDMCount(c *bm.Context) {
	p := c.Request.Form
	aid, err := strconv.ParseInt(p.Get("aid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.FixDMCount(c, aid)
	c.JSON(nil, err)
}
