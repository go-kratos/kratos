package http

import (
	"net/http"

	"go-common/app/service/main/secure/conf"
	"go-common/app/service/main/secure/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	rSrv   *service.Service
	verSvc *verify.Verify
)

// Init init http service
func Init(s *service.Service) {
	initService(s)
	e := bm.DefaultServer(conf.Conf.BM)
	innerRouter(e)
	if err := e.Start(); err != nil {
		log.Error("e.Start() error(%v)", err)
		panic(err)
	}
}

func initService(s *service.Service) {
	rSrv = s
	verSvc = verify.New(conf.Conf.Verify)
}

// innerRouter init local router api path.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	group := e.Group("/x/internal/secure", verSvc.Verify)
	{
		group.GET("/loc", getLoc)
		group.POST("/often/check", oftenCheck)
		expect := group.Group("/expect")
		{
			expect.GET("/status", status)
			expect.POST("/close", closeNotify)
			expect.GET("/loc", loc)
			expect.POST("/feedback", feedback)
			expect.POST("/test", addlog)
		}
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := rSrv.Ping(c); err != nil {
		log.Error("ping error(%+v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// register check server ok.
func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
