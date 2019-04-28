package http

import (
	"net/http"

	"go-common/app/service/main/card/conf"
	"go-common/app/service/main/card/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	srv       *service.Service
	verifySvc *verify.Verify
)

// Init init
func Init(c *conf.Config, svc *service.Service) {
	srv = svc
	verifySvc = verify.New(nil)
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
	group := e.Group("/x/internal/card", verifySvc.Verify)
	{
		group.GET("/bymids", byMids)
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
