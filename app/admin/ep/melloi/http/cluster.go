package http

import (
	bm "go-common/library/net/http/blademaster"
)

//ClusterInfo get cluster infomation
func ClusterInfo(c *bm.Context) {
	c.JSON(srv.ClusterInfo(c))
}
