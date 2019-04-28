package proxy_test

import (
	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/proxy"
)

// This example create several reverse proxy to show how to use proxy middleware.
// We proxy three path to `api.bilibili.com` and return response without any changes.
func Example() {
	proxies := map[string]string{
		"/index":        "http://api.bilibili.com/html/index",
		"/ping":         "http://api.bilibili.com/api/ping",
		"/api/versions": "http://api.bilibili.com/api/web/versions",
	}

	engine := blademaster.Default()
	for path, ep := range proxies {
		engine.GET(path, proxy.NewAlways(ep))
	}
	engine.Run(":18080")
}

// This example create several reverse proxy to show how to use jd proxy middleware.
// The request will be proxied to destination only when request is from specified datacenter.
func ExampleNewZoneProxy() {
	proxies := map[string]string{
		"/index":        "http://api.bilibili.com/html/index",
		"/ping":         "http://api.bilibili.com/api/ping",
		"/api/versions": "http://api.bilibili.com/api/web/versions",
	}

	engine := blademaster.Default()
	// proxy to specified destination
	for path, ep := range proxies {
		engine.GET(path, proxy.NewZoneProxy("sh004", ep), func(ctx *blademaster.Context) {
			ctx.String(200, "Origin")
		})
	}
	// proxy with request path
	ug := engine.Group("/update", proxy.NewZoneProxy("sh004", "http://sh001-api.bilibili.com"))
	ug.POST("/name", func(ctx *blademaster.Context) {
		ctx.String(500, "Should not be accessed")
	})
	ug.POST("/sign", func(ctx *blademaster.Context) {
		ctx.String(500, "Should not be accessed")
	})
	engine.Run(":18080")
}
