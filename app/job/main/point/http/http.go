package http

import (
	"go-common/app/job/main/point/conf"
	"go-common/app/job/main/point/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	svc *service.Service
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	// service init
	svc = s
	// init router
	engineInner := bm.DefaultServer(c.BM)
	innerRouter(engineInner)
	if err := engineInner.Start(); err != nil {
		log.Error("engineInner.Start error(%v)", err)
		panic(err)
	}
}

// outerRouter init outer router api path.
func innerRouter(e *bm.Engine) {
	//init api
	e.Ping(ping)
}

// ping check server ok.
func ping(c *bm.Context) {
	svc.Ping(c)
}
