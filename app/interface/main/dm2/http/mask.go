package http

import (
	"strconv"
	"strings"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func updateMask(c *bm.Context) {
	p := c.Request.Form
	var (
		cid, plat, fps, maskTime int64
		list                     string
		err                      error
	)
	if cid, err = strconv.ParseInt(p.Get("cid"), 10, 64); err != nil || cid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if plat, err = strconv.ParseInt(p.Get("plat"), 10, 64); err != nil || plat < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if p.Get("time") != "" {
		if maskTime, err = strconv.ParseInt(p.Get("time"), 10, 64); err != nil || maskTime < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if p.Get("fps") != "" {
		fps, err = strconv.ParseInt(p.Get("fps"), 10, 64)
		if err != nil || fps <= 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	list = strings.Replace(p.Get("list"), " ", "", -1)
	err = dmSvc.UpdateMask(c, cid, maskTime, int32(fps), int8(plat), list)
	c.JSON(nil, err)
}
