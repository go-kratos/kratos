package http

import (
	"net/http"

	"go-common/app/interface/main/captcha/conf"
	"go-common/app/interface/main/captcha/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/rate"
	verifyx "go-common/library/net/http/blademaster/middleware/verify"
)

var (
	svr       *service.Service
	verifySvc *verifyx.Verify
)

// Init captcha http init.
func Init(c *conf.Config, s *service.Service) (err error) {
	svr = s
	verifySvc = verifyx.New(c.Verify)
	rateLimit := rate.New(c.Rate)
	engineOuter := bm.DefaultServer(c.BM.Outer)
	engineOuter.Use(rateLimit)
	outerRouter(engineOuter)
	interRouter(engineOuter)
	if err := engineOuter.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
	return
}

func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	group := e.Group("/x/v1/captcha")
	group.GET("/get", get)

}

func interRouter(e *bm.Engine) {
	group := e.Group("/x/internal/v1/captcha")
	group.GET("/token", verifySvc.Verify, token)
	group.POST("/verify", verifySvc.Verify, verify)
}

func ping(c *bm.Context) {
	if svr.Ping(c) != nil {
		log.Error("captcha service ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
