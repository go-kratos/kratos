package http

import (
	"go-common/app/admin/openplatform/sug/conf"
	"go-common/app/admin/openplatform/sug/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	srv     *service.Service
	authSrv *permit.Permit
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	srv = s
	authSrv = permit.New(c.Auth)
	engine := bm.DefaultServer(c.BM)
	router(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// outerRouter init outer router api path.
func router(e *bm.Engine) {
	//init api
	e.Ping(ping)
	group := e.Group("/x/admin/sug")
	{
		seasonGroup := group.Group("/season")
		{
			seasonGroup.GET("/source/search", sourceSearch)
			seasonGroup.GET("/match/search", search)
			seasonGroup.POST("/match/operate", matchOperate)
		}
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		c.Error = err
		c.AbortWithStatus(503)
	}
}
