package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// whispers get user's whisper list.
func whispers(c *bm.Context) {
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
	c.JSON(relationSvc.Whispers(c, mid))
}

// addWhisper add whisper.
func addWhisper(c *bm.Context) {
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
	c.JSON(nil, relationSvc.AddWhisper(c, mid, fid, uint8(src), ric))
}

// delWhisper del whisper.
func delWhisper(c *bm.Context) {
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
	c.JSON(nil, relationSvc.DelWhisper(c, mid, fid, uint8(src), ric))
}
