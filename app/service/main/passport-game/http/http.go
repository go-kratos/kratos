package http

import (
	"net/http"

	"go-common/app/service/main/passport-game/conf"
	"go-common/app/service/main/passport-game/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	srv *service.Service
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	srv = s
	// init inner router
	engine := bm.DefaultServer(c.BM)
	initRouter(engine)
	// init inner server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// initRouter init inner router.
func initRouter(e *bm.Engine) {
	// health check
	e.Ping(ping)
	e.Register(register)
	// new defined api lists
	g := e.Group("/x/passport-game")
	{
		g.GET("/oauth", verify, oauth)
		g.GET("/myinfo", verify, myInfo)
		g.GET("/info", verify, info)
		g.GET("/key", verify, getKeyProxy)
		g.GET("/login", verify, loginProxy)
		g.GET("/renewtoken", verify, renewToken)
		g.GET("/regions", verify, regions)
		g.POST("/reg/v3", verify, regV3)
		g.POST("/reg/v2", verify, regV2)
		g.POST("/reg", verify, reg)
		g.POST("/reg/byTel", verify, byTel)

		g.GET("/captcha", verify, captcha)
		g.POST("/sendSms", verify, sendSms)
	}

	inner := e.Group("/cache/pb")
	{
		inner.GET("/token", tokenPBCache)
		inner.GET("/info", infoPBCache)
	}

}

func verify(c *bm.Context) {
	app, err := verifySign(c, srv)
	if err != nil {
		c.JSON(nil, err)
		c.Abort()
	}
	c.Set("app", app)
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
