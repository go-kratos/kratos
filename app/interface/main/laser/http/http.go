package http

import (
	"go-common/app/interface/main/laser/conf"
	"go-common/app/interface/main/laser/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"net/http"
)

var (
	svc     *service.Service
	authSvr *auth.Auth
)

// Init http server
func Init(c *conf.Config) {
	// service
	initService(c)
	engine := bm.DefaultServer(c.BM)
	// init outer router
	outerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// service init
func initService(c *conf.Config) {
	svc = service.New(c)
	authSvr = auth.New(nil)
}

func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	app := e.Group("/x/laser/app", authSvr.UserMobile)
	{
		app.GET("/query", queryTask)
		app.POST("/update", updateTask)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	var err error
	if err = svc.Ping(c); err != nil {
		c.AbortWithStatus(http.StatusServiceUnavailable)
		log.Error("laser-interface ping error(%v)", err)
	}
}
