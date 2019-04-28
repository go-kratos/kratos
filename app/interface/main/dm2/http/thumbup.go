package http

import (
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func thumbupDM(c *bm.Context) {
	p := c.Request.Form
	mid, _ := c.Get("mid")
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil || oid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	dmid, err := strconv.ParseInt(p.Get("dmid"), 10, 64)
	if err != nil || dmid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	op, err := strconv.ParseInt(p.Get("op"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = dmSvc.ThumbupDM(c, oid, dmid, mid.(int64), int8(op)); err != nil {
		log.Error("dmSvc.ThumbupDM(oid:%d,dmid:%d) error(%v)", oid, dmid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func thumbupStats(c *bm.Context) {
	var (
		mid int64
		p   = c.Request.Form
	)
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil || oid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midI, ok := c.Get("mid"); ok {
		mid = midI.(int64)
	}
	dmids, err := xstr.SplitInts(p.Get("ids"))
	if err != nil || len(dmids) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := dmSvc.ThumbupList(c, oid, mid, dmids)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
