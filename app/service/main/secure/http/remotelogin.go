package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func addlog(c *bm.Context) {
	res := c.Request
	mid := res.Form.Get("mid")
	ip := res.Form.Get("ip")
	c.JSON(nil, rSrv.AddLog(c, mid, ip))
}

func getLoc(c *bm.Context) {
	var (
		mid    int64
		params = c.Request.Form
		midStr = params.Get("mid")
	)
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil || mid < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(rSrv.GetLoc(c, mid))
}

func status(c *bm.Context) {
	var (
		mid    int64
		params = c.Request.Form
		midStr = params.Get("mid")
		uuid   = params.Get("uuid")
	)
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil || mid < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(rSrv.Status(c, mid, uuid))
}

func closeNotify(c *bm.Context) {
	var (
		mid    int64
		params = c.Request.Form
		midStr = params.Get("mid")
		uuid   = params.Get("uuid")
	)
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil || mid < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, rSrv.CloseNotify(c, mid, uuid))
}

func loc(c *bm.Context) {
	var (
		mid    int64
		params = c.Request.Form
		midStr = params.Get("mid")
	)
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil || mid < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(rSrv.ExpectionLoc(c, mid))
}

func feedback(c *bm.Context) {
	var (
		mid    int64
		params = c.Request.Form
		midStr = params.Get("mid")
		ip     = params.Get("ip")
		tsStr  = params.Get("logintime")
		tpStr  = params.Get("tp")
	)
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil || mid < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, rSrv.AddFeedBack(c, mid, ts, int8(tp), ip))
}

func oftenCheck(c *bm.Context) {
	var (
		mid    int64
		params = c.Request.Form
		midStr = params.Get("mid")
		ip     = params.Get("ip")
	)
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil || mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	result, err := rSrv.OftenCheck(c, mid, ip)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(result, nil)
}
