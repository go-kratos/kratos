package http

import (
	"strconv"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// dmUpSearch danmu content List by cid
func dmUpSearch(c *bm.Context) {
	p := c.Request.Form
	mid, _ := c.Get("mid")
	params := &model.SearchDMParams{
		Mids:         p.Get("mids"),
		Mode:         p.Get("modes"),
		Pool:         p.Get("pool"),
		ProgressFrom: model.CondIntNil,
		ProgressTo:   model.CondIntNil,
		CtimeFrom:    p.Get("ctime_from"),
		CtimeTo:      p.Get("ctime_to"),
		Pn:           1,
		Ps:           100,
		Sort:         "desc",
		Order:        "ctime",
		Keyword:      p.Get("keyword"),
		Attrs:        p.Get("attrs"),
		State:        "0,2,6",
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
	if p.Get("pn") != "" {
		if params.Pn, err = strconv.ParseInt(p.Get("pn"), 10, 64); err != nil {
			log.Error("param err page number %s %v", p.Get("pn"), err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if p.Get("ps") != "" {
		if params.Ps, err = strconv.ParseInt(p.Get("ps"), 10, 64); err != nil {
			log.Error("param err page_size %s %v", p.Get("page_size"), err)
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
		params.ProgressFrom = params.ProgressFrom * 1000
	}
	if p.Get("progress_to") != "" {
		if params.ProgressTo, err = strconv.ParseInt(p.Get("progress_to"), 10, 64); err != nil {
			log.Error("param err progress_to %s %v", p.Get("progress_to"), err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		params.ProgressTo = params.ProgressTo * 1000
	}
	if p.Get("order") != "" {
		params.Order = p.Get("order")
	}
	if p.Get("sort") != "" {
		params.Sort = p.Get("sort")
	}
	data, err := dmSvc.DMUpSearch(c, mid.(int64), params)
	c.JSON(data, err)
}

// dmUpRecent get dm list by mid fron redis
func dmUpRecent(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			Pn int64 `form:"pn" default:"1"`
			Ps int64 `form:"ps" default:"100"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	mid, _ := c.Get("mid")
	c.JSON(dmSvc.DMUpRecent(c, mid.(int64), v.Pn, v.Ps))
}

// 统计一个每个时间段弹幕数
func dmDistribution(c *bm.Context) {
	p := c.Request.Form
	typ, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	interval, err := strconv.ParseInt(p.Get("interval"), 10, 64)
	if err != nil || interval <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := dmSvc.DMDistribution(c, int32(typ), oid, int32(interval))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func dmUpConfig(c *bm.Context) {
	mid, _ := c.Get("mid")
	advPermit, err := dmSvc.AdvancePermit(c, mid.(int64))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{
		"advance_permit": advPermit,
	}
	c.JSON(data, nil)
}

func upAdvancePermit(c *bm.Context) {
	mid, _ := c.Get("mid")
	p := c.Request.Form
	permit, err := strconv.ParseInt(p.Get("advance_permit"), 10, 64)
	if err != nil || int8(permit) > model.AdvPermitForbid {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, dmSvc.UpdateAdvancePermit(c, mid.(int64), int8(permit)))
}
