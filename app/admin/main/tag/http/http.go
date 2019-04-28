package http

import (
	"net/http"

	"go-common/app/admin/main/tag/conf"
	"go-common/app/admin/main/tag/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	authSvc *permit.Permit
	svc     *service.Service
	vfySvc  *verify.Verify
)

// Init http server .
func Init(c *conf.Config, s *service.Service) {
	svc = s
	authSvc = permit.New(c.Auth)
	vfySvc = verify.New(c.Verify)
	engine := bm.DefaultServer(c.HTTPServer)
	router(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func router(e *bm.Engine) {
	e.Ping(ping)
	group := e.Group("/x/admin/tag")
	{
		// tag action.
		group.GET("list", tagList)
		group.POST("edit", authSvc.Permit(authRouter("manager")), tagEdit)
		group.GET("info", tagInfo)
		group.POST("state", authSvc.Permit(authRouter("manager")), tagState)
		group.POST("verify", authSvc.Permit(authRouter("manager")), tagVerify)
		group.GET("/check", tagCheck)

		// resource action.
		group.GET("/resource/log/list", authSvc.Permit(authRouter("log")), resourceLogList)
		group.POST("/resource/log/state", authSvc.Permit(authRouter("log")), resourceLogState)
		group.GET("/resource/limit/info", authSvc.Permit(authRouter("videoLock")), resourceLimit)
		group.POST("/resource/limit/state", authSvc.Permit(authRouter("videoLock")), resourceLimitState)

		// user limit action.
		group.GET("/user/limit/list", authSvc.Permit(authRouter("whiteList")), limitUserList)
		group.POST("/user/limit/add", authSvc.Permit(authRouter("whiteList")), limitUserAdd)
		group.POST("/user/limit/delete", authSvc.Permit(authRouter("whiteList")), limitUserDel)

		// synonym action.
		group.GET("/synonym/list", authSvc.Permit(authRouter("synonym")), synonymList)
		group.POST("/synonym/edit", authSvc.Permit(authRouter("synonym")), synonymEdit)
		group.GET("/synonym/info", authSvc.Permit(authRouter("synonym")), synonymInfo)
		group.POST("/synonym/delete", authSvc.Permit(authRouter("synonym")), synonymDel)
		group.POST("/synonym/exist", authSvc.Permit(authRouter("synonym")), synonymIsExist)

		// relation action.
		group.GET("/relation/list", authSvc.Permit(authRouter("relation")), relationList)
		group.POST("/relation/add", authSvc.Permit(authRouter("relation")), relationAdd)
		group.POST("/relation/lock", authSvc.Permit(authRouter("relation")), relationLock)
		group.POST("/relation/unlock", authSvc.Permit(authRouter("relation")), relationUnLock)
		group.POST("/relation/delete", authSvc.Permit(authRouter("relation")), relationDelete)

		// report action.
		group.GET("/report/list", authSvc.Permit(authRouter("report")), reportList)
		group.GET("/report/info", authSvc.Permit(authRouter("report")), reportInfo)
		group.POST("/report/handle", authSvc.Permit(authRouter("report")), reportHandle)
		group.POST("/report/state", authSvc.Permit(authRouter("report")), reportState)
		group.POST("/report/ignore", authSvc.Permit(authRouter("report")), reportIgnore)
		group.POST("/report/delete", authSvc.Permit(authRouter("report")), reportDelete)
		group.POST("/report/punishs", authSvc.Permit(authRouter("report")), reportPunish)
		group.GET("/report/log/list", authSvc.Permit(authRouter("report")), reportLogList)
		group.GET("/report/log/info", authSvc.Permit(authRouter("report")), reportLogInfo)

		// hot action.
		group.GET("/hot/list", authSvc.Permit(authRouter("hot")), hotList)
		group.GET("/hot/archive/list", authSvc.Permit(authRouter("hot")), hotArchiveList)
		group.POST("/hot/operate", authSvc.Permit(authRouter("hot")), hotOperate)
		group.POST("/hot/update", authSvc.Permit(authRouter("hot")), updateHotTag)
		group.POST("/region/archive/refresh", authSvc.Permit(authRouter("hot")), regionArcRefresh)

		// business
		group.GET("/business/list", listBusiness)
		group.GET("/business/get", getBusiness)
		group.POST("/business/add", addBusiness)
		group.POST("/business/update", upBusiness)
		group.POST("/business/state", upBusiState)
	}
	groupChan := e.Group("/x/admin/channel")
	{
		groupChan.GET("/card/list", vfySvc.Verify, channeList) // APP卡片管理后台使用
		groupChan.GET("/card/all", vfySvc.Verify, channelAll)  // APP卡片管理后台使用
		groupChan.GET("/card/category/list", vfySvc.Verify, categoryList)

		groupChan.GET("/all", channelAll)
		groupChan.GET("/list", channeList)
		groupChan.GET("/info", channelInfo)
		groupChan.POST("/edit", authSvc.Permit(authRouter("channel")), channelEdit)
		groupChan.POST("/delete", authSvc.Permit(authRouter("channel")), channelDelete)
		groupChan.POST("/state", authSvc.Permit(authRouter("channel")), channelState)
		groupChan.POST("/migrate", authSvc.Permit(authRouter("channel")), migrateChannel)
		groupChan.POST("/sort", authSvc.Permit(authRouter("channel")), sortChannel)
		groupChan.POST("/shield/international", authSvc.Permit(authRouter("channel")), channelShieldINT)

		groupChan.GET("/recommend", recommandChannel)
		groupChan.POST("/recommend/sort", authSvc.Permit(authRouter("channel")), sortRecommandChannel)
		groupChan.POST("/recommend/migrate", authSvc.Permit(authRouter("channel")), migrateRecommandChannel)

		groupChan.GET("/category/list", categoryList)
		groupChan.POST("/category/add", authSvc.Permit(authRouter("channel")), categoryAdd)
		groupChan.POST("/category/delete", authSvc.Permit(authRouter("channel")), categoryDelete)
		groupChan.GET("/category/channels", categoryChannels)
		groupChan.POST("/category/sort", authSvc.Permit(authRouter("channel")), categorySort)
		groupChan.POST("/category/shield/international", authSvc.Permit(authRouter("channel")), categoryShieldINT)

		groupChan.GET("/rule/check", channelRuleCheck)
	}
}

func ping(c *bm.Context) {
	if svc.Ping(c) != nil {
		log.Error("tag-admin ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func authRouter(name string) string {
	if perm, ok := conf.Conf.Perms[name]; ok {
		return perm
	}
	return ""
}

func managerInfo(c *bm.Context) (uid int64, username string) {
	if nameInter, ok := c.Get("username"); ok {
		username = nameInter.(string)
	}
	if uidInter, ok := c.Get("uid"); ok {
		uid = uidInter.(int64)
	}
	return
}
