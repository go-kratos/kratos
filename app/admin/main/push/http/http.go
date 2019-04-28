package http

import (
	"net/http"

	"go-common/app/admin/main/push/conf"
	"go-common/app/admin/main/push/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	pushSrv *service.Service
	authSrv *permit.Permit
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	pushSrv = s
	authSrv = permit.New(c.Auth)
	engine := bm.DefaultServer(c.HTTPServer)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/x/admin/push")
	{
		gapp := g.Group("/apps")
		{
			gapp.GET("/list", authSrv.Permit("PUSH_APPS_LIST"), appList)
			gapp.GET("/info", authSrv.Permit("PUSH_APPS"), appInfo)
			gapp.POST("/add", authSrv.Permit("PUSH_APPS"), addApp)
			gapp.POST("/save", authSrv.Permit("PUSH_APPS"), saveApp)
			gapp.POST("/delete", authSrv.Permit("PUSH_APPS"), delApp)
		}
		gauth := g.Group("/auths", authSrv.Permit("PUSH_AUTH"))
		{
			gauth.GET("/list", authList)
			gauth.GET("/info", authInfo)
			gauth.POST("/add", addAuth)
			gauth.POST("/save", saveAuth)
			gauth.POST("/delete", delAuth)
		}
		gbiz := g.Group("/business")
		{
			gbiz.GET("/list", authSrv.Permit("PUSH_BUSINESS_LIST"), businessList)
			gbiz.GET("/info", authSrv.Permit("PUSH_BUSINESS"), businessInfo)
			gbiz.POST("/add", authSrv.Permit("PUSH_BUSINESS"), addBusiness)
			gbiz.POST("/save", authSrv.Permit("PUSH_BUSINESS"), saveBusiness)
			gbiz.POST("/delete", authSrv.Permit("PUSH_BUSINESS"), delBusiness)
		}
		gtask := g.Group("/tasks")
		{
			gtask.GET("/list", authSrv.Permit("PUSH_TASK"), taskList)
			gtask.GET("/info", authSrv.Permit("PUSH_TASK"), taskInfo)
			gtask.POST("/add", authSrv.Permit("PUSH_TASK"), addTask)
			gtask.POST("/save", authSrv.Permit("PUSH_TASK"), saveTask)
			gtask.POST("/delete", authSrv.Permit("PUSH_TASK"), delTask)
			gtask.POST("/upload", authSrv.Permit("PUSH_TASK"), upload)
			gtask.POST("/upimg", authSrv.Permit("PUSH_TASK"), upimg)
			gtask.POST("/stop", authSrv.Permit("PUSH_TASK"), stopTask)
			gtask.POST("/confirm", authSrv.Permit("PUSH_CONFIRM"), confirmTask)
			gdp := gtask.Group("dataplatform", authSrv.Permit("PUSH_TASK"))
			{
				gdp.POST("/add", addDPTask)
				gdp.GET("/info", dpTaskInfo)
				gdp.GET("/check", checkDpData)
			}
		}
		gtest := g.Group("/test", authSrv.Permit("PUSH_TEST"))
		{
			gtest.POST("/mid", testPushMid)     // 按 mids 测试推送
			gtest.POST("/token", testPushToken) // 按单个 token 测试推送
		}
	}
}

func ping(ctx *bm.Context) {
	if err := pushSrv.Ping(ctx); err != nil {
		log.Error("push-admin ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
