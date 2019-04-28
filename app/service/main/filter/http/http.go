package http

import (
	"net/http"

	"go-common/app/service/main/filter/conf"
	"go-common/app/service/main/filter/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	verifySvc *verify.Verify
	svc       *service.Service
)

// Init init http service
func Init(s *service.Service) {
	svc = s
	verifySvc = verify.New(nil)

	engine := bm.DefaultServer(conf.Conf.BM)
	internalRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("%+v", err)
		panic(err)
	}
}

func internalRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	group := e.Group("/x/internal/filter", verifySvc.Verify)
	{
		group.GET("", filter)
		group.GET("/multi", mfilter)
		group.POST("/post", filter)
		group.POST("/mpost", mfilter)
		group.POST("/area/mpost", areaMfilter)
		group.POST("/article", article)
		group.POST("/test", filterTest)
		group.GET("/key/test", testKey)

		group1 := group.Group("/v2")
		{
			group1.POST("/hit", hit)
		}
		groupV3 := group.Group("/v3")
		{
			groupV3.POST("/hit", hitV3)
		}
	}
}

// register check server ok.
func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("ping error(%+v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
