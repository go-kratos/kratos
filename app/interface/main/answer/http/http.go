package http

import (
	"go-common/app/interface/main/answer/conf"
	"go-common/app/interface/main/answer/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
)

var (
	authSvc *auth.Auth
	antSvc  *antispam.Antispam
	svc     *service.Service
)

// Init init
func Init(c *conf.Config) {
	initService(c)
	engineOuter := bm.DefaultServer(c.BM)
	outerRouter(engineOuter)
	if err := engineOuter.Start(); err != nil {
		log.Error("engineOuter.Start() error(%v)", err)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	authSvc = auth.New(c.AuthN)
	antSvc = antispam.New(c.Antispam)
	svc = service.New(c)
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	og := e.Group("/x/answer", bm.CORS())
	{
		// 答题流程排序
		og.GET("/v3/base", authSvc.UserWeb, antSvc.Handler(), baseQus)
		og.POST("/v3/base/check", authSvc.User, checkBase)
		og.GET("/v3/extra", authSvc.User, antSvc.Handler(), extraQus)
		og.POST("/v3/extra/check", authSvc.UserWeb, checkExtra)
		og.GET("/v3/extra/score", authSvc.User, extraScore)
		og.GET("/user/birthday", authSvc.User, checkBirthDay)
		og.GET("/v3/pro/type", authSvc.User, proType)
		og.GET("/v3/pro", authSvc.UserWeb, antSvc.Handler(), proQus)
		og.POST("/v3/pro/check", authSvc.User, checkPro)
		og.GET("/v3/captcha/gt", authSvc.User, antSvc.Handler(), captcha)
		og.POST("/v3/captcha/check", authSvc.User, validate)
		og.GET("/v3/result", authSvc.GuestWeb, cool)
		og.POST("/rec/pendant", authSvc.User, pendantRec)
	}
}

func ping(c *bm.Context) {
}
