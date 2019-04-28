package http

import (
	"net/http"

	"go-common/app/job/main/favorite/conf"
	"go-common/app/job/main/favorite/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var favSvc *service.Service

// Init init http
func Init(c *conf.Config, s *service.Service) {
	favSvc = s
	// init external router
	engineOut := bm.DefaultServer(c.BM)
	outerRouter(engineOut)
	// init Outer serve
	if err := engineOut.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	e.Ping(ping)
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := favSvc.Ping(c); err != nil {
		log.Error("favorite http service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
