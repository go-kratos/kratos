package http

import (
	"net/http"

	"go-common/app/interface/main/esports/conf"
	"go-common/app/interface/main/esports/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	authn  *auth.Auth
	vfySvr *verify.Verify
	eSvc   *service.Service
)

// Init init http server
func Init(c *conf.Config, s *service.Service) {
	authn = auth.New(c.Auth)
	vfySvr = verify.New(c.Verify)
	eSvc = s
	engine := bm.DefaultServer(c.BM)
	outerRouter(engine)
	internalRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("httpx.Serve error(%v)", err)
		panic(err)
	}
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	group := e.Group("/x/esports", bm.CORS())
	{
		group.GET("/season", season)
		group.GET("/app/season", appSeason)
		matGroup := group.Group("/matchs")
		{
			matGroup.GET("/filter", filterMatch)
			matGroup.GET("/list", authn.Guest, listContest)
			matGroup.GET("/app/list", authn.Guest, appContest)
			matGroup.GET("/calendar", calendar)
			matGroup.GET("/active", authn.Guest, active)
			matGroup.GET("/videos", authn.Guest, actVideos)
			matGroup.GET("/points", authn.Guest, actPoints)
			matGroup.GET("/top", authn.Guest, actTop)
			matGroup.GET("/knockout", authn.Guest, actKnockout)
			matGroup.GET("/info", authn.Guest, contest)
			matGroup.GET("/recent", authn.Guest, recent)
		}
		videoGroup := group.Group("/video")
		{
			videoGroup.GET("/filter", filterVideo)
			videoGroup.GET("/list", listVideo)
			videoGroup.GET("/search", authn.Guest, search)
		}
		favGroup := group.Group("/fav")
		{
			favGroup.GET("", authn.Guest, listFav)
			favGroup.GET("/season", authn.Guest, seasonFav)
			favGroup.GET("/stime", authn.Guest, stimeFav)
			favGroup.GET("/list", authn.Guest, appListFav)
			favGroup.POST("/add", authn.User, addFav)
			favGroup.POST("/del", authn.User, delFav)
		}
		pointGroup := group.Group("/leida")
		{
			pointGroup.GET("/game", game)
			pointGroup.GET("/items", items)
			pointGroup.GET("/heroes", heroes)
			pointGroup.GET("/abilities", abilities)
			pointGroup.GET("/players", players)
		}
	}
}

func internalRouter(e *bm.Engine) {
	group := e.Group("/x/internal/esports")
	{
		group.GET("/season", vfySvr.Verify, season)
		group.GET("/app/season", vfySvr.Verify, appSeason)
		matGroup := group.Group("/matchs")
		{
			matGroup.GET("/filter", vfySvr.Verify, filterMatch)
			matGroup.GET("/list", vfySvr.Verify, listContest)
			matGroup.GET("/app/list", vfySvr.Verify, appContest)
			matGroup.GET("/calendar", vfySvr.Verify, calendar)
			matGroup.GET("/active", vfySvr.Verify, active)
			matGroup.GET("/videos", vfySvr.Verify, actVideos)
			matGroup.GET("/points", vfySvr.Verify, actPoints)
			matGroup.GET("/top", vfySvr.Verify, actTop)
			matGroup.GET("/knockout", vfySvr.Verify, actKnockout)
			matGroup.GET("/info", vfySvr.Verify, contest)
			matGroup.GET("/recent", vfySvr.Verify, recent)
		}
		videoGroup := group.Group("/video")
		{
			videoGroup.GET("/filter", vfySvr.Verify, filterVideo)
			videoGroup.GET("/list", vfySvr.Verify, listVideo)
			videoGroup.GET("/search", vfySvr.Verify, search)
		}
		favGroup := group.Group("/fav")
		{
			favGroup.GET("", vfySvr.Verify, listFav)
			favGroup.GET("/season", vfySvr.Verify, seasonFav)
			favGroup.GET("/stime", vfySvr.Verify, stimeFav)
			favGroup.GET("/list", vfySvr.Verify, appListFav)
		}
	}
}

func ping(c *bm.Context) {
	if err := eSvc.Ping(c); err != nil {
		log.Error("esports interface ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
