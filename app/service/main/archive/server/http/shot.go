package http

import (
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func videoshot(c *bm.Context) {
	params := c.Request.Form
	cidStr := params.Get("cid")
	aidStr := params.Get("aid")
	// check params
	cid, err := strconv.ParseInt(cidStr, 10, 64)
	if err != nil || cid == 0 {
		log.Warn("query (cid) must be number and > 0 but (%s) error(%v)", cidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid == 0 {
		log.Warn("videoshot aid(%s) error", aidStr)
	}
	c.JSON(arcSvc.Videoshot(c, aid, cid))
}
