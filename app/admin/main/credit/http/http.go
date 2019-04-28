package http

import (
	"net/http"

	"go-common/app/admin/main/credit/conf"
	"go-common/app/admin/main/credit/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	authSvc *permit.Permit
	creSvc  *service.Service
)

// Init http server
func Init(c *conf.Config) {
	initService(c)
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	// init internal server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// initService init service
func initService(c *conf.Config) {
	creSvc = service.New(c)
	authSvc = permit.New(c.Auth)
}

// innerRouter
func innerRouter(e *bm.Engine) {
	// ping monitor
	e.Ping(ping)
	// internal api
	bg := e.Group("/x/admin/credit/blocked")
	{
		// info
		bg.GET("/info", authSvc.Permit("BLACK_HOUSE_BLOCKED_INFO"), infos)                             // 封禁列表
		bg.GET("/info/id", authSvc.Permit("BLACK_HOUSE_BLOCKED_INFO_ID"), infoByID)                    // 封禁信息详情
		bg.POST("/info/up", authSvc.Permit("BLACK_HOUSE_BLOCKED_INFO_UP"), upInfo)                     // 编辑封禁列表
		bg.POST("/info/del", authSvc.Permit("BLACK_HOUSE_BLOCKED_INFO_DEL"), delInfo)                  // 删除封禁列表
		bg.POST("/info/status/up", authSvc.Permit("BLACK_HOUSE_BLOCKED_INFO_STATUS_UP"), upInfoStatus) // 更新封禁列表状态
		bg.GET("/info.so", authSvc.Permit("BLACK_HOUSE_BLOCKED_INFO_SO"), infosEx)                     // 封禁列表导出
		// publish
		bg.GET("/publish", authSvc.Permit("BLACK_HOUSE_BLOCKED_PUBLISH"), publishs)            // 公告列表
		bg.GET("/publish/id", authSvc.Permit("BLACK_HOUSE_BLOCKED_PUBLISH_ID"), publishByID)   // 公告信息详情
		bg.POST("/publish/add", authSvc.Permit("BLACK_HOUSE_BLOCKED_PUBLISH_ADD"), addPublish) // 增加公告
		bg.POST("/publish/up", authSvc.Permit("BLACK_HOUSE_BLOCKED_PUBLISH_UP"), upPublish)    // 更新公告
		bg.POST("/publish/del", authSvc.Permit("BLACK_HOUSE_BLOCKED_PUBLISH_DEL"), delPublish) // 删除公告
		// notice
		bg.GET("/notice", authSvc.Permit("BLACK_HOUSE_BLOCKED_NOTICE"), notices)                             // 通知列表
		bg.POST("/notice/add", authSvc.Permit("BLACK_HOUSE_BLOCKED_NOTICE_ADD"), addNotice)                  // 增加通知
		bg.POST("/notice/status/up", authSvc.Permit("BLACK_HOUSE_BLOCKED_NOTICE_STATUS_UP"), upNoticeStatus) // 更新通知
	}
	jg := e.Group("/x/admin/credit/jury")
	{
		// case
		jg.GET("/case", authSvc.Permit("BLACK_HOUSE_JURY_CASE"), cases)                 // 案件列表
		jg.GET("/case/id", authSvc.Permit("BLACK_HOUSE_JURY_CASE_ID"), caseByID)        // 案件信息详情
		jg.GET("/case/reason", authSvc.Permit("BLACK_HOUSE_JURY_CASE_REASON"), reasons) // 案件举报理由
		jg.POST("/case/add", authSvc.Permit("BLACK_HOUSE_JURY_CASE_ADD"), addCase)
		jg.POST("/case/up", authSvc.Permit("BLACK_HOUSE_JURY_CASE_UP"), upCase)
		jg.POST("/case/vote/add", authSvc.Permit("BLACK_HOUSE_JURY_CASE_VOTE_ADD"), addCaseVote)
		jg.POST("/case/status/up", authSvc.Permit("BLACK_HOUSE_JURY_CASE_STATUS_UP"), upCaseStatus)
		jg.POST("/case/type/add", authSvc.Permit("BLACK_HOUSE_JURY_CASE_TYPE_ADD"), addCaseType) // 添加众裁稿件
		jg.GET("/case/auto/conf", authSvc.Permit("BLACK_HOUSE_JURY_CASE_AUTO_CONF"), autoCaseConfig)
		jg.POST("/case/auto/conf/set", authSvc.Permit("BLACK_HOUSE_JURY_CASE_AUTO_CONF_SET"), setAutoCaseConfig)
		// opinions
		jg.GET("/opinion", authSvc.Permit("BLACK_HOUSE_JURY_OPINION"), opinions)          // 观点列表
		jg.GET("/opinion/id", authSvc.Permit("BLACK_HOUSE_JURY_OPINION_ID"), opinionByID) // 观点信息详情
		jg.POST("/opinion/del", authSvc.Permit("BLACK_HOUSE_JURY_OPINION_DEL"), delOpinions)
		// users
		jg.GET("/users", authSvc.Permit("BLACK_HOUSE_JURY_USERS"), users)          // 委员列表
		jg.GET("/user/id", authSvc.Permit("BLACK_HOUSE_JURY_USER_ID"), userByID)   // 委员信息详情
		jg.POST("/user/add", authSvc.Permit("BLACK_HOUSE_JURY_USER_ADD"), userAdd) // 新增委员
		jg.POST("/users/status/up", authSvc.Permit("BLACK_HOUSE_JURY_USERS_STATUS_UP"), upUserStatus)
		jg.POST("/users/blackwhite", authSvc.Permit("BLACK_HOUSE_JURY_USERS_BLACKWHITE"), blackWhite)
		jg.GET("/users.so", authSvc.Permit("BLACK_HOUSE_JURY_USERS_SO"), usersEx)
		// kpi
		jg.GET("/kpi", authSvc.Permit("BLACK_HOUSE_JURY_KPI"), kpis)
		jg.GET("/kpi/point", authSvc.Permit("BLACK_HOUSE_JURY_KPI_POINT"), kpiPoints)
		jg.GET("/kpi.so", authSvc.Permit("BLACK_HOUSE_JURY_KPI_SO"), kpisEx)
		// config
		jg.GET("/config", authSvc.Permit("BLACK_HOUSE_JURY_CONFIG"), caseConf)
		jg.POST("/config/set", authSvc.Permit("BLACK_HOUSE_JURY_CONFIG_SET"), setCaseConf)
		jg.GET("/votenum/conf", authSvc.Permit("BLACK_HOUSE_VOTENUM_CONFIG"), votenumConf)
		jg.POST("/votenum/conf/set", authSvc.Permit("BLACK_HOUSE_VOTENUM_CONFIG_SET"), setVotenumConf)

	}
	ug := e.Group("/x/admin/credit/upload")
	{
		ug.POST("", upload)
		ug.POST("/coins", annualCoins)
	}
	lg := e.Group("/x/admin/credit/labour")
	{
		/*
			lg.GET("/quest", question)
			lg.POST("/quest/statistics", statQuestion)
		*/
		lg.POST("/quest/oper", operQuestion)
		lg.POST("/quest/del", delQuestion)
	}
	ig := e.Group("/x/admin/credit/jury")
	{
		// appeal webhook for internal .
		ig.POST("/webhook", webHook)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	err := creSvc.Ping(c)
	if err != nil {
		log.Error("credit admin ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
