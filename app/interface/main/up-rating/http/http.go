package http

import (
	"net/http"

	"go-common/app/interface/main/up-rating/conf"
	"go-common/app/interface/main/up-rating/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
)

var (
	svc     *service.Service
	authSvr *auth.Auth
)

// Init http server
func Init(c *conf.Config) {
	initService(c)
	engine := bm.DefaultServer(c.BM)
	externalRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func initService(c *conf.Config) {
	svc = service.New(c)
	authSvr = auth.New(nil)
}

func externalRouter(e *bm.Engine) {
	e.Ping(ping)
	// define routers
	group := e.Group("/studio/up-rating", authSvr.User)
	{
		group.GET("/info", upRating)
	}
	cache := group.Group("/cache")
	{
		cache.GET("/expire/up", expireUpRating)
	}
}

func ping(c *bm.Context) {
	var err error
	if err = svc.Ping(c); err != nil {
		log.Error("service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
