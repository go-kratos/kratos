package http

import (
	bm "go-common/library/net/http/blademaster"
)

func userNotice(c *bm.Context) {
	midInter, _ := c.Get("mid")
	mid := midInter.(int64)
	c.JSON(artSrv.UserNoticeState(c, mid))
}

func updateUserNotice(c *bm.Context) {
	midInter, _ := c.Get("mid")
	mid := midInter.(int64)
	typ := c.Request.Form.Get("type")
	c.JSON(nil, artSrv.UpdateUserNoticeState(c, mid, typ))
}
