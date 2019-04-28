package http

import (
	"context"

	"go-common/app/job/main/search/conf"
	"go-common/app/job/main/search/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	ctx = context.TODO()
	svr *service.Service
)

// Init init http.
func Init(c *conf.Config, s *service.Service) {
	svr = s
	// init local router
	engine := bm.DefaultServer(c.HTTPServer)
	route(engine)
	// init local server
	if err := engine.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

// innerRouter init local router api path.
func route(e *bm.Engine) {
	e.GET("/x/search-job/action", action)
	e.GET("/x/search-job/stat", stat)
	e.Ping(ping)
}

// ping check server ok.
func ping(ctx *bm.Context) {
	if err := svr.Ping(ctx); err != nil {
		log.Error("search job ping error(%v)", err)
		ctx.AbortWithStatus(503)
	}
}
