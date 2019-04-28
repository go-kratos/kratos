package http

import (
	"go-common/app/service/main/upcredit/conf"
	"go-common/app/service/main/upcredit/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	idfSvc *verify.Verify
	//Svc service.
	Svc     *service.Service
	authSrc *permit.Permit
)

// Init init account service.
func Init(c *conf.Config) {
	// service
	initService(c)
	// init internal router
	innerEngine := bm.DefaultServer(c.BM.Inner)
	setupInnerEngine(innerEngine)
	// init internal server
	if err := innerEngine.Start(); err != nil {
		log.Error("httpx.Serve2 error(%v)", err)
		panic(err)
	}
}

func initService(c *conf.Config) {
	idfSvc = verify.New(nil)
	Svc = service.New(c)
	authSrc = permit.New(c.Auth)
}

// innerRouter
func setupInnerEngine(e *bm.Engine) {
	// monitor ping
	e.GET("/monitor/ping", ping)
	// base
	var base *bm.RouterGroup
	if conf.Conf.IsTest {
		base = e.Group("/x/internal/upcredit")

	} else {
		base = e.Group("/x/internal/upcredit", idfSvc.Verify)
	}
	{
		base.GET("/test", test)

		base.POST("/log/add", logCredit)
		base.GET("/log/get", logGet)

		base.GET("/score/get", scoreGet)
		base.POST("/score/recalc", recalc)
		base.POST("/score/calc_section", calcSection)

	}

}

// ping check server ok.
func ping(ctx *bm.Context) {
}
