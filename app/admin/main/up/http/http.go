package http

import (
	"go-common/app/admin/main/up/conf"
	"go-common/app/admin/main/up/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	//Svc service.
	Svc     *service.Service
	authSrc *permit.Permit
	idfSvc  *verify.Verify
)

// Init init account service.
func Init(c *conf.Config) {
	// service
	initService(c)
	// init internal router
	engine := bm.DefaultServer(c.HTTPServer)
	setupInnerEngine(engine)
	// init internal server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func initService(c *conf.Config) {
	idfSvc = verify.New(nil)
	Svc = service.New(c)
	authSrc = permit.New(c.Auth)
}

// innerRouter
func setupInnerEngine(e *bm.Engine) {
	// monitor ping
	e.Ping(ping)
	e.Register(disRegister)
	// base
	var adminUpProfit *bm.RouterGroup
	var noAdminUpProfit *bm.RouterGroup
	var identifyUpProfit *bm.RouterGroup
	if conf.Conf.IsTest {
		adminUpProfit = e.Group("/allowance/api/x/admin/uper")
	} else {
		// 现在只要登录，默认放过
		adminUpProfit = e.Group("/allowance/api/x/admin/uper", authSrc.Verify(), authSrc.Permit(""))
	}
	// 因为经常出现-401，所以把这些接口的验证去掉
	noAdminUpProfit = e.Group("/allowance/api/x/admin/uper")
	{
		//noAdminUpProfit.GET("/score/query", crmScoreQuery) // 这个接口需要干掉
		noAdminUpProfit.GET("/score/query_section", crmScoreQuery)
		noAdminUpProfit.GET("/score/query_up", crmScoreQueryUp)
		noAdminUpProfit.GET("/score/query_up_history", crmScoreQueryUpHistory)

		noAdminUpProfit.GET("/play/query", crmPlayQueryInfo)

		noAdminUpProfit.GET("/info/query", crmInfoQueryUp)
		noAdminUpProfit.GET("/info/account_info", crmInfoAccountInfo)
		noAdminUpProfit.POST("/info/search", crmInfoSearch)

		noAdminUpProfit.GET("/creditlog/query", crmCreditLogQueryUp)

		noAdminUpProfit.GET("/rank/query_list", crmRankQueryList)

		noAdminUpProfit.POST("/file/upload", upload)

		noAdminUpProfit.GET("/data/batch_query_data", crmQueryUpInfoWithViewerData)
		noAdminUpProfit.GET("/data/fan_summary", dataGetFanSummary)
		noAdminUpProfit.GET("/data/fan_relation_history", dataRelationFansHistory)
		noAdminUpProfit.GET("/data/up_archive_info", dataGetUpArchiveInfo)
		noAdminUpProfit.GET("/data/up_archive_tag_info", dataGetUpArchiveTagInfo)
		noAdminUpProfit.GET("/data/up_view_info", dataGetUpViewInfo)
	}

	if conf.Conf.IsTest {
		identifyUpProfit = e.Group("/allowance/api/x/admin/uper")
	} else {
		identifyUpProfit = e.Group("/allowance/api/x/admin/uper", idfSvc.Verify)
	}
	{
		identifyUpProfit.GET("/service/batch_query_data", crmQueryUpInfoWithViewerData)
		identifyUpProfit.GET("/service/data/fan_summary", dataGetFanSummary)
		identifyUpProfit.GET("/service/data/fan_relation_history", dataRelationFansHistory)
		identifyUpProfit.GET("/service/data/up_archive_info", dataGetUpArchiveInfo)
		identifyUpProfit.GET("/service/data/up_archive_tag_info", dataGetUpArchiveTagInfo)

		noAdminUpProfit.GET("/test/get_view_base", testGetViewBase)

	}
	dashboard := noAdminUpProfit.Group("/dashboard")
	{
		dashboard.GET("/yesterday", yesterday)
		dashboard.GET("/trend", trend)
		dashboard.GET("/trend/detail", trendDetail)
	}

	// sign 需要admin验证，这里需要admin的名字和id
	sign := adminUpProfit.Group("/sign")
	{
		sign.POST("/add", signAdd)
		sign.POST("/update", signUpdate)
		sign.POST("/violation/add", violationAdd)
		sign.POST("/violation/retract", violationRetract)
		sign.GET("/violation/list", violationList)
		sign.POST("/absence/add", absenceAdd)
		sign.POST("/absence/retract", absenceRetract)
		sign.GET("/absence/list", absenceList)
		sign.GET("/up/view/check", viewCheck)
		sign.GET("/query", signQuery)
		sign.GET("/query/id", signQueryID)
		sign.GET("/up/aduit/log", signUpAuditLogs)
		sign.GET("/country/list", countrys)
		sign.GET("/tid/list", tids)
		sign.POST("/pay/complete", signPayComplete)

	}
	signNoAdmin := noAdminUpProfit.Group("/sign")
	{
		signNoAdmin.GET("/check_exist", signCheckExist)
	}

	commandNoAdmin := noAdminUpProfit.Group("/command")
	{
		commandNoAdmin.GET("/refresh_up_rank", commandRefreshUpRank)
	}
	//{
	//	admin.GET("/special/get", specialGet)
	//	admin.GET("/special/get_by_mid", specialGetByMid)
	//	admin.POST("/special/delete", specialDel)
	//	admin.POST("/special/add", specialAdd)
	//	admin.POST("/special/edit", specialEdit)
	//	admin.GET("/group/get", getGroup)
	//	admin.POST("/group/add", authSrc.Permit("UPGROUP_ADD"), addGroup)
	//	admin.POST("/group/update", updateGroup)
	//	admin.POST("/group/delete", authSrc.Permit("UPGROUP_ADD"), removeGroup)
	//}

}

// ping check server ok.
func ping(ctx *bm.Context) {
	if err := Svc.Ping(ctx); err != nil {
		ctx.Error = err
		ctx.AbortWithStatus(503)
	}
}

// disRegister check server ok.
func disRegister(ctx *bm.Context) {
	ctx.JSON(map[string]interface{}{}, nil)
}
