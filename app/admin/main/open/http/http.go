package http

import (
	"net/http"

	"go-common/app/admin/main/open/conf"
	"go-common/app/admin/main/open/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	mngSvc *service.Service
	idfSvc *verify.Verify
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	mngSvc = s
	idfSvc = verify.New(c.Verify)
	// init inner router
	engineIn := bm.DefaultServer(nil)
	innerRouter(engineIn)
	// init inner server
	if err := engineIn.Start(); err != nil {
		log.Error("enginIn.Start.error", err)
		panic(err)
	}
}

// innerRouter.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/x/admin/open")
	{
		gapp := g.Group("/app", idfSvc.Verify)
		{
			gapp.POST("/add", addApp)
			gapp.POST("/delete", delApp)
			gapp.POST("/update", updateApp)
			gapp.GET("/list", listApp)
		}
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := mngSvc.Ping(c); err != nil {
		log.Error("service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
