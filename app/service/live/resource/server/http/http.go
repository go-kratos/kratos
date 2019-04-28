package http

import (
	"go-common/app/service/live/resource/api/http/v1"
	"go-common/app/service/live/resource/conf"
	"go-common/app/service/live/resource/service"
	v12 "go-common/app/service/live/resource/service/v1"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"net/http"
)

var (
	srv           *service.Service
	vfy           *verify.Verify
	titansService *v12.TitansService
)

// Init init
func Init(c *conf.Config, srv *service.Service) {
	vfy = verify.New(c.Verify)
	initService(c)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func initService(c *conf.Config) {
	srv = service.New(c)
	titansService = v12.NewTitansService(c)
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	e.GET("/xlive/internal/resource/v1/titans/getMyTreeApps", getNodes)
	g := e.Group("/xlive/resource")
	{
		g.GET("/start", vfy.Verify, howToStart)
	}
	v1.RegisterV1TitansService(e, titansService, map[string]bm.HandlerFunc{})
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

// example for http request handler
func howToStart(c *bm.Context) {
	c.String(0, "Golang 大法好 !!!")
}
