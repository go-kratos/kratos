package http

import (
	bm "go-common/library/net/http/blademaster"
)

func onlineInfo(c *bm.Context) {
	c.JSON(webSvc.OnlineArchiveCount(c), nil)
}

func onlineList(c *bm.Context) {
	c.JSON(webSvc.OnlineList(c))
}
