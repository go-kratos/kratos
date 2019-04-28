package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func face(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	if midStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fromStr, toStr := params.Get("from"), params.Get("to")
	if fromStr == "" || toStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mid < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	from, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	to, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(passportSvc.FaceApplies(c, mid, from, to, params.Get("status"), params.Get("operator")))
}
