package http

import (
	"net/http"

	"go-common/app/interface/main/web-goblin/conf"
	"go-common/app/interface/main/web-goblin/service/share"
	"go-common/app/interface/main/web-goblin/service/web"
	"go-common/app/interface/main/web-goblin/service/wechat"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
)

var (
	srvWeb    *web.Service
	srvShare  *share.Service
	srvWechat *wechat.Service
	authSvr   *auth.Auth
)

// Init init .
func Init(c *conf.Config) {
	authSvr = auth.New(c.Auth)
	srvWeb = web.New(c)
	srvShare = share.New(c)
	srvWechat = wechat.New(c)
	engine := bm.DefaultServer(c.BM)
	router(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func router(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	group := e.Group("/x/web-goblin")
	{
		miGroup := group.Group("/mi")
		{
			miGroup.GET("/full", fullshort)
		}
		channelGroup := group.Group("/channel")
		{
			channelGroup.GET("", authSvr.Guest, channel)
		}
		ugcGroup := group.Group("ugc")
		{
			ugcGroup.GET("/full", ugcfull)
			ugcGroup.GET("/increment", ugcincre)
		}
		pgcGroup := group.Group("pgc")
		{
			pgcGroup.GET("/full", pgcfull)
			pgcGroup.GET("/increment", pgcincre)
		}
		group.GET("/share/encourage", authSvr.User, encourage)
		group.GET("/recruit", recruit)
		weChatGroup := group.Group("/wechat")
		{
			weChatGroup.GET("/qrcode", qrcode)
		}
	}
}

func ping(c *bm.Context) {
	if err := srvWeb.Ping(c); err != nil {
		log.Error("web-goblin ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
