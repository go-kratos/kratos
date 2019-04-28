package http

import (
	"strconv"
	"strings"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// relation get relation between mid and fid.
func relation(c *bm.Context) {
	var (
		err      error
		mid, fid int64
		params   = c.Request.Form
		midStr   = params.Get("mid")
		fidStr   = params.Get("fid")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if fid, err = strconv.ParseInt(fidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(relationSvc.Relation(c, mid, fid))
}

// relations get relations between mid and fids.
func relations(c *bm.Context) {
	var (
		err      error
		mid, fid int64
		fids     []int64
		params   = c.Request.Form
		midStr   = params.Get("mid")
		fidsStr  = params.Get("fids")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fidsStrArr := strings.Split(fidsStr, ",")
	for _, v := range fidsStrArr {
		if fid, err = strconv.ParseInt(v, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		fids = append(fids, fid)
	}
	c.JSON(relationSvc.Relations(c, mid, fids))
}
