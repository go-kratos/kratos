package http

import (
	bm "go-common/library/net/http/blademaster"
)

// join
func join(c *bm.Context) {
	c.JSON(jobSvc.Jobs(c), nil)
}
