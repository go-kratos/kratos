package http

import (
	"net/http"

	"go-common/app/service/main/msm/conf"
	"go-common/app/service/main/msm/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	vfySvr *verify.Verify
	svr    *service.Service
)

// Init init config.
func Init(c *conf.Config, s *service.Service) {
	svr = s
	vfySvr = verify.New(nil)
	engine := bm.DefaultServer(c.BM)
	oldRouter(engine)
	innerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func oldRouter(e *bm.Engine) {
	e.Ping(ping)
	group := e.Group("/x/v1/msm", vfySvr.Verify)
	{
		group.GET("/codes/2", codes)
		group.POST("/conf/push", push)
		group.POST("/conf/setToken", setToken)
		group.GET("/codes/langs", codesLangs)
	}
}

func innerRouter(e *bm.Engine) {
	group := e.Group("/x/internal/msm/v1")
	{
		group.GET("/codes/2", vfySvr.Verify, codes)
		group.GET("/auth/scope", credential, scope)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if svr.Ping() != nil {
		log.Error("service ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
