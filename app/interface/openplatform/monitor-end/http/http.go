package http

import (
	"go-common/app/interface/openplatform/monitor-end/conf"
	"go-common/app/interface/openplatform/monitor-end/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"net/http"
)

var (
	mfSvc   *service.Service
	authSvr *auth.Auth
)

// Init init account service.
func Init(c *conf.Config, s *service.Service) {
	authSvr = auth.New(nil)
	mfSvc = s
	// engine
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	// init inner server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start() error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	group := e.Group("/open/monitor", authSvr.Guest)
	{
		group.GET("/report", report)
		group.GET("/err/collect", collect)
	}
	group = e.Group("open/internal/monitor", authSvr.Guest)
	{
		group.GET("/report", report)
		group.GET("/consume/start", startConsume)
		group.GET("/consume/stop", stopConsume)
		group.GET("/consume/pause", pauseConsume)
		r := group.Group("/alert")
		{
			cr := r.Group("/group")
			{
				cr.GET("/list", groupList)
				cr.POST("/update", groupUpdate)
				cr.POST("/add", groupAdd)
				cr.POST("/delete", groupDelete)
			}
			cr = r.Group("/target")
			{
				cr.GET("/list", targetList)
				cr.POST("/update", targetUpdate)
				cr.POST("/add", targetAdd)
				cr.POST("/sync", targetSync)
			}
			cr = r.Group("/product")
			{
				cr.POST("/add", productAdd)
				cr.POST("/update", productUpdate)
				cr.POST("/delete", productDelete)
				cr.GET("/list", productList)
			}
		}
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := mfSvc.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// register support discovery.
func register(c *bm.Context) {
	c.JSON(map[string]struct{}{}, nil)
}
