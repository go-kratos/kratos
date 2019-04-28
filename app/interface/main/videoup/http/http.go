package http

import (
	"go-common/app/interface/main/videoup/conf"
	"go-common/app/interface/main/videoup/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	vdpSvc    *service.Service
	antiSvc   *antispam.Antispam
	verifySvc *verify.Verify
	authSvc   *auth.Auth
)

// Init fn
func Init(c *conf.Config, s *service.Service) {
	initService(c)
	engineOuter := bm.DefaultServer(c.BM)
	outerRouter(engineOuter)
	if err := engineOuter.Start(); err != nil {
		log.Error("engineOuter.Start() error(%v) | config(%v)", err, c)
		panic(err)
	}
}

// initService init service
func initService(c *conf.Config) {
	verifySvc = verify.New(nil)
	authSvc = auth.New(nil)
	vdpSvc = service.New(c)
	antiSvc = antispam.New(c.UpCoverAnti)
}

// outerRouter ForLogic port:6321
func outerRouter(e *bm.Engine) {
	e.GET("/monitor/ping", ping)
	g := e.Group("/x/vu")
	{
		client := g.Group("/client")
		{
			client.POST("/add", authSvc.UserMobile, clientAdd)
			client.POST("/edit", authSvc.UserMobile, clientEdit)
			client.POST("/cover/up", authSvc.UserMobile, antiSvc.Handler(), clientUpCover)
		}
		web := g.Group("/web")
		{
			web.POST("/add", authSvc.UserWeb, webAdd)
			web.POST("/edit", authSvc.UserWeb, webEdit)
			web.POST("/cover/up", authSvc.UserWeb, antiSvc.Handler(), webUpCover)
			web.POST("/filter", authSvc.UserWeb, webFilter)
			web.POST("/staff-title/filter", authSvc.UserWeb, webStaffTitleFilter)
			web.POST("/cm/add", authSvc.UserWeb, webCmAdd)
			web.POST("/v2/add", authSvc.UserWeb, webV2Add)
		}
		app := g.Group("/app")
		{
			app.POST("/edit", authSvc.UserMobile, appEdit)
			//new feature
			app.POST("/add", authSvc.UserMobile, appAdd)
			app.POST("/edit/full", authSvc.UserMobile, appEditFull)
			app.POST("/cover/up", authSvc.UserMobile, antiSvc.Handler(), appUpCover)
		}
		creator := g.Group("/creator")
		{
			creator.POST("/add", authSvc.UserMobile, creatorAdd)
			creator.POST("/edit", authSvc.UserMobile, creatorEdit)
			creator.POST("/cover/up", authSvc.UserMobile, antiSvc.Handler(), creatorUpCover)
		}
	}
}

func getBuildInfo(c *bm.Context) (build, buvid string) {
	buvid = c.Request.Header.Get("Buvid")
	if buvid == "" {
		buvidCookie, _ := c.Request.Cookie("buvid3")
		if buvidCookie != nil {
			buvid = buvidCookie.Value
		}
	}
	build = c.Request.Form.Get("build")
	return
}
