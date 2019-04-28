package http

import (
	"go-common/app/interface/main/app-tag/conf"
	pingSvr "go-common/app/interface/main/app-tag/service/ping"
	"go-common/app/interface/main/app-tag/service/tag"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
)

var (
	// depend service
	tagSvr  *tag.Service
	pingSvc *pingSvr.Service
	authSvc *auth.Auth
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
	tagSvr = tag.New(c)
	pingSvc = pingSvr.New(c)
	authSvc = auth.New(nil)
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	ft := e.Group("/x/feed/tag", authSvc.GuestMobile)
	{
		ft.GET("/detail", tagDetail)
		ft.GET("/change/default", tagDefault)
		ft.GET("/change/new", tagNew)
	}
	t := e.Group("/x/v2/tag", authSvc.GuestMobile)
	{
		t.GET("/tab", tagTab)
		t.GET("/dynamic", tagDynamic)
		t.GET("/dynamic/index", tagDynamicIndex)
		t.GET("/dynamic/list", tagDynamicList)
		t.GET("/dynamic/rank/list", tagRankList)
	}
}
