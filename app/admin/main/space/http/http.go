package http

import (
	"net/http"

	"go-common/app/admin/main/space/conf"
	"go-common/app/admin/main/space/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	spcSvc *service.Service
	//idfSvc  *identify.Identify
	permitSvc *permit.Permit
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	spcSvc = s
	permitSvc = permit.New(c.Permit)
	engine := bm.DefaultServer(c.BM)
	authRouter(engine)
	// init internal server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func authRouter(e *bm.Engine) {
	e.Ping(ping)
	group := e.Group("/x/admin/space")
	{
		noticeGroup := group.Group("/notice", permitSvc.Permit("SPACE_NOTICE"))
		{
			noticeGroup.GET("", notice)
			noticeGroup.POST("/up", noticeUp)
		}
		group.GET("/relation", relation)
		blacklist := group.Group("/blacklist", permitSvc.Permit("SPACE_BLACKLIST"))
		{
			blacklist.GET("", blacklistIndex)
			blacklist.POST("add", blacklistAdd)
			blacklist.POST("update", blacklistUp)
		}
	}
}

func ping(c *bm.Context) {
	if err := spcSvc.Ping(c); err != nil {
		log.Error("space-admin ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
