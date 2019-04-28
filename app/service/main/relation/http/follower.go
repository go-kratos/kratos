package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// followers get user's follower list.
func followers(c *bm.Context) {
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
	c.JSON(relationSvc.Followers(c, mid))
}

// delFollower del follower.
func delFollower(c *bm.Context) {
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
	// mid移除粉丝fid：对fid对mid的关注状态进行更改
	c.JSON(nil, relationSvc.DelFollower(c, fid, mid, uint8(src), ric))
}

// delFollowerCache del follower cache.
func delFollowerCache(c *bm.Context) {
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
	c.JSON(nil, relationSvc.DelFollowerCache(c, mid))
}
