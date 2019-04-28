package http

import (
	"net/http"

	"go-common/app/interface/main/web-feed/conf"
	"go-common/app/interface/main/web-feed/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/limit/aqm"
)

var (
	feedSrv *service.Service
	authSvr *auth.Auth
)

// Init init
func Init(c *conf.Config, srv *service.Service) {
	authSvr = auth.New(c.Auth)
	feedSrv = srv
	engine := bm.DefaultServer(c.BM)
	outerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	group := e.Group("/x/web-feed", authSvr.User)
	{
		group.GET("/feed", aqm.New(nil).Limit(), feed)
		group.GET("/feed/unread", aqm.New(nil).Limit(), feedUnread)
		art := group.Group("/article")
		{
			art.GET("/feed", aqm.New(nil).Limit(), articleFeed)
			art.GET("/unread", aqm.New(nil).Limit(), articleFeedUnread)
		}
	}
}

// Ping check server ok.
func ping(c *bm.Context) {
	if err := feedSrv.Ping(c); err != nil {
		log.Error("web-feed ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
