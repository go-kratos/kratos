package http

import (
	bm "go-common/library/net/http/blademaster"
)

func debugCache(c *bm.Context) {
	opt := new(struct {
		Keys string `form:"keys" validate:"required"`
	})
	if err := c.Bind(opt); err != nil {
		return
	}
	c.JSONMap(srv.DebugCache(opt.Keys), nil)
}
