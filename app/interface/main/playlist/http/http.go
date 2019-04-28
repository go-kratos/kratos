package http

import (
	"net/http"

	"go-common/app/interface/main/playlist/conf"
	"go-common/app/interface/main/playlist/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	plSvc   *service.Service
	authSvr *auth.Auth
	vfySvr  *verify.Verify
)

// Init init http server
func Init(c *conf.Config, s *service.Service) {
	authSvr = auth.New(c.Auth)
	vfySvr = verify.New(c.Verify)
	plSvc = s
	engine := bm.DefaultServer(c.HTTPServer)
	outerRouter(engine)
	internalRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	// init api
	e.Ping(ping)
	group := e.Group("/x/playlist", bm.CORS())
	{
		group.GET("/whitelist", authSvr.Guest, whiteList)
		group.GET("/report", report)
		group.GET("/share/report", reportShare)
		group.GET("", authSvr.Guest, list)
		group.GET("/info", authSvr.Guest, info)
		group.POST("/add", authSvr.User, add)
		group.POST("/del", authSvr.User, del)
		group.POST("/update", authSvr.User, update)
		videoGroup := group.Group("/video")
		{
			videoGroup.GET("", videoList)
			videoGroup.POST("/check", authSvr.User, check)
			videoGroup.GET("/toview", authSvr.Guest, toView)
			videoGroup.POST("/add", authSvr.User, addVideo)
			videoGroup.POST("/del", authSvr.User, delVideo)
			videoGroup.POST("/sort", authSvr.User, sortVideo)
			videoGroup.POST("/desc/edit", authSvr.User, editVideoDesc)
			videoGroup.GET("/search", authSvr.User, searchVideo)
		}
		favGroup := group.Group("/fav")
		{
			favGroup.GET("", authSvr.Guest, listFavorite)
			favGroup.POST("/add", authSvr.User, addFavorite)
			favGroup.POST("/del", authSvr.User, delFavorite)
		}
	}
}
func internalRouter(e *bm.Engine) {
	group := e.Group("/x/internal/playlist")
	{
		group.GET("/whitelist", vfySvr.Verify, whiteList)
		group.GET("/report", vfySvr.Verify, report)
		group.GET("/share/report", vfySvr.Verify, reportShare)
		group.GET("", vfySvr.Verify, list)
		group.GET("/info", vfySvr.Verify, info)
		videoGroup := group.Group("/video")
		{
			videoGroup.GET("", vfySvr.Verify, videoList)
			videoGroup.GET("/toview", vfySvr.Verify, toView)
			videoGroup.GET("/search", vfySvr.Verify, searchVideo)
		}
		favGroup := group.Group("/fav")
		{
			favGroup.GET("", vfySvr.Verify, listFavorite)
		}
	}
}
func ping(c *bm.Context) {
	if err := plSvc.Ping(c); err != nil {
		log.Error("playlist interface ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
