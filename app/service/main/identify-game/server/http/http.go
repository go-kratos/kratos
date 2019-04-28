package http

import (
	"net/http"

	"go-common/app/service/main/identify-game/conf"
	"go-common/app/service/main/identify-game/service"
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
	// engine
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	// init inner server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	group := e.Group("/x/internal/identify-game", vfy.Verify)
	{
		group.GET("/oauth", oauth)
		group.GET("/renewtoken", renewToken)
		group.GET("/regions", regions)
		group.GET("/cookies", getCookieByToken)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// register support discovery.
func register(c *bm.Context) {
	c.JSON(map[string]struct{}{}, nil)
}
