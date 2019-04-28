package http

import (
	"go-common/library/ecode"
	"strconv"

	bm "go-common/library/net/http/blademaster"
)

// Blacks get user's black list.
func blacks(c *bm.Context) {
	var (
		err    error
		mid    int64
		params = c.Request.Form
		midStr = params.Get("mid")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(relationSvc.Blacks(c, mid))
}

// addBlack add black.
func addBlack(c *bm.Context) {
	var (
		err      error
		mid, fid int64
		src      uint64
		params   = c.Request.Form
		midStr   = params.Get("mid")
		fidStr   = params.Get("fid")
		srcStr   = params.Get("src")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if fid, err = strconv.ParseInt(fidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mid <= 0 || fid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if src, err = strconv.ParseUint(srcStr, 10, 8); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ric := infocArg(c)
	c.JSON(nil, relationSvc.AddBlack(c, mid, fid, uint8(src), ric))
}

// delBlack del black.
func delBlack(c *bm.Context) {
	var (
		err      error
		mid, fid int64
		src      uint64
		params   = c.Request.Form
		midStr   = params.Get("mid")
		fidStr   = params.Get("fid")
		srcStr   = params.Get("src")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if fid, err = strconv.ParseInt(fidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if src, err = strconv.ParseUint(srcStr, 10, 8); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ric := infocArg(c)
	c.JSON(nil, relationSvc.DelBlack(c, mid, fid, uint8(src), ric))
}
