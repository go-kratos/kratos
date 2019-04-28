package http

import (
	"net/http"

	"go-common/app/interface/main/kvo/conf"
	"go-common/app/interface/main/kvo/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
)

var (
	kvoSvr  *service.Service
	authSvr *auth.Auth
)

// Init init http
func Init(c *conf.Config) {
	kvoSvr = service.New(c)
	authSvr = auth.New(c.Auth)
	// init outer router
	engineOut := bm.DefaultServer(c.BM)
	outerRouter(engineOut)
	if err := engineOut.Start(); err != nil {
		log.Error("engineOut.Start error(%v)", err)
		panic(err)
	}
}

func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	group := e.Group("/x/kvo", bm.CORS())
	{
		group.GET("/web/doc/get", authSvr.UserWeb, doc)
		group.POST("/web/doc/add", authSvr.UserWeb, addDoc)
		group.GET("/app/doc/get", authSvr.UserMobile, doc)
		group.POST("/app/doc/add", authSvr.UserMobile, addDoc)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	var err error
	if err = kvoSvr.Ping(c); err != nil {
		log.Error("kvo service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
