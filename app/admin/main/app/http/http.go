package http

import (
	"go-common/app/admin/main/app/conf"
	aidssvr "go-common/app/admin/main/app/service/aids"
	auditsvr "go-common/app/admin/main/app/service/audit"
	bfssvr "go-common/app/admin/main/app/service/bfs"
	bottomsvr "go-common/app/admin/main/app/service/bottom"
	langsvr "go-common/app/admin/main/app/service/language"
	noticesvr "go-common/app/admin/main/app/service/notice"
	pingsvr "go-common/app/admin/main/app/service/ping"
	wallsvr "go-common/app/admin/main/app/service/wall"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	authSvc   *permit.Permit
	auditSvc  *auditsvr.Service
	noticeSvc *noticesvr.Service
	langSvc   *langsvr.Service
	wallSvc   *wallsvr.Service
	bottomSvc *bottomsvr.Service
	aidsSvc   *aidssvr.Service
	pingSvc   *pingsvr.Service
	bfsSvc    *bfssvr.Service
)

// Init init
func Init(c *conf.Config) {
	initService(c)
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	authSvc = permit.New(c.Auth)
	auditSvc = auditsvr.New(c)
	noticeSvc = noticesvr.New(c)
	langSvc = langsvr.New(c)
	wallSvc = wallsvr.New(c)
	bottomSvc = bottomsvr.New(c)
	aidsSvc = aidssvr.New(c)
	pingSvc = pingsvr.New(c)
	bfsSvc = bfssvr.New(c)
}

func innerRouter(e *bm.Engine) {
	e.GET("/monitor/ping", moPing)
	b := e.Group("/x/admin/app", authSvc.Verify())
	{
		b.POST("/upload/cover", clientUpCover)
		cb := b.Group("/audit", authSvc.Permit(authRouter("audit")))
		{
			cb.GET("", audits)
			cb.GET("/detail", auditByID)
			cb.POST("/save", auditSave)
			cb.POST("/del", auditDelByIDs)
		}
		cb = b.Group("/notice", authSvc.Permit(authRouter("notice")))
		{
			cb.GET("", notices)
			cb.GET("/detail", noticeByID)
			cb.POST("/save", addOrupdate)
			cb.POST("/modifybuild", updateBuild)
			cb.POST("/modifystate", updateState)
		}
		cb = b.Group("/language", authSvc.Permit(authRouter("language")))
		{
			cb.GET("", languages)
			cb.GET("/detail", langByID)
			cb.POST("/save", addOrup)
		}
		cb = b.Group("/wall", authSvc.Permit(authRouter("wall")))
		{
			cb.GET("", walls)
			cb.GET("/detail", wallByID)
			cb.POST("/save", saveWall)
			cb.POST("/publish", publish)
			cb.POST("/publishtest", publish)
		}
		cb = b.Group("/bottom", authSvc.Permit(authRouter("bottom")))
		{
			cb.GET("", bottoms)
			cb.GET("/detail", bottomByID)
			cb.POST("/save", bottomInsert)
			cb.POST("/publish", publishBottom)
			cb.POST("/publishtest", publishBottom)
			cb.POST("/delbottom", delBottom)
		}
		cb = b.Group("/aids", authSvc.Permit(authRouter("aids")))
		{
			cb.POST("/save", saveAids)
		}
	}
}

func authRouter(name string) string {
	if perm, ok := conf.Conf.Perms.Perm[name]; ok {
		return perm
	}
	return ""
}
