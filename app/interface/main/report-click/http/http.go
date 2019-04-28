package http

import (
	"time"

	"go-common/app/interface/main/report-click/conf"
	"go-common/app/interface/main/report-click/service"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	clickSvr        *service.Service
	authSvc         *auth.Auth
	verifySvc       *verify.Verify
	infocRealTime   *infoc.Infoc
	infocStatistics *infoc.Infoc
	fromMap         = make(map[int64]bool)
	fromInlineMap   = make(map[int64]bool)
	inlineDuration  int64
)

// New http init.
func New(c *conf.Config) (engine *bm.Engine) {
	clickSvr = service.New(c)
	authSvc = auth.New(c.Auth)
	verifySvc = verify.New(c.Verify)
	infocRealTime = infoc.New(c.Infoc2.RealTime)
	infocStatistics = infoc.New(c.Infoc2.Statistics)
	for _, v := range c.Click.From {
		fromMap[v] = true
	}
	for _, v := range c.Click.FromInline { // init inline play "from"
		fromInlineMap[v] = true
	}
	inlineDuration = c.Click.InlineDuration // inline play duration line
	engine = bm.DefaultServer(c.BM)
	engine.Use(bm.Recovery(), bm.Logger())
	outerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start() error(%v)", err)
		panic(err)
	}
	return
}

func outerRouter(e *bm.Engine) {
	e.GET("/monitor/ping", ping)
	e.POST("/x/report/click/web", authSvc.GuestWeb, webClick)
	e.POST("/x/report/click/outer", authSvc.GuestWeb, outerClick)
	e.POST("/x/stat/web", authSvc.GuestWeb, webClick)
	e.POST("/x/stat/outer", authSvc.GuestWeb, outerClick)
	click := e.Group("/x/report/click")
	{
		click.GET("/now", serverNow)
		click.POST("/h5", authSvc.Guest, h5Click)
		click.POST("/h5/outer", authSvc.Guest, outerClickH5) // nocsrf
		click.POST("/ios", authSvc.Guest, iosClick)
		click.POST("/android", authSvc.Guest, androidClick)
		click.POST("/android2", authSvc.Guest, android2Click)
		click.POST("/web/h5", authSvc.Guest, webH5Click)
		click.POST("/android/tv", authSvc.Guest, androidTV)
	}
	report := e.Group("/x/report/")
	{
		report.POST("/player", verifySvc.Verify, reportPlayer)       // old 30s heart
		report.POST("/heartbeat", verifySvc.Verify, reportHeartbeat) // new app 30s heart
		report.POST("/heartbeat/mobile", verifySvc.Verify, heartbeatMobile)
		report.POST("/web/heartbeat", authSvc.Guest, webHeartbeat) // web 30s heart

	}
	stat := e.Group("/x/stat")
	{
		stat.GET("/now", serverNow)
		stat.POST("/err_report", errReport)
		stat.POST("/h5", authSvc.Guest, h5Click)
		stat.POST("/ios", authSvc.Guest, iosClick)
		stat.POST("/android", authSvc.Guest, androidClick)
		stat.POST("/android2", authSvc.Guest, android2Click)
	}
}

func ping(c *bm.Context) {}
func serverNow(c *bm.Context) {
	data := map[string]int64{"now": time.Now().Unix()}
	c.JSON(data, nil)
}
