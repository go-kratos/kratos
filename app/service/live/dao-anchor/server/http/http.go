package http

import (
	"go-common/app/service/live/dao-anchor/conf"
	"go-common/app/service/live/dao-anchor/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	srv *service.Service
	vfy *verify.Verify
)

// Init init
func Init(c *conf.Config) {
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
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/xlive/dao-anchor")
	{
		g.GET("/start", vfy.Verify, howToStart)
	}
}

func ping(c *bm.Context) {
	// TODO
	/*if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}*/
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

// example for http request handler
func howToStart(c *bm.Context) {
	c.String(0, "Golang 大法好 !!!")
}
