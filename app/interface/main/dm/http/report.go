package http

import (
	"strconv"

	"go-common/app/interface/main/dm/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

const (
	_maxContentLen = 100
)

func addReport(c *bm.Context) {
	var (
		p         = c.Request.Form
		err       error
		cid, dmid int64
		reason    int64
	)
	mid, _ := c.Get("mid")
	if cid, err = strconv.ParseInt(p.Get("cid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if dmid, err = strconv.ParseInt(p.Get("dmid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if reason, err = strconv.ParseInt(p.Get("reason"), 10, 8); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if int8(reason) > model.ReportReasonTeenagers {
		c.JSON(nil, ecode.DMReportReasonError)
		return
	}
	content := p.Get("content")
	if len([]rune(content)) > _maxContentLen {
		c.JSON(nil, ecode.DMReportReasonTooLong)
		return
	}
	if _, err = dmSvc.AddReport(c, cid, dmid, mid.(int64), int8(reason), content); err != nil {
		log.Error("dmSvc.AddReport(cid:%d, dmid:%d, mid:%v) error(%v)", cid, dmid, mid, err)
		c.JSON(nil, err)
		return
	}
	res := map[string]interface{}{}
	res["message"] = ""
	c.JSONMap(res, err)
}

func addReport2(c *bm.Context) {
	var (
		p      = c.Request.Form
		err    error
		cid    int64
		reason int64
		dmids  []int64
		ok2    bool
	)
	mid, _ := c.Get("mid")
	if cid, err = strconv.ParseInt(p.Get("cid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if dmids, err = xstr.SplitInts(p.Get("dmids")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if reason, err = strconv.ParseInt(p.Get("reason"), 10, 8); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if int8(reason) > model.ReportReasonOther {
		c.JSON(nil, ecode.DMReportReasonError)
		return
	}
	content := p.Get("content")
	if len([]rune(content)) > _maxContentLen {
		c.JSON(nil, ecode.DMReportReasonTooLong)
		return
	}
	for _, dmid := range dmids {
		if _, err = dmSvc.AddReport(c, cid, dmid, mid.(int64), int8(reason), content); err != nil {
			log.Error("dmSvc.AddReport(cid:%d, dmid:%d, mid:%v) error(%v)", cid, dmid, mid, err)
		}
		if err == nil {
			ok2 = true
		}
	}
	if ok2 {
		res := map[string]interface{}{}
		res["message"] = ""
		c.JSONMap(res, nil)
		return
	}
	c.JSON(nil, err)
}

func editReport(c *bm.Context) {
	var (
		err           error
		cid, dmid, op int64
		p             = c.Request.Form
	)
	mid, _ := c.Get("mid")
	if dmid, err = strconv.ParseInt(p.Get("dmid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if cid, err = strconv.ParseInt(p.Get("cid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if op, err = strconv.ParseInt(p.Get("op"), 10, 8); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if int8(op) != model.StatUpperIgnore && int8(op) != model.StatUpperDelete {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	_, err = dmSvc.EditReport(c, 1, cid, mid.(int64), dmid, int8(op))
	c.JSON(nil, err)
}

func reportList(c *bm.Context) {
	var (
		err            error
		aid            int64 = -1
		page, pageSize int64
		p              = c.Request.Form
		upOp           = model.StatUpperInit
	)
	mid, _ := c.Get("mid")
	aidStr := p.Get("aid")
	if aidStr != "" {
		if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if page, err = strconv.ParseInt(p.Get("page"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pageSize, err = strconv.ParseInt(p.Get("size"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := dmSvc.ReportList(c, mid.(int64), aid, page, pageSize, upOp, []int64{int64(model.StatFirstInit), int64(model.StatSecondInit)})
	if err != nil {
		log.Error("dmSvc.ReportList(mid:%v, aid:%d) error(%v)", mid, aid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func rptArchives(c *bm.Context) {
	var (
		err    error
		upOp         = model.StatUpperInit
		states       = []int8{model.StatFirstInit, model.StatSecondInit}
		pn, ps int64 = 1, 20
	)
	mid, _ := c.Get("mid")
	data, err := dmSvc.ReportArchives(c, mid.(int64), upOp, states, pn, ps)
	if err != nil {
		log.Error("dmSvc.ReportArchives(mid:%v) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
