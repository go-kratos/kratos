package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const (
	_defaultLimit = 1
)

func loginLog(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	if midStr == "" {
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

	limit := _defaultLimit
	limitStr := params.Get("limit")
	if limitStr != "" {
		if limit, err = strconv.Atoi(limitStr); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	c.JSON(passportSvc.FormattedLoginLogs(c, mid, limit))
}
