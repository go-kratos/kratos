package http

import (
	"net/url"
	"strconv"
	"strings"

	"go-common/app/admin/main/dm/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// checkState check admin operation.
func checkState(state int8) (ok bool) {
	if state != model.StatFirstInit &&
		state != model.StatFirstDelete &&
		state != model.StatFirstIgnore &&
		state != model.StatSecondInit &&
		state != model.StatSecondIgnore &&
		state != model.StatSecondAutoDelete &&
		state != model.StatSecondDelete {
		ok = false
	} else {
		ok = true
	}
	return
}

func reportList2(c *bm.Context) {
	var (
		v = new(model.ReportListParams)
	)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(dmSvc.ReportList2(c, v))
}

func changeReportStat(c *bm.Context) {
	var (
		reason, notice, block, blockReason, moral int64
		cidDmids                                  = map[int64][]int64{}
		params                                    = c.Request.Form
		data                                      struct {
			Affect int64 `json:"affect"`
		}
	)
	uid, err := strconv.ParseInt(params.Get("adminId"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	state, err := strconv.ParseInt(params.Get("state"), 10, 8)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	uname := params.Get("uname")
	remark := params.Get("remark")
	noticeStr := params.Get("notice")
	if noticeStr != "" {
		if notice, err = strconv.ParseInt(noticeStr, 10, 8); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	ids := strings.Split(params.Get("ids"), "|")
	if len(ids) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, idStr := range ids {
		var (
			cid   int64
			dmids []int64
		)
		s := strings.Split(idStr, ":")
		if len(s) != 2 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if cid, err = strconv.ParseInt(s[0], 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if dmids, err = xstr.SplitInts(s[1]); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if !checkState(int8(state)) {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		cidDmids[cid] = dmids
	}
	if state == int64(model.StatSecondDelete) || state == int64(model.StatFirstDelete) {
		blockStr := params.Get("block")
		if blockStr != "" {
			if block, err = strconv.ParseInt(blockStr, 10, 8); err != nil {
				c.JSON(nil, ecode.RequestErr)
				return
			}
		}
		MoralStr := params.Get("moral")
		if MoralStr != "" {
			if moral, err = strconv.ParseInt(MoralStr, 10, 8); err != nil {
				c.JSON(nil, ecode.RequestErr)
				return
			}
		}
		blockReasonStr := params.Get("block_reason")
		if blockReasonStr != "" {
			if blockReason, err = strconv.ParseInt(blockReasonStr, 10, 8); err != nil {
				c.JSON(nil, ecode.RequestErr)
				return
			}
		}
		reasonStr := params.Get("reason")
		if reasonStr != "" {
			if reason, err = strconv.ParseInt(reasonStr, 10, 8); err != nil {
				c.JSON(nil, ecode.RequestErr)
				return
			}
		}
	}
	data.Affect, err = dmSvc.ChangeReportStat(c, cidDmids, int8(state), int8(reason), int8(notice), uid, block, blockReason, moral, remark, uname)
	if err != nil {
		log.Error("dmSvc.ChangeReportStat(id:%+v, uid:%d) error(%v)", cidDmids, uid, err)
		c.JSON(nil, err)
	}
	res := map[string]interface{}{}
	res["data"] = data
	c.JSONMap(res, err)
}

func reportList(c *bm.Context) {
	var (
		tid, rpID                        []int64
		rt                               *model.Report
		p                                = c.Request.Form
		start, end, sort, order, keyword string
	)
	rt = &model.Report{
		Aid:    -1,
		UID:    -1,
		RpUID:  -1,
		RpType: -1,
		Cid:    -1,
	}
	state, err := xstr.SplitInts(p.Get("state"))
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	upOp, err := xstr.SplitInts(p.Get("up_op"))
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	page, err := strconv.ParseInt(p.Get("page"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tidStr := p.Get("tid")
	if tidStr != "" {
		if tid, err = xstr.SplitInts(tidStr); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	aidStr := p.Get("aid")
	if aidStr != "" {
		if rt.Aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	cidStr := p.Get("cid")
	if cidStr != "" {
		if rt.Cid, err = strconv.ParseInt(cidStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	uidStr := p.Get("uid")
	if uidStr != "" {
		if rt.UID, err = strconv.ParseInt(uidStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	userStr := p.Get("rp_user")
	if userStr != "" {
		if rt.RpUID, err = strconv.ParseInt(userStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	typeStr := p.Get("rp_type")
	if typeStr != "" {
		if rpID, err = xstr.SplitInts(typeStr); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	startStr := p.Get("start")
	start, err = url.QueryUnescape(startStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	endStr := p.Get("end")
	end, err = url.QueryUnescape(endStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pageSizeStr := p.Get("page_size")
	var pageSize int64
	if pageSizeStr != "" {
		if pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if pageSize > 100 {
			pageSize = 100
		}
	} else {
		pageSize = 100
	}
	// TODO: swap order&sort
	order = p.Get("sort")
	sort = p.Get("order")
	keyword = p.Get("keyword")
	rpts, err := dmSvc.ReportList(c, page, pageSize, start, end, order, sort, keyword, tid, rpID, state, upOp, rt)
	res := map[string]interface{}{}
	res["data"] = rpts
	c.JSONMap(res, err)
}

func reportLog(c *bm.Context) {
	p := c.Request.Form
	dmid, err := strconv.ParseInt(p.Get("dmid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := dmSvc.ReportLog(c, dmid)
	res := map[string]interface{}{}
	res["data"] = data
	c.JSONMap(res, err)
}

func changeReportUserStat(c *bm.Context) {
	p := c.Request.Form
	dmids, err := xstr.SplitInts(p.Get("dmids"))
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.ChangeReportUserStat(c, dmids)
	c.JSON(nil, err)
}

func transferJudge(c *bm.Context) {
	var (
		err      error
		uname    string
		cidDmids = map[int64][]int64{}
		p        = c.Request.Form
	)

	ids := strings.Split(p.Get("ids"), "|")
	if len(ids) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, idStr := range ids {
		var (
			cid   int64
			dmids []int64
		)
		s := strings.Split(idStr, ":")
		if len(s) != 2 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if cid, err = strconv.ParseInt(s[0], 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if dmids, err = xstr.SplitInts(s[1]); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		cidDmids[cid] = dmids
	}
	uname = p.Get("uname")
	uid, err := strconv.ParseInt(p.Get("adminId"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.DMReportJudge(c, cidDmids, uid, uname)
	c.JSON(nil, err)
}

// JudgeResult post judgement result
func JudgeResult(c *bm.Context) {
	p := c.Request.Form
	cid, err := strconv.ParseInt(p.Get("cid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	dmid, err := strconv.ParseInt(p.Get("dmid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	result, err := strconv.ParseInt(p.Get("result"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.JudgeResult(c, cid, dmid, result)
	c.JSON(nil, err)
}

func logList(c *bm.Context) {
	p := c.Request.Form
	dmid, err := strconv.ParseInt(p.Get("dmid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := dmSvc.QueryOpLogs(c, dmid)
	res := map[string]interface{}{}
	res["data"] = data
	c.JSONMap(res, err)
}
