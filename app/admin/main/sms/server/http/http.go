package http

import (
	"net/http"

	"go-common/app/admin/main/sms/conf"
	"go-common/app/admin/main/sms/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	authSrv *permit.Permit
	svc     *service.Service
)

// Init http server
func Init(c *conf.Config, s *service.Service) {
	svc = s
	authSrv = permit.New(c.Auth)
	engine := bm.DefaultServer(c.HTTPServer)
	route(engine)
	if err := engine.Start(); err != nil {
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/x/admin/sms")
	{
		tplg := g.Group("/template", authSrv.Permit("SMS_TEMPLATE"))
		{
			tplg.POST("/add", addTemplate)
			tplg.POST("/update", updateTemplate)
			tplg.GET("/list", templateList)
		}
		taskg := g.Group("/task", authSrv.Permit("SMS_TASK"))
		{
			taskg.POST("/add", addTask)
			taskg.POST("/update", updateTask)
			taskg.POST("/delete", deleteTask)
			taskg.GET("/info", taskInfo)
			taskg.GET("/list", taskList)
			taskg.POST("/upload", upload)
		}
	}
}

func ping(ctx *bm.Context) {
	if err := svc.Ping(ctx); err != nil {
		log.Error("sms-admin ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
