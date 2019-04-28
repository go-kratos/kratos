package http

import (
	bm "go-common/library/net/http/blademaster"
)

func indexIcon(c *bm.Context) {
	c.JSON(webSvc.IndexIcon(), nil)
}
