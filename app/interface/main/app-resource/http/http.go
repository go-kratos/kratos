package http

import (
	"go-common/app/interface/main/app-resource/conf"
	absvr "go-common/app/interface/main/app-resource/service/abtest"
	auditsvr "go-common/app/interface/main/app-resource/service/audit"
	broadcastsvr "go-common/app/interface/main/app-resource/service/broadcast"
	domainsvr "go-common/app/interface/main/app-resource/service/domain"
	guidesvc "go-common/app/interface/main/app-resource/service/guide"
	modulesvr "go-common/app/interface/main/app-resource/service/module"
	"go-common/app/interface/main/app-resource/service/notice"
	"go-common/app/interface/main/app-resource/service/param"
	pingsvr "go-common/app/interface/main/app-resource/service/ping"
	pluginsvr "go-common/app/interface/main/app-resource/service/plugin"
	showsvr "go-common/app/interface/main/app-resource/service/show"
	sidesvr "go-common/app/interface/main/app-resource/service/sidebar"
	"go-common/app/interface/main/app-resource/service/splash"
	staticsvr "go-common/app/interface/main/app-resource/service/static"
	"go-common/app/interface/main/app-resource/service/version"
	whitesvr "go-common/app/interface/main/app-resource/service/white"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
)

var (
	// depend service
	authSvc *auth.Auth
	// self service
	pgSvr        *pluginsvr.Service
	pingSvr      *pingsvr.Service
	sideSvr      *sidesvr.Service
	verSvc       *version.Service
	paramSvc     *param.Service
	ntcSvc       *notice.Service
	splashSvc    *splash.Service
	auditSvc     *auditsvr.Service
	abSvc        *absvr.Service
	moduleSvc    *modulesvr.Service
	guideSvc     *guidesvc.Service
	staticSvc    *staticsvr.Service
	domainSvc    *domainsvr.Service
	whiteSvc     *whitesvr.Service
	showSvc      *showsvr.Service
	broadcastSvc *broadcastsvr.Service
)

type Server struct {
	// depend service
	AuthSvc *auth.Auth
	// self service
	PgSvr        *pluginsvr.Service
	PingSvr      *pingsvr.Service
	SideSvr      *sidesvr.Service
	VerSvc       *version.Service
	ParamSvc     *param.Service
	NtcSvc       *notice.Service
	SplashSvc    *splash.Service
	AuditSvc     *auditsvr.Service
	AbSvc        *absvr.Service
	ModuleSvc    *modulesvr.Service
	GuideSvc     *guidesvc.Service
	StaticSvc    *staticsvr.Service
	DomainSvc    *domainsvr.Service
	WhiteSvc     *whitesvr.Service
	ShowSvc      *showsvr.Service
	BroadcastSvc *broadcastsvr.Service
}

// Init is
func Init(c *conf.Config, svr *Server) {
	initService(c, svr)
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
func initService(c *conf.Config, svr *Server) {
	// init self service
	authSvc = svr.AuthSvc
	pgSvr = svr.PgSvr
	pingSvr = svr.PingSvr
	sideSvr = svr.SideSvr
	verSvc = svr.VerSvc
	paramSvc = svr.ParamSvc
	ntcSvc = svr.NtcSvc
	splashSvc = svr.SplashSvc
	auditSvc = svr.AuditSvc
	abSvc = svr.AbSvc
	moduleSvc = svr.ModuleSvc
	guideSvc = svr.GuideSvc
	staticSvc = svr.StaticSvc
	domainSvc = svr.DomainSvc
	broadcastSvc = svr.BroadcastSvc
	whiteSvc = svr.WhiteSvc
	showSvc = svr.ShowSvc
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	r := e.Group("/x/resource")
	{
		r.GET("/plugin", plugin)
		r.GET("/sidebar", authSvc.GuestMobile, sidebar)
		r.GET("/topbar", topbar)
		r.GET("/abtest", abTest)
		r.GET("/abtest/v2", abTestV2)
		r.GET("/abtest/abserver", authSvc.GuestMobile, abserver)
		m := r.Group("/module")
		{
			m.POST("", module)
			m.POST("/list", list)
		}
		g := r.Group("/guide", authSvc.GuestMobile)
		{
			g.GET("/interest", interest)
			g.GET("/interest2", interest2)
		}
		r.GET("/static", getStatic)
		r.GET("/domain", domain)
		r.GET("/broadcast/servers", serverList)
		r.GET("/white/list", whiteList)
		r.GET("/show/tab", authSvc.GuestMobile, tabs)
	}
	v := e.Group("/x/v2/version")
	{
		v.GET("", getVersion)
		v.GET("/update", versionUpdate)
		v.GET("/update.pb", versionUpdatePb)
		v.GET("/so", versionSo)
		v.GET("/rn/update", versionRn)
	}
	p := e.Group("/x/v2/param", authSvc.GuestMobile)
	{
		p.GET("", getParam)
	}
	n := e.Group("/x/v2/notice", authSvc.GuestMobile)
	{
		n.GET("", getNotice)
	}
	s := e.Group("/x/v2/splash")
	{
		s.GET("", splashs)
		s.GET("/birthday", birthSplash)
		s.GET("/list", authSvc.GuestMobile, splashList)
	}
	a := e.Group("/x/v2/audit")
	{
		a.GET("", audit)
	}
}
