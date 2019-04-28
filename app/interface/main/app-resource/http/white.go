package http

import bm "go-common/library/net/http/blademaster"

func whiteList(c *bm.Context) {
	c.JSON(whiteSvc.List(), nil)
}
