package http

import (
	"net/http"

	"go-common/app/service/live/recommend/api/grpc/v1"
	"go-common/app/service/live/recommend/internal/conf"
	"go-common/app/service/live/recommend/internal/service"
	v12 "go-common/app/service/live/recommend/internal/service/v1"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	vfy *verify.Verify
	svc *service.Service
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	svc = s
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
	g := e.Group("/x/recommend")
	{
		g.GET("/start", vfy.Verify, howToStart)
	}
	v1.RegisterV1RecommendService(e, v12.NewRecommendService(conf.Conf), nil)
}

func ping(ctx *bm.Context) {
	if err := svc.Ping(ctx); err != nil {
		log.Error("ping error(%+v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

// example for http request handler
func howToStart(c *bm.Context) {
	c.String(0, "Golang 大法好 !!!")
}
