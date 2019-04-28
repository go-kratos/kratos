package http

import (
	"time"

	bm "go-common/library/net/http/blademaster"
)

func displayID(c *bm.Context) {
	header := c.Request.Header
	buvid := header.Get("Buvid")
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	id := displaySvr.DisplayID(c, mid, buvid, time.Now())

	c.JSON(struct {
		ID string `json:"id"`
	}{ID: id}, nil)
}

func zone(c *bm.Context) {
	c.JSON(displaySvr.Zone(c, time.Now()), nil)
}
