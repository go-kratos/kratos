package http

import (
	bm "go-common/library/net/http/blademaster"
)

func ipZone(c *bm.Context) {
	c.JSON(webSvc.IPZone(c))
}
