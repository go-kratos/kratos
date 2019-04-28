package http

import (
	"go-common/app/admin/main/point/conf"
	"go-common/app/admin/main/point/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	svc     *service.Service
	authSvc *permit.Permit
)

// Init init
func Init(c *conf.Config) {
	initService(c)
	// init router
	engine := bm.DefaultServer(c.BM)
	initRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	authSvc = permit.New(c.Auth)
	svc = service.New(c)
}

// initRouter init outer router api path.
func initRouter(e *bm.Engine) {
	//init api
	e.Ping(ping)
	group := e.Group("/x/admin/point", authSvc.Permit("VIP_POINT"))
	{
		group.GET("/conf/list", pointConfList)
		group.GET("/conf/info", pointConfInfo)
		group.POST("/conf/add", pointConfAdd)
		group.POST("/conf/edit", pointConfEdit)
		group.GET("/history/list", pointHistory)
		group.POST("/user/add", pointUserAdd)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	svc.Ping(c)
}
