package http

import (
	"go-common/app/service/main/passport-auth/conf"
	"go-common/app/service/main/passport-auth/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	svr *service.Service
	vfy *verify.Verify
)

// Init init config
func Init(c *conf.Config, srv *service.Service) {
	initService(c, srv)
	engineOut := bm.DefaultServer(c.BM)
	outerRouter(engineOut)
	if err := engineOut.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
	}
}

func initService(c *conf.Config, srv *service.Service) {
	svr = srv
	vfy = verify.New(c.VerifyConfig)

}

// outerRouter init outer router
func outerRouter(e *bm.Engine) {
	e.GET("/monitor/ping", ping)
	group := e.Group("/x/internal/passport-auth", vfy.Verify)
	{
		group.GET("/cookie_info", cookieInfo)
		group.GET("/token_info", tokenInfo)
		group.GET("/refresh_info", refreshInfo)
		group.GET("/old_token_info", oldTokenInfo)
	}
}

func ping(c *bm.Context) {
}
