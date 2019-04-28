package http

import (
	"strconv"

	bm "go-common/library/net/http/blademaster"
)

func notice(c *bm.Context) {
	params := c.Request.Form
	build, _ := strconv.Atoi(params.Get("build"))
	c.JSON(artSrv.Notice(params.Get("platform"), build), nil)
}
