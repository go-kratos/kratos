package http

import (
	bm "go-common/library/net/http/blademaster"
)

func anniversaryInfo(c *bm.Context) {
	var (
		mid int64
	)
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	c.JSON(artSrv.AnniversaryInfo(c, mid))
}
