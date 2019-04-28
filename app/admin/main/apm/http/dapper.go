package http

import (
	bm "go-common/library/net/http/blademaster"
)

func dapperProxy(c *bm.Context) {
	apmSvc.DapperProxy(c.Writer, c.Request)
}
