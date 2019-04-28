package http

import (
	bm "go-common/library/net/http/blademaster"
)

func broadServer(c *bm.Context) {
	v := new(struct {
		Platform string `form:"platform" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(webSvc.BroadServers(c, v.Platform))
}
