package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func elecShow(c *bm.Context) {
	var (
		mid, loginID, aid int64
		err               error
	)

	params := c.Request.Form
	midStr := params.Get("mid")
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil || mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aidStr := params.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// login mid
	if loginIDStr, ok := c.Get("mid"); ok {
		loginID = loginIDStr.(int64)
	}
	c.JSON(webSvc.ElecShow(c, mid, aid, loginID))
}
