package http

import (
	bm "go-common/library/net/http/blademaster"
)

func shopInfo(c *bm.Context) {
	var (
		mid int64
	)
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	c.JSON(spcSvc.ShopInfo(c, mid))
}
