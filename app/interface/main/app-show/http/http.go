package http

import (
	"net/http"

	"go-common/app/interface/main/app-show/conf"
	"go-common/app/interface/main/app-show/service/banner"
	"go-common/app/interface/main/app-show/service/daily"
	pingSvr "go-common/app/interface/main/app-show/service/ping"
	"go-common/app/interface/main/app-show/service/rank"
	"go-common/app/interface/main/app-show/service/region"
	"go-common/app/interface/main/app-show/service/show"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/render"
)

var (
	// depend service
	authSvc *auth.Auth
	// self service
	bannerSvc *banner.Service
	regionSvc *region.Service
	showSvc   *show.Service
	pingSvc   *pingSvr.Service
	rankSvc   *rank.Service
	dailySvc  *daily.Service
)

func Init(c *conf.Config) {
	initService(c)
	// init external router
	engineOut := bm.DefaultServer(c.BM.Outer)
	outerRouter(engineOut)
	// init Outer server
	if err := engineOut.Start(); err != nil {
		log.Error("engineOut.Start() error(%v) | config(%v)", err, c)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	authSvc = auth.New(nil)
	bannerSvc = banner.New(c)
	regionSvc = region.New(c)
	showSvc = show.New(c)
	pingSvc = pingSvr.New(c)
	rankSvc = rank.New(c)
	dailySvc = daily.New(c)
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	bnnr := e.Group("/x/v2/banner", authSvc.GuestMobile)
	{
		bnnr.GET("", banners)
	}
	region := e.Group("/x/v2/region", authSvc.GuestMobile)
	{
		region.GET("", regions)
		region.GET("/list", regionsList)
		region.GET("/index", regionsIndex)
		region.GET("/show", regionShow)
		region.GET("/show/dynamic", regionShowDynamic)
		region.GET("/show/child", regionChildShow)
		region.GET("/show/child/list", regionChildListShow)
		region.GET("/dynamic", regionDynamic)
		region.GET("/dynamic/list", regionDynamicList)
		region.GET("/dynamic/child", regionDynamicChild)
		region.GET("/dynamic/child/list", regionDynamicChildList)
	}
	rank := e.Group("/x/v2/rank", authSvc.GuestMobile)
	{
		rank.GET("", rankAll)
		rank.GET("/region", rankRegion)
	}
	show := e.Group("/x/v2/show", authSvc.GuestMobile)
	{
		show.GET("", shows)
		show.GET("/region", showsRegion)
		show.GET("/index", showsIndex)
		show.GET("/widget", showWidget)
		show.GET("/temp", showTemps)
		show.GET("/change", showChange)
		show.GET("/change/live", showLiveChange)
		show.GET("/change/region", showRegionChange)
		show.GET("/change/bangumi", showBangumiChange)
		show.GET("/change/dislike", showDislike)
		show.GET("/change/article", showArticleChange)
		show.GET("/popular", popular)
		show.GET("/popular/index", popular2)
	}
	daily := e.Group("/x/v2/daily", authSvc.GuestMobile)
	{
		daily.GET("/list", dailyID)
	}
	column := e.Group("/x/v2/column", authSvc.GuestMobile)
	{
		column.GET("", columnList)
	}
	cg := e.Group("/x/v2/category", authSvc.GuestMobile)
	{
		cg.GET("", category)
	}
}

//returnJSON return json no message
func returnJSON(c *bm.Context, data interface{}, err error) {
	code := http.StatusOK
	c.Error = err
	bcode := ecode.Cause(err)
	c.Render(code, render.JSON{
		Code:    bcode.Code(),
		Message: "",
		Data:    data,
	})
}

//returnDataJSON return json no message
func returnDataJSON(c *bm.Context, data map[string]interface{}, ttl int, err error) {
	code := http.StatusOK
	if ttl < 1 {
		ttl = 1
	}
	if err != nil {
		c.JSON(nil, err)
	} else {
		if data != nil {
			data["code"] = 0
			data["message"] = ""
			data["ttl"] = ttl
		}
		c.Render(code, render.MapJSON(data))
	}
}
