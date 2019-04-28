package http

import (
	"go-common/library/net/http/blademaster"
)

func regions(c *blademaster.Context) {
	c.JSON(srv.Regions(c), nil)
}
