package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func coinVideo(c *bm.Context) {
	var (
		mid, vmid int64
		err       error
	)
	vmidStr := c.Request.Form.Get("vmid")
	if vmid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || vmid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(spcSvc.CoinVideo(c, mid, vmid))
}
