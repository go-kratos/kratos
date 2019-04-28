package http

import (
	"go-common/app/job/main/click/conf"
	"go-common/app/job/main/click/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var srv *service.Service

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

// ping check server ok.
func ping(c *bm.Context) {}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	// path
	g := e.Group("/x/internal/click")
	{
		g.GET("", click)
		g.GET("/lock", lock)
		g.GET("/lock/mid", lockMid)
	}
}
