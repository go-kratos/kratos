package http

import (
	"net/http"

	"go-common/app/service/bbq/search/api/grpc/v1"
	"go-common/app/service/bbq/search/conf"
	"go-common/app/service/bbq/search/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	srv *service.Service
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	srv = s
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
	g := e.Group("/bbq/internal/search")
	{
		g.GET("/start", howToStart)
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

// example for http request handler
func howToStart(c *bm.Context) {
	arg := new(v1.RecVideoDataRequest)

	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(srv.RecVideoData(c, arg))
}
