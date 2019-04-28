package http

import (
	"strconv"

	"go-common/app/interface/main/creative/model/search"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func replyList(c *bm.Context) {
	req := c.Request
	params := req.Form
	kw := params.Get("keyword")
	order := params.Get("order")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	var (
		err          error
		oid          int64
		isReport, tp int
	)
	isReportStr := params.Get("is_report")
	if isReportStr != "" {
		isReport, err = strconv.Atoi(isReportStr)
		if err != nil {
			log.Error("strconv.Atoi replyList isReportStr(%s)|error(%v)", isReportStr, err)
			c.JSON(nil, ecode.RequestErr)
		}
	}
	oidStr := params.Get("oid")
	if oidStr != "" {
		oid, err = strconv.ParseInt(oidStr, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt replyList oidStr(%s)|error(%v)", oidStr, err)
			c.JSON(nil, ecode.RequestErr)
		}
	}
	typeStr := params.Get("type")
	if typeStr != "" {
		tp, err = strconv.Atoi(typeStr)
		if err != nil {
			log.Error("strconv.ParseInt replyList typeStr(%s)|error(%v)", typeStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	filterStr := params.Get("filter")
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps <= 0 || pn > 10 {
		ps = 10
	}
	tmidStr := params.Get("tmid")
	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}
	p := &search.ReplyParam{
		Ak:          params.Get("access_key"),
		Ck:          c.Request.Header.Get("cookie"),
		OMID:        mid,
		OID:         oid,
		Pn:          pn,
		Ps:          ps,
		IP:          metadata.String(c, metadata.RemoteIP),
		IsReport:    int8(isReport),
		Type:        int8(tp),
		FilterCtime: filterStr,
		Kw:          kw,
		Order:       order,
	}
	replies, err := replySvc.Replies(c, p)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{
		"data": replies.Result,
		"pager": map[string]int{
			"current": p.Pn,
			"size":    p.Ps,
			"total":   replies.Total,
		},
	}, nil)
}
