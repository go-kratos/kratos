package http

import (
	"strings"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func discoveryProxy(c *bm.Context) {
	req := c.Request
	idx := strings.Index(req.URL.Path, "/x/admin/apm/discovery/")
	if idx == -1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := apmSvc.DiscoveryProxy(c, req.Method, req.URL.Path[idx+len("/x/admin/apm/discovery/"):], req.Form)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func discoveryProxyNoAuth(c *bm.Context) {
	req := c.Request
	data, err := apmSvc.DiscoveryProxy(c, req.Method, "fetch", req.Form)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
