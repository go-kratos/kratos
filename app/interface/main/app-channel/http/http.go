package http

import (
	"go-common/app/interface/main/app-channel/conf"
	channelSvr "go-common/app/interface/main/app-channel/service/channel"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/proxy"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	// depend service
	channelSvc *channelSvr.Service
	verifySvc  *verify.Verify
	authSvc    *auth.Auth
)

func Init(c *conf.Config) {
	initService(c)
	// init external router
	engineOut := bm.DefaultServer(c.BM.Outer)
	outerRouter(engineOut)
	// init Outer server
	if err := engineOut.Start(); err != nil {
		log.Error("engineOut.Start() error(%v) | config(%v)", err, c)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	channelSvc = channelSvr.New(c)
	verifySvc = verify.New(nil)
	authSvc = auth.New(nil)
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	proxyHandler := proxy.NewZoneProxy("sh004", "http://sh001-app.bilibili.com")
	cl := e.Group("/x/channel", verifySvc.Verify)
	{
		feed := cl.Group("/feed", authSvc.GuestMobile)
		{
			feed.GET("", index)
			feed.GET("/index", proxyHandler, index2)
			feed.GET("/tab", tab)
			feed.GET("/tab/list", tablist)
		}
		cl.POST("/add", authSvc.UserMobile, subscribeAdd)
		cl.POST("/cancel", authSvc.UserMobile, subscribeCancel)
		cl.POST("/update", authSvc.UserMobile, subscribeUpdate)
		cl.GET("/list", authSvc.GuestMobile, list)
		cl.GET("/subscribe", authSvc.UserMobile, subscribe)
		cl.GET("/discover", authSvc.GuestMobile, discover)
		cl.GET("/category", authSvc.GuestMobile, category)
		cl.GET("/square", authSvc.GuestMobile, square)
		cl.GET("/mysub", authSvc.UserMobile, mysub)
	}
}
