package http

import (
	"kratos/tool/bmproto/examples/api"
	"kratos/tool/bmproto/examples/internal/service"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/time"

	xtime "time"
)

var svc *service.GreeterService

// New new a bm server.
func New(s *service.GreeterService) (engine *bm.Engine) {
	svc = s
	engine = bm.DefaultServer(&bm.ServerConfig{
		Addr:    "0.0.0.0:8000",
		Timeout: time.Duration(1 * xtime.Second),
	})
	initRouter(engine, s)
	if err := engine.Start(); err != nil {
		panic(err)
	}
	return
}

func initRouter(e *bm.Engine, s *service.GreeterService) {
	e.Ping(ping)
	e.Register(register)
	api.RegisterGreeterBMServer(e, s)
}

func ping(c *bm.Context) {
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
