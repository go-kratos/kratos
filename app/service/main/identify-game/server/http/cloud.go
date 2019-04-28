package http

import (
	bm "go-common/library/net/http/blademaster"
)

func regions(c *bm.Context) {
	c.JSON(srv.Regions, nil)
}
