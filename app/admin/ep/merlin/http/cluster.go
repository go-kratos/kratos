package http

import (
	bm "go-common/library/net/http/blademaster"
)

func queryCluster(c *bm.Context) {
	c.JSON(svc.QueryCluster(c))
}
