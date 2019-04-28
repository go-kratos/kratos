package http

import (
	"net/http"

	"go-common/app/service/live/xroom/internal/service"
	"go-common/library/conf/paladin"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	svc *service.Service
)

// New new a bm server.
func New(s *service.Service) (engine *bm.Engine) {
	var (
		hc struct {
			Server *bm.ServerConfig
		}
	)
	if err := paladin.Get("http.toml").UnmarshalTOML(&hc); err != nil {
		if err != paladin.ErrNotExist {
			panic(err)
		}
	}
	svc = s
	engine = bm.DefaultServer(hc.Server)
	initRouter(engine, verify.New(nil))
	if err := engine.Start(); err != nil {
		panic(err)
	}
	return
}

func initRouter(e *bm.Engine, v *verify.Verify) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/x/xroom")
	{
		g.GET("/start", v.Verify, howToStart)
	}
}

func ping(ctx *bm.Context) {
	if err := svc.Ping(ctx); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

// example for http request handler.
func howToStart(c *bm.Context) {
	c.String(0, "Golang 大法好 !!!")
}
