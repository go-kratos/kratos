package http

import (
	"net/http"

	"go-common/app/service/live/xrewardcenter/conf"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	//midAuth *auth.Auth
	vfy *verify.Verify
)

// Init init
func Init(c *conf.Config) {
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
	//midMap := map[string]bm.HandlerFunc{"auth": midAuth.User, "verify": vfy.Verify}

	g := e.Group("/x/xrewardcenter")
	{
		g.GET("/start", vfy.Verify, howToStart)
	}
}

func ping(c *bm.Context) {
	c.AbortWithStatus(http.StatusOK)
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

// example for http request handler
func howToStart(c *bm.Context) {
	c.String(0, "Golang 大法好 !!!")
}
