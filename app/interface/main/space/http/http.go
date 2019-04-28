package http

import (
	"net/http"

	"go-common/app/interface/main/space/conf"
	"go-common/app/interface/main/space/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/supervisor"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	authSvr *auth.Auth
	spcSvc  *service.Service
	spvSvc  *supervisor.Supervisor
	vfySvc  *verify.Verify
)

// Init init http server
func Init(c *conf.Config, s *service.Service) {
	authSvr = auth.New(c.Auth)
	spvSvc = supervisor.New(c.Supervisor)
	vfySvc = verify.New(c.Verify)
	spcSvc = s
	// init http server
	engine := bm.DefaultServer(c.HTTPServer)
	outerRouter(engine)
	internalRouter(engine)
	// init Outer serve
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func outerRouter(e *bm.Engine) {
	e.GET("/monitor/ping", ping)
	group := e.Group("/x/space", bm.CORS())
	{
		chGroup := group.Group("/channel")
		{
			chGroup.GET("", authSvr.Guest, channel)
			chGroup.GET("/index", authSvr.Guest, channelIndex)
			chGroup.GET("/list", authSvr.Guest, channelList)
			chGroup.POST("/add", spvSvc.ServeHTTP, authSvr.User, addChannel)
			chGroup.POST("/del", authSvr.User, delChannel)
			chGroup.POST("/edit", spvSvc.ServeHTTP, authSvr.User, editChannel)
		}
		chvGroup := group.Group("/channel/video")
		{
			chvGroup.GET("", authSvr.Guest, channelVideo)
			chvGroup.POST("/add", authSvr.User, addChannelVideo)
			chvGroup.POST("/del", authSvr.User, delChannelVideo)
			chvGroup.POST("/sort", authSvr.User, sortChannelVideo)
			chvGroup.GET("/check", authSvr.User, checkChannelVideo)
		}
		riderGroup := group.Group("/rider")
		{
			riderGroup.GET("/list", authSvr.User, riderList)
			riderGroup.POST("/exit", authSvr.User, exitRider)
		}
		tagGroup := group.Group("/tag")
		{
			tagGroup.POST("/sub", authSvr.User, tagSub)
			tagGroup.POST("/sub/cancel", authSvr.User, tagCancelSub)
			tagGroup.GET("/sub/list", authSvr.Guest, tagSubList)
		}
		bgmGroup := group.Group("/bangumi")
		{
			bgmGroup.POST("/concern", authSvr.User, bangumiConcern)
			bgmGroup.POST("/unconcern", authSvr.User, bangumiUnConcern)
			bgmGroup.GET("/concern/list", authSvr.Guest, bangumiList)
		}
		topGroup := group.Group("/top")
		{
			topGroup.GET("/arc", authSvr.Guest, topArc)
			topGroup.POST("/arc/set", authSvr.User, setTopArc)
			topGroup.POST("/arc/cancel", authSvr.User, cancelTopArc)
			topGroup.POST("/dynamic/set", authSvr.User, setTopDynamic)
			topGroup.POST("/dynamic/cancel", authSvr.User, cancelTopDynamic)
		}
		mpGroup := group.Group("/masterpiece")
		{
			mpGroup.GET("", authSvr.Guest, masterpiece)
			mpGroup.POST("/add", authSvr.User, addMasterpiece)
			mpGroup.POST("/edit", authSvr.User, editMasterpiece)
			mpGroup.POST("/cancel", authSvr.User, cancelMasterpiece)
		}
		noticeGroup := group.Group("/notice")
		{
			noticeGroup.GET("", notice)
			noticeGroup.POST("/set", authSvr.User, setNotice)
		}
		accGroup := group.Group("/acc")
		{
			accGroup.GET("/info", authSvr.Guest, accInfo)
			accGroup.GET("/tags", accTags)
			accGroup.POST("/tags/set", authSvr.User, setAccTags)
			accGroup.GET("/relation", authSvr.User, relation)
		}
		themeGroup := group.Group("theme")
		{
			themeGroup.GET("/list", authSvr.User, themeList)
			themeGroup.POST("/active", authSvr.User, themeActive)
		}
		appGroup := group.Group("/app")
		{
			appGroup.GET("/index", authSvr.Guest, appIndex)
			appGroup.GET("/dynamic/list", authSvr.Guest, dynamicList)
			appGroup.GET("/played/game", authSvr.Guest, appPlayedGame)
			appGroup.GET("/top/photo", authSvr.Guest, appTopPhoto)
		}
		arcGroup := group.Group("/arc")
		{
			arcGroup.GET("/search", arcSearch)
			arcGroup.GET("/list", arcList)
		}
		group.GET("/setting", settingInfo)
		group.GET("/article", article)
		group.GET("/navnum", authSvr.Guest, navNum)
		group.GET("/upstat", upStat)
		group.GET("/shop", authSvr.User, shopInfo)
		group.GET("/album/index", albumIndex)
		group.GET("/fav/nav", authSvr.Guest, favNav)
		group.GET("/fav/arc", authSvr.Guest, favArc)
		group.GET("/coin/video", authSvr.Guest, coinVideo)
		group.GET("/myinfo", authSvr.User, myInfo)
		group.GET("/lastplaygame", authSvr.Guest, lastPlayGame)
		group.POST("/privacy/modify", authSvr.User, privacyModify)
		group.POST("/index/order/modify", authSvr.User, indexOrderModify)
	}
}

func internalRouter(e *bm.Engine) {
	group := e.Group("/x/internal/space")
	{
		group.GET("/setting", vfySvc.Verify, settingInfo)
		group.GET("/myinfo", vfySvc.Verify, authSvr.User, myInfo)
		group.POST("/privacy/modify", authSvr.User, privacyModify)
		group.POST("/privacy/batch/modify", authSvr.User, privacyBatchModify)
		group.POST("/index/order/modify", authSvr.User, indexOrderModify)
		accGroup := group.Group("/acc")
		{
			accGroup.GET("/info", vfySvc.Verify, authSvr.Guest, accInfo)
		}
		appGroup := group.Group("/app")
		{
			appGroup.GET("/index", vfySvc.Verify, authSvr.Guest, appIndex)
		}
		group.GET("/web/index", vfySvc.Verify, authSvr.Guest, webIndex)
		group.POST("/cache/clear", clearCache)
		group.GET("/blacklist", vfySvc.Verify, blacklist)
	}
}

func ping(c *bm.Context) {
	if err := spcSvc.Ping(c); err != nil {
		log.Error("space-interface ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
