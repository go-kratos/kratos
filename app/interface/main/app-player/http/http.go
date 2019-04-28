package http

import (
	"go-common/app/interface/main/app-player/conf"
	"go-common/app/interface/main/app-player/service"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/stat/prom"
)

var (
	svr      *service.Service
	ver      *verify.Verify
	ah       *auth.Auth
	errCount = prom.BusinessErrCount
)

// Init init http
func Init(c *conf.Config) {
	initService(c)
	engine := bm.DefaultServer(nil)
	outerRouter(engine)
	if err := engine.Start(); err != nil {
		panic(err)
	}
}

func initService(c *conf.Config) {
	svr = service.New(c)
	ver = verify.New(nil)
	ah = auth.New(nil)
}

func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.GET("/x/playurl", ver.Verify, ah.GuestMobile, playurl)
}

// Ping is
func ping(ctx *bm.Context) {

}
