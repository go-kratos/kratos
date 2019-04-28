package http

import (
	"strconv"

	"go-common/app/admin/main/dm/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func archiveList(c *bm.Context) {
	var (
		p   = c.Request.Form
		req = &model.ArchiveListReq{
			Pn:     1,
			Ps:     20,
			IDType: p.Get("type"),
			Sort:   "desc",
			Order:  "mtime",
			Page:   int64(model.CondIntNil),
			Attrs:  make([]int64, 0),
			State:  int64(model.CondIntNil),
		}
		err error
	)
	if idStr := p.Get("id"); len(idStr) > 0 {
		if req.ID, err = strconv.ParseInt(idStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if pageStr := p.Get("page"); len(pageStr) > 0 {
		if req.Page, err = strconv.ParseInt(pageStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if attrStr := p.Get("attrs"); len(attrStr) > 0 {
		req.Attrs, err = xstr.SplitInts(attrStr)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if stateStr := p.Get("state"); len(stateStr) > 0 {
		req.State, err = strconv.ParseInt(stateStr, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if p.Get("sort") != "" {
		req.Sort = p.Get("sort")
	}
	if p.Get("order") != "" {
		req.Order = p.Get("order")
	}
	if pnStr := p.Get("pn"); len(pnStr) > 0 {
		req.Pn, err = strconv.ParseInt(pnStr, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if psStr := p.Get("ps"); len(psStr) > 0 {
		req.Ps, err = strconv.ParseInt(psStr, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	data, err := dmSvc.ArchiveList(c, req)
	c.JSON(data, err)
}

func uptSubjectsState(c *bm.Context) {
	var (
		uid, _   = c.Get("uid")
		uname, _ = c.Get("username")
		p        = c.Request.Form
		comment  = p.Get("comment")
	)
	oids, err := xstr.SplitInts(p.Get("oids"))
	if err != nil || len(oids) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	state, err := strconv.ParseInt(p.Get("state"), 10, 64)
	if err != nil || (int32(state) != model.SubStateOpen && int32(state) != model.SubStateClosed) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.UptSubjectsState(c, int32(tp), uid.(int64), uname.(string), oids, int32(state), comment)
	c.JSON(nil, err)
}

func upSubjectMaxLimit(c *bm.Context) {
	var (
		tp            int64
		p             = c.Request.Form
		cid, maxlimit int64
		err           error
	)
	if tp, err = strconv.ParseInt(p.Get("type"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if cid, err = strconv.ParseInt(p.Get("cid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if maxlimit, err = strconv.ParseInt(p.Get("limit"), 10, 64); err != nil || maxlimit > 20000 || maxlimit < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.UpSubjectMaxLimit(c, int32(tp), cid, maxlimit)
	c.JSON(nil, err)
}

func subjectLog(c *bm.Context) {
	var (
		p       = c.Request.Form
		oid, tp int64
		err     error
	)
	if oid, err = strconv.ParseInt(p.Get("oid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tp, err = strconv.ParseInt(p.Get("type"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := dmSvc.SubjectLog(c, int32(tp), oid)
	c.JSON(data, err)
}
