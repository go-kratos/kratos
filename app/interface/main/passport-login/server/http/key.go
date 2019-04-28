package http

import (
	bm "go-common/library/net/http/blademaster"
)

func getKey(c *bm.Context) {
	c.JSON(srv.RSAKey(c), nil)
}
