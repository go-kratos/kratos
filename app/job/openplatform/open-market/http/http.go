package http

import (
	"net/http"

	"go-common/app/job/openplatform/open-market/conf"
	"go-common/app/job/openplatform/open-market/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var srv *service.Service

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	srv = s
	// init external router
	engineIn := bm.DefaultServer(c.BM.Inner)
	innerRouter(engineIn)
	// init Inner server
	if err := engineIn.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
}

// ping check server ok.
func ping(c *bm.Context) {
	err := srv.Ping(c)
	if err != nil {
		log.Error("app-job service ping error(%+v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
