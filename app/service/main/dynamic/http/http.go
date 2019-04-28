package http

import (
	"net/http"

	"go-common/app/service/main/dynamic/conf"
	"go-common/app/service/main/dynamic/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	dySvc  *service.Service
	vfySvr *verify.Verify
)

// Init init.
func Init(c *conf.Config, s *service.Service) {
	vfySvr = verify.New(c.Verify)
	dySvc = s
	engineInner := bm.DefaultServer(c.HTTPServer)
	innerRouter(engineInner)
	// init inner serve
	if err := engineInner.Start(); err != nil {
		log.Error("engineInner.Start error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	group := e.Group("/x/dynamic", bm.CORS())
	{
		group.GET("/tag", vfySvr.Verify, regionTagArcs)
		group.GET("/region", vfySvr.Verify, regionArcs)
		group.GET("/regions", vfySvr.Verify, regionsArcs)
		group.GET("/region/total", vfySvr.Verify, regionTotal)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := dySvc.Ping(c); err != nil {
		log.Error("dynamic service ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(nil, nil)
}
