package http

import (
	"strconv"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_format = "2006-01-02"
)

var (
	_loc = time.Now().Location()
)

// generate generate invite code in batches.
func generate(c *bm.Context) {
	var (
		err    error
		params = c.Request.Form
	)

	numStr := params.Get("num")
	if numStr == "" {
		httpCode(c, ecode.RequestErr)
		return
	}
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		httpCode(c, ecode.RequestErr)
		return
	}

	midStr := params.Get("mid")
	if midStr == "" {
		httpCode(c, ecode.RequestErr)
		return
	}
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		httpCode(c, ecode.RequestErr)
		return
	}

	expireDayStr := params.Get("expire")
	if expireDayStr == "" {
		httpCode(c, ecode.RequestErr)
		return
	}
	expireDay, err := strconv.ParseInt(expireDayStr, 10, 64)
	if err != nil {
		httpCode(c, ecode.RequestErr)
		return
	}

	invs, err := svc.Generate(c, mid, num, expireDay)
	if err != nil {
		log.Error("service.Generate(%d, %d, %d) error(%v)", mid, num, expireDay, err)
		httpCode(c, err)
		return
	}
	httpData(c, invs, nil)
}

// list get invite codes range from and to.
func list(c *bm.Context) {
	params := c.Request.Form

	midStr := params.Get("mid")
	if midStr == "" {
		httpCode(c, ecode.RequestErr)
		return
	}
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		httpCode(c, ecode.RequestErr)
		return
	}

	fromStr := params.Get("from")
	toStr := params.Get("to")
	if fromStr == "" || toStr == "" {
		httpCode(c, ecode.RequestErr)
		return
	}

	from, to, ok := rangeDate(fromStr, toStr)
	if !ok {
		httpCode(c, ecode.RequestErr)
		return
	}

	fromTs := from.Unix()
	toTs := to.Unix()

	if fromTs < 0 || toTs < 0 || fromTs > toTs {
		httpCode(c, ecode.RequestErr)
		return
	}
	invs, err := svc.List(c, mid, fromTs, toTs)
	if err != nil {
		log.Error("service.List(%d, %d, %d) error(%v)", mid, from, to, err)
		httpCode(c, err)
		return
	}
	httpData(c, invs, nil)
}

func rangeDate(fromStr, toStr string) (from, to time.Time, ok bool) {
	ok = true
	if fromStr == "" || toStr == "" {
		ok = false
		return
	}
	fromDate, err := time.ParseInLocation(_format, fromStr, _loc)
	if err != nil {
		ok = false
		return
	}
	toDate, err := time.ParseInLocation(_format, toStr, _loc)
	if err != nil {
		ok = false
		return
	}
	from = time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, fromDate.Location())
	to = time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 23, 59, 59, 0, toDate.Location())
	return
}
