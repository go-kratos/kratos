package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func replyRecord(c *bm.Context) {
	var (
		err               error
		types             []int64
		mid, stime, etime int64
		pn                = int64(1)
		ps                = int64(10)
	)
	params := c.Request.Form
	stimeStr := params.Get("stime")
	etimeStr := params.Get("etime")
	typesStr := params.Get("types")
	sortStr := params.Get("sort")
	orderStr := params.Get("order")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	midStr := params.Get("mid")
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if typesStr != "" {
		if types, err = xstr.SplitInts(typesStr); err != nil {
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	if stime, err = strconv.ParseInt(stimeStr, 10, 64); err != nil {
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if etime, err = strconv.ParseInt(etimeStr, 10, 64); err != nil {
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if pnStr != "" {
		if pn, err = strconv.ParseInt(pnStr, 10, 32); err != nil {
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	if psStr != "" {
		if ps, err = strconv.ParseInt(psStr, 10, 32); err != nil {
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	if orderStr == "" {
		orderStr = "ctime"
	}
	if sortStr == "" {
		sortStr = "desc"
	}
	records, total, err := rpSvr.Records(c, types, mid, stime, etime, orderStr, sortStr, int32(pn), int32(ps))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{
		"page": map[string]int64{
			"num":   pn,
			"size":  ps,
			"total": int64(total),
		},
		"records": records,
	}
	c.JSON(data, err)
}
