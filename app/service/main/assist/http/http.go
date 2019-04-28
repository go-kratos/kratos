package http

import (
	"go-common/app/service/main/assist/conf"
	"go-common/app/service/main/assist/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	assSvc    *service.Service
	verifySvc *verify.Verify
)

// Init init account service.
func Init(c *conf.Config, s *service.Service) {
	assSvc = service.New(c)
	verifySvc = verify.New(nil)
	engineInner := bm.DefaultServer(c.BM)
	innerRouter(engineInner)
	if err := engineInner.Start(); err != nil {
		log.Error("engineInner.Start() error(%v) | config(%v)", err, c)
		panic(err)
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	e.GET("/monitor/ping", ping)
	cr := e.Group("/x/internal/assist", verifySvc.Verify)
	{
		cr.GET("/assists", assists)
		cr.GET("/ups", assistUps)
		cr.GET("/stat", assistsMids)
		cr.GET("/info", assistInfo)
		cr.GET("/ids", assistIDs)
		cr.POST("/add", assistAdd)
		cr.POST("/del", assistDel)
		cr.POST("/exit", assistExit)
		cr.GET("/logs", assistLogs)
		cr.GET("/log/info", assistLogInfo)
		cr.POST("/log/add", assistLogAdd)
		cr.POST("/log/cancel", assistLogCancel)
		cr.GET("/log/obj", assistLogObj)
	}
}
