package http

import (
	"net/http"

	"go-common/app/admin/ep/tapd/conf"
	"go-common/app/admin/ep/tapd/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	svc     *service.Service
	authSvc *permit.Permit
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	svc = s
	authSvc = permit.New(c.Auth)
	engine := bm.DefaultServer(c.BM)
	engine.Ping(ping)
	outerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {

	inner := e.Group("/internal/ep/tapd")
	{
		inner.GET("/version", getVersion)
		inner.POST("/test", test)
		inner.POST("/testform", testform)

		version1 := inner.Group("/v1", authSvc.Permit(""))
		{
			version1.POST("/hook/update", updateHook)
			version1.POST("/hook/query", queryHook)

			version1.GET("/hook/event/query", queryURLEvent)
			version1.GET("/hook/cache/save", saveHookUrlInCache)
			version1.GET("/hook/cache/query", queryHookUrlInCache)

			version1.POST("/eventlog/query", queryEventLog)
		}
	}

	outer := e.Group("/ep/tapd")
	{
		outer.POST("/callback", tapdCallback)
		outer.GET("/version", getVersion)
	}
}

func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func getVersion(c *bm.Context) {
	v := new(struct {
		Version string `json:"version"`
	})
	v.Version = "v.0.0.7"
	c.JSON(v, nil)

}
