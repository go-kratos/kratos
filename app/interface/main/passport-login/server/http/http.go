package http

import (
	"net/http"
	"strconv"

	"go-common/app/interface/main/passport-login/conf"
	"go-common/app/interface/main/passport-login/service"
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
	srv = service.New(c)
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
	g := e.Group("/x/passport-login")
	{
		g.GET("/key", getKey)
	}
	proxy := e.Group("/x/internal/passport-login", vfy.Verify)
	{
		// common
		proxy.GET("/key", getKey)
		// user
		proxy.GET("/user/check", proxyCheckUserData)
		// cookie
		proxy.POST("/cookie/add", proxyAddCookie)
		proxy.POST("/cookie/delete", proxyDeleteCookie)
		proxy.POST("/cookies/delete", proxyDeleteCookies)
		// token
		proxy.POST("/token/add", proxyAddToken)
		proxy.POST("/token/delete", proxyDeleteToken)
		proxy.POST("/tokens/delete", proxyDeleteTokens)
		proxy.POST("/tokens/game/delete", proxyDeleteGameTokens)
		proxy.POST("/token/renew", proxyRenewToken)
		// login
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

func parseInt(midStr string) (mid int64) {
	if len(midStr) == 0 {
		return 0
	}
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		return 0
	}
	return
}
