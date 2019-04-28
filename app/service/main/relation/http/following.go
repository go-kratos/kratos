package http

import (
	"strconv"

	"go-common/app/service/main/relation/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/time"
)

// followings get user's following list.
func followings(c *bm.Context) {
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
	c.JSON(relationSvc.Followings(c, mid))
}

// addFollowing add following.
func addFollowing(c *bm.Context) {
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
	c.JSON(nil, relationSvc.AddFollowing(c, mid, fid, uint8(src), ric))
}

// delFollowing del following.
func delFollowing(c *bm.Context) {
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
	c.JSON(nil, relationSvc.DelFollowing(c, mid, fid, uint8(src), ric))
}

// delFollowingCache del following cache.
func delFollowingCache(c *bm.Context) {
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
	c.JSON(nil, relationSvc.DelFollowingCache(c, mid))
}

// updateFollowingCache update following cache.
func updateFollowingCache(c *bm.Context) {
	var (
		err       error
		mid, fid  int64
		na        uint64
		mts       int64
		following *model.Following
		params    = c.Request.Form
		midStr    = params.Get("mid")
		fidStr    = params.Get("fid")
		naStr     = params.Get("attribute")
		mtStr     = params.Get("mtime")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if fid, err = strconv.ParseInt(fidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if na, err = strconv.ParseUint(naStr, 10, 32); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mts, err = strconv.ParseInt(mtStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	following = &model.Following{
		Mid:       fid,
		Attribute: uint32(na),
		MTime:     time.Time(mts),
	}
	c.JSON(nil, relationSvc.UpdateFollowingCache(c, mid, following))
}

func sameFollowings(c *bm.Context) {
	arg := new(model.ArgSameFollowing)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(relationSvc.SameFollowings(c, arg))
}
