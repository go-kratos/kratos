package http

import (
	"go-common/app/interface/main/shorturl/conf"
	"go-common/app/interface/main/shorturl/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	// depend service
	idfSvc *verify.Verify
	suSvr  *service.Service
)

// initService .
func initService(c *conf.Config) {
	suSvr = service.New(c)
	idfSvc = verify.New(c.Verify)
}

// Init init http
func Init(c *conf.Config) {
	initService(c)
	// init internal router
	engineInner := bm.NewServer(c.BM)
	engineInner.Use(bm.Recovery(), bm.Trace(), bm.CSRF(), bm.Mobile(), logger())
	innerRouter(engineInner)
	if err := engineInner.Start(); err != nil {
		log.Error("engineInner.Start error(%v)", err)
		panic(err)
	}
}

// innerRouter .
func innerRouter(e *bm.Engine) {
	e.GET("/monitor/ping", ping)
	e.GET("/", jump)
	b := e.Group("/x/internal/shorturl")
	{
		b.POST("/add", idfSvc.Verify, add)
		b.POST("/update", idfSvc.Verify, shortUpdate)
		b.GET("/detail", idfSvc.Verify, shortByID)
		b.GET("/list", idfSvc.Verify, shortAll)
		b.POST("/del", idfSvc.Verify, shortDel)
	}
}
