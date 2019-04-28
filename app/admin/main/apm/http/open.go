package http

import (
	"strings"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func openProxyGet(c *bm.Context) {
	req := c.Request
	idx := strings.Index(req.URL.Path, "/x/admin/apm/open/get/")
	if idx == -1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := apmSvc.OpenProxy(c, "GET", req.URL.Path[idx+len("/x/admin/apm/open/get/"):], req.Form)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func openProxyPost(c *bm.Context) {
	req := c.Request
	idx := strings.Index(req.URL.Path, "/x/admin/apm/open/post/")
	if idx == -1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := apmSvc.OpenProxy(c, "POST", req.URL.Path[idx+len("/x/admin/apm/open/post/"):], req.Form)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
