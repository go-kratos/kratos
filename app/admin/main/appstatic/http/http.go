package http

import (
	"go-common/app/admin/main/appstatic/conf"
	"go-common/app/admin/main/appstatic/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	vfySvc  *verify.Verify
	authSvc *permit.Permit
	apsSvc  *service.Service
)

// Init http server
func Init(c *conf.Config, s *service.Service) {
	initService(c, s)
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	// init internal server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// initService init service
func initService(c *conf.Config, s *service.Service) {
	apsSvc = s
	authSvc = permit.New(c.Auth)
	vfySvc = verify.New(nil)
}

// innerRouter
func innerRouter(e *bm.Engine) {
	// ping monitor
	e.GET("/monitor/ping", ping)
	// internal api
	bg := e.Group("/x/admin/appstatic/res")
	{
		bg.POST("/add_ver", authSvc.Permit("APP_RESOURCE_POOL_MGT"), addVer)           // 从mgr上传，正式权限
		bg.POST("/add_ver_test", authSvc.Permit("APP_RESOURCE_POOL_MGT_EDIT"), addVer) // 从mgr上传，测试权限
		bg.POST("/upload", vfySvc.Verify, addVer)                                      // 从其他系统上传
		bg.POST("/publish", vfySvc.Verify, publish)                                    // 告知某资源包的第一次发布，用于触发增量包补充计算
	}
}

// ping check server ok.
func ping(c *bm.Context) {}
