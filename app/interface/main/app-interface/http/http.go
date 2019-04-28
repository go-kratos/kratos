package http

import (
	"go-common/app/interface/main/app-interface/conf"
	acc "go-common/app/interface/main/app-interface/service/account"
	"go-common/app/interface/main/app-interface/service/dataflow"
	"go-common/app/interface/main/app-interface/service/display"
	"go-common/app/interface/main/app-interface/service/favorite"
	"go-common/app/interface/main/app-interface/service/history"
	"go-common/app/interface/main/app-interface/service/relation"
	"go-common/app/interface/main/app-interface/service/search"
	"go-common/app/interface/main/app-interface/service/space"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/proxy"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/queue/databus"
)

var (
	verifySvc   *verify.Verify
	authSvc     *auth.Auth
	spaceSvr    *space.Service
	srcSvr      *search.Service
	displaySvr  *display.Service
	favSvr      *favorite.Service
	accSvr      *acc.Service
	relSvr      *relation.Service
	historySvr  *history.Service
	dataflowSvr *dataflow.Service
	// databus
	userActPub *databus.Databus
	config     *conf.Config
)

// Init init http
func Init(c *conf.Config) {
	initService(c)
	// init external router
	engineOut := bm.DefaultServer(c.BM.Outer)
	// init outer router
	outerRouter(engineOut)
	if err := engineOut.Start(); err != nil {
		log.Error("engineOut.Start() error(%v) | config(%v)", err, c)
		panic(err)
	}
}

func initService(c *conf.Config) {
	verifySvc = verify.New(nil)
	authSvc = auth.New(nil)
	spaceSvr = space.New(c)
	srcSvr = search.New(c)
	displaySvr = display.New(c)
	favSvr = favorite.New(c)
	accSvr = acc.New(c)
	relSvr = relation.New(c)
	historySvr = history.New(c)
	dataflowSvr = dataflow.New(c)
	userActPub = databus.New(c.UseractPub)
	config = c
}

func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	proxyHandler := proxy.NewZoneProxy("sh004", "http://sh001-app.bilibili.com")
	account := e.Group("/x/v2/account", verifySvc.Verify)
	account.GET("/myinfo", myinfo)
	account.GET("/mine", authSvc.GuestMobile, mine)
	account.GET("/mine/ipad", authSvc.GuestMobile, mineIpad)

	space := e.Group("/x/v2/space")
	space.GET("", authSvc.GuestMobile, spaceAll)
	space.GET("/archive", authSvc.GuestMobile, upArchive)
	space.GET("/article", authSvc.GuestMobile, upArticle)
	space.GET("/bangumi", authSvc.GuestMobile, bangumi)
	space.GET("/coinarc", authSvc.GuestMobile, coinArc)
	space.GET("/likearc", authSvc.GuestMobile, likeArc)
	space.GET("/community", authSvc.GuestMobile, community)
	space.GET("/contribute", proxyHandler, authSvc.GuestMobile, contribute)
	space.GET("/contribute/cursor", proxyHandler, authSvc.GuestMobile, contribution)
	space.GET("/clips", authSvc.GuestMobile, clips)
	space.GET("/albums", authSvc.GuestMobile, albums)
	space.POST("/report", verifySvc.Verify, report)
	space.POST("/upContribute", proxyHandler, verifySvc.Verify, upContribute)

	search := e.Group("/x/v2/search")
	search.GET("", authSvc.GuestMobile, searchAll)
	search.GET("/type", authSvc.GuestMobile, searchByType)
	search.GET("/episodes", authSvc.GuestMobile, searchEpisodes)
	search.GET("/live", authSvc.GuestMobile, searchLive)
	search.GET("/hot", authSvc.GuestMobile, hotSearch)
	search.GET("/suggest", authSvc.GuestMobile, suggest)
	search.GET("/suggest2", authSvc.GuestMobile, suggest2)
	search.GET("/suggest3", authSvc.GuestMobile, suggest3)
	search.GET("/defaultwords", authSvc.GuestMobile, defaultWords)
	search.GET("/user", authSvc.GuestMobile, searchUser)
	search.GET("/recommend", authSvc.GuestMobile, recommend)
	search.GET("/recommend/noresult", authSvc.GuestMobile, recommendNoResult)
	search.GET("/recommend/pre", authSvc.GuestMobile, recommendPre)
	search.GET("/resource", authSvc.GuestMobile, resource)

	display := e.Group("/x/v2/display", verifySvc.Verify)
	display.GET("/zone", zone)
	display.GET("/id", authSvc.GuestMobile, displayID)

	favorite := e.Group("/x/v2/favorite", verifySvc.Verify)
	favorite.GET("", authSvc.GuestMobile, folder)
	favorite.GET("/video", authSvc.GuestMobile, favoriteVideo)
	favorite.GET("/topic", authSvc.GuestMobile, topic)
	favorite.GET("/article", authSvc.GuestMobile, article)
	favorite.GET("/clips", authSvc.GuestMobile, favClips)
	favorite.GET("/albums", authSvc.GuestMobile, favAlbums)
	favorite.GET("/sp", specil)
	favorite.GET("/audio", authSvc.GuestMobile, audio)
	favorite.GET("/tab", authSvc.UserMobile, tab)

	relation := e.Group("/x/v2/relation")
	relation.GET("/followings", authSvc.GuestMobile, followings)
	relation.GET("/tag", authSvc.UserMobile, tag)

	history := e.Group("/x/v2/history", verifySvc.Verify)
	history.GET("", authSvc.UserMobile, historyList)
	history.GET("/live", live)
	history.GET("/liveList", authSvc.UserMobile, liveList)
	history.GET("/cursor", authSvc.UserMobile, historyCursor)
	history.POST("/del", authSvc.UserMobile, historyDel)
	history.POST("/clear", authSvc.UserMobile, historyClear)

	dataflow := e.Group("/x/v2/dataflow")
	dataflow.POST("/report", reportInfoc)
}
