package http

import (
	"go-common/app/job/main/archive/conf"
	"go-common/app/job/main/archive/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// Init init http router.
func Init(c *conf.Config, s *service.Service) {
	e := bm.DefaultServer(c.BM)
	innerRouter(e)
	// init internal server
	if err := e.Start(); err != nil {
		log.Error("xhttp.Serve error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
}

// ping check server ok.
func ping(c *bm.Context) {}
