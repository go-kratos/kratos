package http

import (
	"go-common/app/infra/canal/conf"
	"go-common/app/infra/canal/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	cs *service.Canal
)

// Init int http service
func Init(c *conf.Config, cs *service.Canal) {
	initService(cs)
	// init router
	eg := bm.DefaultServer(c.BM)
	initRouter(eg)
	// init Outer serve
	if err := eg.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

func initService(canal *service.Canal) {
	cs = canal
}

// initRouter init outer router api path.
func initRouter(e *bm.Engine) {
	// init api
	e.Ping(ping)
	group := e.Group("/x/internal/canal")
	{
		group.GET("/infoc/post", infocPost)
		group.GET("/infoc/current", infocCurrent)
		group.GET("/errors", errors)
		group.POST("/master/check", checkMaster)
		group.POST("/test/sync", syncPos)
	}
}

func ping(c *bm.Context) {
}
