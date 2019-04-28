package http

import (
	"net/http"

	"go-common/app/service/main/spy/conf"
	"go-common/app/service/main/spy/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	spySvc  *service.Service
	verifyN *verify.Verify
)

// Init init http server.
func Init(c *conf.Config, s *service.Service) {
	spySvc = s
	verifyN = verify.New(nil)
	// init inner router
	engineInner := bm.DefaultServer(c.BM)
	innerRouter(engineInner)
	// init local server
	if err := engineInner.Start(); err != nil {
		log.Error("xhttp.Serve error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	// health check
	e.Ping(ping)
	e.Register(register)
	//new defined api lists
	group := e.Group("/x/internal/v1/spy", verifyN.Verify)
	{
		group.GET("/info", info)
		group.GET("/purge", purgeUser)
		group.POST("/purge2", purgeUser2)
		group.GET("/stat", stat)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if spySvc.Ping(c) != nil {
		log.Error("spy-service service ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// register check server ok.
func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
