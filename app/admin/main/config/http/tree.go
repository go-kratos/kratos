package http

import (
	bm "go-common/library/net/http/blademaster"
)

func syncTree(c *bm.Context) {
	svr.SyncTree(c, user(c), c.Request.Header.Get("Cookie"))
	c.JSON(nil, nil)
}
