package http

import (
	"net/http"

	"go-common/app/job/main/usersuit/conf"
	"go-common/app/job/main/usersuit/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var svc *service.Service

// Init init
func Init(c *conf.Config) {
	initService(c)
	// init  inner router
	engineInner := bm.DefaultServer(c.BM)
	innerRouter(engineInner)
	if err := engineInner.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	svc = service.New(c)
}

// innerRouter init local router api path.
func innerRouter(e *bm.Engine) {
	//init api
	e.Ping(ping)
}

func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("usersuit-job ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
