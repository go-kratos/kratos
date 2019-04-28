package http

import (
	"net/http"

	"go-common/app/job/main/stat/conf"
	"go-common/app/job/main/stat/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var srv *service.Service

// Init init http service
func Init(c *conf.Config, s *service.Service) {
	srv = s
	e := bm.DefaultServer(c.BM)
	innerRouter(e)
	// init internal server
	if err := e.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

// innerRouter init local router api path.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
