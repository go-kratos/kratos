package http

import (
	"strings"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func searchProxyGet(c *bm.Context) {
	req := c.Request
	idx := strings.Index(req.URL.Path, "/x/admin/apm/platform/search/get/")
	if idx == -1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := apmSvc.PlatformProxy(c, "GET", req.URL.Path[idx+len("/x/admin/apm/platform/search/get/"):], req.Form)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func searchProxyPost(c *bm.Context) {
	req := c.Request
	idx := strings.Index(req.URL.Path, "/x/admin/apm/platform/search/post/")
	if idx == -1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := apmSvc.PlatformProxy(c, "POST", req.URL.Path[idx+len("/x/admin/apm/platform/search/post/"):], req.Form)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func replyProxyGet(c *bm.Context) {
	req := c.Request
	idx := strings.Index(req.URL.Path, "/x/admin/apm/platform/reply/get/")
	if idx == -1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := apmSvc.PlatformProxy(c, "GET", req.URL.Path[idx+len("/x/admin/apm/platform/reply/get/"):], req.Form)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func replyProxyPost(c *bm.Context) {
	req := c.Request
	idx := strings.Index(req.URL.Path, "/x/admin/apm/platform/reply/post/")
	if idx == -1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := apmSvc.PlatformProxy(c, "POST", req.URL.Path[idx+len("/x/admin/apm/platform/reply/post/"):], req.Form)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func tagProxyGet(c *bm.Context) {
	req := c.Request
	idx := strings.Index(req.URL.Path, "/x/admin/apm/platform/tag/get/")
	if idx == -1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := apmSvc.PlatformProxy(c, "GET", req.URL.Path[idx+len("/x/admin/apm/platform/tag/get/"):], req.Form)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func tagProxyPost(c *bm.Context) {
	req := c.Request
	idx := strings.Index(req.URL.Path, "/x/admin/apm/platform/tag/post/")
	if idx == -1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := apmSvc.PlatformProxy(c, "POST", req.URL.Path[idx+len("/x/admin/apm/platform/tag/post/"):], req.Form)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func bfsProxyGet(c *bm.Context) {
	req := c.Request
	idx := strings.Index(req.URL.Path, "/x/admin/apm/platform/bfs/get/")
	if idx == -1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := apmSvc.PlatformProxy(c, "GET", req.URL.Path[idx+len("/x/admin/apm/platform/bfs/get/"):], req.Form)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func bfsProxyPost(c *bm.Context) {
	req := c.Request
	idx := strings.Index(req.URL.Path, "/x/admin/apm/platform/bfs/post/")
	if idx == -1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := apmSvc.PlatformProxy(c, "POST", req.URL.Path[idx+len("/x/admin/apm/platform/bfs/post/"):], req.Form)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func replyFeedProxyGet(c *bm.Context) {
	req := c.Request
	idx := strings.Index(req.URL.Path, "/x/admin/apm/platform/reply/feed/get/")
	if idx == -1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := apmSvc.PlatformProxy(c, "GET", req.URL.Path[idx+len("/x/admin/apm/platform/reply/feed/get/"):], req.Form)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func replyFeedProxyPost(c *bm.Context) {
	req := c.Request
	idx := strings.Index(req.URL.Path, "/x/admin/apm/platform/reply/feed/post/")
	if idx == -1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := apmSvc.PlatformProxy(c, "POST", req.URL.Path[idx+len("/x/admin/apm/platform/reply/feed/post/"):], req.Form)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
