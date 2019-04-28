package http

import bm "go-common/library/net/http/blademaster"

func domain(c *bm.Context) {
	c.JSON(domainSvc.Domain(), nil)
}
