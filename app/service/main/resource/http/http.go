package http

import (
	"go-common/app/service/main/resource/conf"
	"go-common/app/service/main/resource/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	vfySvc *verify.Verify
	resSvc *service.Service
)

// Init int http service
func Init(c *conf.Config, s *service.Service) {
	vfySvc = verify.New(c.Verify)
	resSvc = s
	// init internal router
	engineInner := bm.DefaultServer(c.BM.Inner)
	innerRouter(engineInner)
	// init internal server
	if err := engineInner.Start(); err != nil {
		log.Error("engineInner.Start() error(%v) | config(%v)", err, c)
		panic(err)
	}
	// init external router
	engineLocal := bm.DefaultServer(c.BM.Local)
	localRouter(engineLocal)
	// init external server
	if err := engineLocal.Start(); err != nil {
		log.Error("engineLocal.Start() error(%v) | config(%v)", err, c)
		panic(err)
	}
}

// innerRouter init outer router api path.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	rs := e.Group("/x/internal/resource")

	bn := rs.Group("/banner")
	bn.GET("", banner)

	ads := rs.Group("/ads")
	ads.GET("/paster/app", vfySvc.Verify, pasterAPP)
	ads.GET("/paster/pgc", vfySvc.Verify, pasterPGC)

	res := rs.Group("/res")
	res.GET("", vfySvc.Verify, resource)
	res.GET("/resources", vfySvc.Verify, resources)
	res.GET("/indexIcon", vfySvc.Verify, indexIcon)
	res.GET("/playerIcon", vfySvc.Verify, playerIcon)
	res.GET("/cmtbox", vfySvc.Verify, cmtbox)
	res.GET("/regionCard", vfySvc.Verify, regionCard)
	res.GET("/audit", vfySvc.Verify, audit)
}

// localRouter init local router api path.
func localRouter(e *bm.Engine) {
	e.GET("/x/resource/version", version)
	e.GET("/x/resource/monitor", monitor)
}
