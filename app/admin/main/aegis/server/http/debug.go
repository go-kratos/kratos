package http

import (
	bm "go-common/library/net/http/blademaster"
)

func cache(c *bm.Context) {
	c.JSONMap(srv.Cache(), nil)
}
