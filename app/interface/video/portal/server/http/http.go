package http

import (
	"net/http"

	"go-common/app/interface/video/portal/conf"
	"go-common/app/interface/video/portal/middleware"
	"go-common/app/interface/video/portal/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	srv *service.Service
	vfy *verify.Verify
)

// Init init
func Init(c *conf.Config, svc *service.Service) {
	srv = svc
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
	g := e.Group("/video/portal")
	{
		g.GET("/", howToStart)
		g.GET("/whitelist", middleware.TXCheckSecretkey, txStreamForwardingConf)
		g.GET("/whitelist/conf", middleware.BvcCheckSecret, bvcStreamForwardingConf)
		g.GET("/leveltwosourcestation", streamLevelTwoSourceStation)
		g.GET("/lplallstar", streamLplAllStar)
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
	//c.String(0, "Golang 大法好 !!!")
	c.JSON(nil, nil)
}
