package http

import (
	"net/url"
	"strconv"

	"go-common/app/interface/main/tv/conf"
	appsrv "go-common/app/interface/main/tv/service/app"
	auditsrv "go-common/app/interface/main/tv/service/audit"
	"go-common/app/interface/main/tv/service/favorite"
	gobsrv "go-common/app/interface/main/tv/service/goblin"
	hissrv "go-common/app/interface/main/tv/service/history"
	"go-common/app/interface/main/tv/service/pgc"
	secsrv "go-common/app/interface/main/tv/service/search"
	"go-common/app/interface/main/tv/service/thirdp"
	"go-common/app/interface/main/tv/service/tvvip"
	viewsrv "go-common/app/interface/main/tv/service/view"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	favSvc    *favorite.Service
	tvSvc     *appsrv.Service
	viewSvc   *viewsrv.Service
	auditSvc  *auditsrv.Service
	gobSvc    *gobsrv.Service
	secSvc    *secsrv.Service
	thirdpSvc *thirdp.Service
	tvVipSvc  *tvvip.Service
	authSvc   *auth.Auth
	vfySvc    *verify.Verify
	hisSvc    *hissrv.Service
	pgcSvc    *pgc.Service
	signCfg   *conf.AuditSign
)

// Init init http sever instance.
func Init(c *conf.Config) {
	signCfg = c.Cfg.AuditSign
	initService(c)
	// init outer router
	engineOut := bm.NewServer(c.HTTPServer)
	engineOut.Use(bm.Recovery(), bm.Trace(), bm.Logger(), bm.Mobile())
	outerRouter(engineOut)
	if err := engineOut.Start(); err != nil {
		log.Error("engineOut.Start error(%v)", err)
		panic(err)
	}
}

func parseInt(value string) int64 {
	intval, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		intval = 0
	}
	return intval
}

func takeBuild(req url.Values) {
	buildStr := req.Get("build")
	if buildStr != "" {
		if tvSvc.TVAppInfo.Build != buildStr {
			tvSvc.TVAppInfo.Build = buildStr
		}
	}
	platStr := req.Get("platform")
	if platStr != "" {
		if tvSvc.TVAppInfo.Platform != platStr {
			tvSvc.TVAppInfo.Platform = buildStr
		}
	}
	mobiStr := req.Get("mobi_app")
	if mobiStr != "" {
		if tvSvc.TVAppInfo.MobiApp != mobiStr {
			tvSvc.TVAppInfo.MobiApp = mobiStr
		}
	}
}

func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	tv := e.Group("/x/tv", bm.CORS(), bm.CSRF())
	e.GET("/x/tv/vip/order/guest_create", authSvc.User, createGuestOrder)
	{
		app := tv.Group("", authSvc.Guest) // the public group
		{
			// app pages
			app.GET("/homepage", homepage)
			app.GET("/zonepage", zonePage)
			app.GET("/zone_index", zoneIdx)
			app.GET("/media_detail", mediaDetail)
			app.GET("/modpage", modpage)
			// app functions
			app.GET("/upgrade", upgrade)
			app.GET("/splash", splash)
			app.GET("/recommend", recommend)
			app.GET("/suggest", searchSug)
			app.GET("/hotword", hotword)
			app.GET("/history", history)
			// dangbei page
			app.GET("/dangbei", dbeiPage)
			// video audit status check
			app.GET("/loadep", loadEP)
			app.GET("/labels", labels)
		}
		aud := e.Group("/x/tv/audit", bm.CSRF()) // license owner audit related functions
		{
			aud.POST("", audit)
			aud.POST("/transcode", vfySvc.Verify, transcode)
			aud.POST("/apply/pgc", vfySvc.Verify, applyPGC)
		}
		pgc := e.Group("/x/tv/pgc", bm.CSRF(), authSvc.Guest)
		{
			pgc.GET("/view", mDetailV2)
		}
		ugc := e.Group("/x/tv/ugc", bm.CSRF(), authSvc.Guest) // the APIs dedicated for ugc
		{
			ugc.GET("/view", view)
			ugc.GET("/load_video", loadVideo)
			ugc.GET("/playurl", ugcPlayurl)
		}
		search := e.Group("/x/tv/search", bm.CSRF(), authSvc.Guest) // the APIs for search
		{
			search.GET("/types", searchTypes)
			search.GET("", searchResult)
			wild := search.Group("/wild")
			{
				wild.GET("", searchAll)      // 综合搜索
				wild.GET("user", userSearch) // 按用户搜索
				wild.GET("pgc", pgcSearch)   // pgc番剧影视
			}
		}
		fav := e.Group("/x/tv/favorites", bm.CSRF(), authSvc.Guest)
		{
			fav.GET("", favorites)
			fav.POST("/act", favAct)
		}
		mango := e.Group("/x/tv/mango", bm.CSRF())
		{
			mango.GET("/recom", mangoRecom)
		}
		third := e.Group("/x/tv/third", bm.CSRF())
		{
			third.GET("/pgc/season", mangoSnPage)
			third.GET("/pgc/ep", mangoEpPage)
			third.GET("/ugc/archive", mangoArcPage)
			third.GET("/ugc/video", mangoVideoPage)
		}
		idx := e.Group("/x/tv/index", bm.CSRF(), authSvc.Guest)
		{
			idx.GET("/pgc", pgcIdx)
			idx.GET("/ugc", ugcIdx)
		}
		tv.GET("/region", region) // all region info
		vip := e.Group("/x/tv/vip", bm.CSRF())
		{
			vip.GET("/user/info", authSvc.UserMobile, vipInfo)
			vip.GET("/user/yst_info", ystVipInfo)
			vip.GET("/panel/user", authSvc.UserMobile, panelInfo)
			vip.GET("/panel/guest", authSvc.Guest, guestPanelInfo)

			vip.POST("/order/qr", authSvc.UserMobile, createQr)
			vip.POST("/order/guest_qr", authSvc.Guest, createGuestQr)
			vip.GET("/order/create", authSvc.Guest, createOrder)

			vip.GET("/token/info", authSvc.UserMobile, tokenStatus)

			vip.POST("/callback/pay", payCallback)
			vip.POST("/callback/wx_contract", wxContractCallback)

		}
	}
}

// ping check db server ok.
func ping(c *bm.Context) {}

func initService(c *conf.Config) {
	tvSvc = appsrv.New(c)
	viewSvc = viewsrv.New(c)
	favSvc = favorite.New(c)
	auditSvc = auditsrv.New(c)
	gobSvc = gobsrv.New(c)
	secSvc = secsrv.New(c)
	authSvc = auth.New(c.Auth)
	vfySvc = verify.New(c.Verify)
	hisSvc = hissrv.New(c)
	thirdpSvc = thirdp.New(c)
	pgcSvc = pgc.New(c)
	tvVipSvc = tvvip.New(c)
}
