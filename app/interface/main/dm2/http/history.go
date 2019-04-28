package http

import (
	"net/http"
	"strconv"
	"time"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func dmHistory(c *bm.Context) {
	var (
		p           = c.Request.Form
		contextType = "text/xml"
	)
	tp, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	date, err := time.Parse("2006-01-02", p.Get("date"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// convert 2006-01-02-->2016-01-02 23:59:59
	tm := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.Local)
	data, err := dmSvc.SearchDMHistory(c, int32(tp), oid, tm.Unix())
	if err != nil {
		c.AbortWithStatus(httpCode(err))
		return
	}
	c.Writer.Header().Set("Content-Encoding", "deflate")
	c.Bytes(200, contextType, data)
}

func dmHistoryV2(c *bm.Context) {
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
	date, err := time.Parse("2006-01-02", p.Get("date"))
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// convert 2006-01-02-->2016-01-02 23:59:59
	tm := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.Local)
	c.JSON(dmSvc.SearchDMHistoryV2(c, int32(tp), oid, tm.Unix()))
}

func dmHistoryIndex(c *bm.Context) {
	var (
		p   = c.Request.Form
		now = time.Now()
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
	month := p.Get("month")
	date, err := time.Parse("2006-01", month)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// only allow recent one year query
	if now.Year()-date.Year() >= 1 && now.Month()-date.Month() > 12 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := dmSvc.SearchDMHisIndex(c, int32(tp), oid, month)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
