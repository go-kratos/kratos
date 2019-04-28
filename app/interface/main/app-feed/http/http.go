package http

import (
	"go-common/app/interface/main/app-feed/conf"
	"go-common/app/interface/main/app-feed/service/external"
	"go-common/app/interface/main/app-feed/service/feed"
	pingsvc "go-common/app/interface/main/app-feed/service/ping"
	"go-common/app/interface/main/app-feed/service/region"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	// depend service
	authSvc   *auth.Auth
	verifySvc *verify.Verify
	// self service
	regionSvc   *region.Service
	feedSvc     *feed.Service
	pingSvc     *pingsvc.Service
	externalSvc *external.Service
)

// Init is
func Init(c *conf.Config) {
	initService(c)
	// init external router
	engineOut := bm.DefaultServer(c.BM.Outer)
	outerRouter(engineOut)
	// init outer server
	if err := engineOut.Start(); err != nil {
		log.Error("engineOut.Start() error(%v)", err)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	authSvc = auth.New(nil)
	verifySvc = verify.New(nil)
	// init self service
	regionSvc = region.New(c)
	feedSvc = feed.New(c)
	pingSvc = pingsvc.New(c)
	externalSvc = external.New(c)
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	// formal api
	feed := e.Group("/x/feed")
	feed.GET("/region/tags", authSvc.GuestMobile, tags)
	feed.GET("/subscribe/tags", authSvc.UserMobile, subTags)
	feed.POST("/subscribe/tags/add", authSvc.UserMobile, addTag)
	feed.POST("/subscribe/tags/cancel", authSvc.UserMobile, cancelTag)
	feed.GET("/index", authSvc.GuestMobile, feedIndex)
	feed.GET("/index/tab", authSvc.GuestMobile, feedIndexTab)
	feed.GET("/upper", authSvc.UserMobile, feedUpper)
	feed.GET("/upper/archive", authSvc.UserMobile, feedUpperArchive)
	feed.GET("/upper/bangumi", authSvc.UserMobile, feedUpperBangumi)
	feed.GET("/upper/recent", authSvc.UserMobile, feedUpperRecent)
	feed.GET("/upper/article", authSvc.UserMobile, feedUpperArticle)
	feed.GET("/upper/unread/count", authSvc.UserMobile, feedUnreadCount)
	feed.GET("/dislike", authSvc.GuestMobile, feedDislike)
	feed.GET("/dislike/cancel", authSvc.GuestMobile, feedDislikeCancel)
	feed.POST("/rcmd/up", verifySvc.Verify, upRcmd)

	feedV2 := e.Group("/x/v2/feed")
	feedV2.GET("/index", authSvc.Guest, feedIndex2)
	feedV2.GET("/index/tab", authSvc.Guest, feedIndexTab2)
	feedV2.GET("/index/converge", authSvc.Guest, feedIndexConverge)

	// live dynamic
	external := e.Group("/x/feed/external")
	external.GET("/dynamic/count", dynamicCount)
	external.GET("/dynamic/new", dynamicNew)
	external.GET("/dynamic/history", dynamicHistory)
}
