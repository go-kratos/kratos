package http

import (
	"net/http"

	"go-common/app/service/bbq/recsys/conf"
	"go-common/app/service/bbq/recsys/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	srv *service.Service
	vfy *verify.Verify
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	srv = s
	vfy = verify.New(c.Verify)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/bbq/admin/recsys")
	{
		g.GET("/verify/start", vfy.Verify, start)
		g.GET("/start", start)
	}

	ga := e.Group("/bbq/admin/recsys")
	{
		ga.POST("/check/rec/message", reqRecsys)
		ga.POST("/check/related/message", relatedRecsys)
		ga.POST("/check/ups/message", upsRecsys)
	}
}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
