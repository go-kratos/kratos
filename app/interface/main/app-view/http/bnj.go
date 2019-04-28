package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// viewIndex view handler
func bnj2019(c *bm.Context) {
	var (
		mid      int64
		relateID int64
		params   = c.Request.Form
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if !viewSvr.CheckAccess(mid) {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	relateID, _ = strconv.ParseInt(params.Get("relate_id"), 10, 64)
	c.JSON(viewSvr.Bnj2019(c, mid, relateID))
}

// viewPage view page handler.
func bnjList(c *bm.Context) {
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if !viewSvr.CheckAccess(mid) {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	c.JSON(viewSvr.BnjList(c, mid))
}

// videoShot video shot .
func bnjItem(c *bm.Context) {
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if !viewSvr.CheckAccess(mid) {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	params := c.Request.Form
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(viewSvr.BnjItem(c, aid, mid))
}
