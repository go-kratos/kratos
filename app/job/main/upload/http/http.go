package http

import (
	"go-common/app/job/main/upload/conf"
	"go-common/app/job/main/upload/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	svc *service.Service
)

// Init init
func Init(c *conf.Config) {
	initService(c)
	// init router
	engineInner := bm.DefaultServer(c.BM.Inner)
	outerRouter(engineInner)
	if err := engineInner.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
	engineLocal := bm.DefaultServer(c.BM.Local)
	localRouter(engineLocal)
	if err := engineLocal.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	svc = service.New(c)
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	//init api
	e.GET("/monitor/ping", ping)
	group := e.Group("/x/job/main/upload")
	{
		group.GET("")
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	svc.Ping(c)
}

// innerRouter init local router api path.
func localRouter(e *bm.Engine) {
	group := e.Group("")
	{
		group.GET("")
	}
}
