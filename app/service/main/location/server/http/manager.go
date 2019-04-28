package http

import (
	bm "go-common/library/net/http/blademaster"
)

func flushCache(c *bm.Context) {
	if err := svr.LoadPolicy(); err != nil {
		c.JSON(nil, err)
	}
}
