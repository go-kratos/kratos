package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func maskState(c *bm.Context) {
	var (
		p       = c.Request.Form
		oid, tp int64
		err     error
	)
	if oid, err = strconv.ParseInt(p.Get("oid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tp, err = strconv.ParseInt(p.Get("type"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	open, mobile, web, err := dmSvc.MaskState(c, int32(tp), oid)
	res := map[string]interface{}{}
	res["open"] = open
	res["mobile"] = mobile
	res["web"] = web
	c.JSON(res, err)
}

func updateMaskState(c *bm.Context) {
	var (
		p                    = c.Request.Form
		oid, tp, plat, state int64
		err                  error
	)
	if oid, err = strconv.ParseInt(p.Get("oid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tp, err = strconv.ParseInt(p.Get("type"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if plat, err = strconv.ParseInt(p.Get("plat"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if state, err = strconv.ParseInt(p.Get("state"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.UpdateMaskState(c, int32(tp), oid, int8(plat), int32(state))
	c.JSON(nil, err)
}

func generateMask(c *bm.Context) {
	var (
		p             = c.Request.Form
		oid, tp, plat int64
		err           error
	)
	if oid, err = strconv.ParseInt(p.Get("oid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tp, err = strconv.ParseInt(p.Get("type"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if plat, err = strconv.ParseInt(p.Get("plat"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.GenerateMask(c, int32(tp), oid, int8(plat))
	c.JSON(nil, err)
}

func maskUps(c *bm.Context) {
	var (
		p   = c.Request.Form
		pn  = int64(1)
		ps  = int64(50)
		err error
	)
	if p.Get("pn") != "" {
		if pn, err = strconv.ParseInt(p.Get("pn"), 10, 64); err != nil || pn <= 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if p.Get("ps") != "" {
		if ps, err = strconv.ParseInt(p.Get("ps"), 10, 64); err != nil || ps <= 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	c.JSON(dmSvc.MaskUps(c, pn, ps))
}

func maskUpOpen(c *bm.Context) {
	var (
		p       = c.Request.Form
		comment = p.Get("comment")
		mids    []int64
		state   int64
		err     error
	)
	if mids, err = xstr.SplitInts(p.Get("mids")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if state, err = strconv.ParseInt(p.Get("state"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.MaskUpOpen(c, mids, int32(state), comment)
	c.JSON(nil, err)
}
