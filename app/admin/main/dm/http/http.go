package http

import (
	"net/http"

	"go-common/app/admin/main/dm/conf"
	"go-common/app/admin/main/dm/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	//dmSvc service
	dmSvc     *service.Service
	verifySvc *verify.Verify
	// authSvc auth service
	authSvc *permit.Permit
)

// Init http init
func Init(c *conf.Config, s *service.Service) {
	dmSvc = s
	verifySvc = verify.New(c.Verify)
	authSvc = permit.New(c.ManagerAuth)
	engine := bm.DefaultServer(c.HTTPServer)
	authRouter(engine)
	innerRouter(engine)
	// init internal server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func authRouter(engine *bm.Engine) {
	engine.GET("/monitor/ping", ping)
	group := engine.Group("/x/admin/dm")
	{
		// up filter
		group.GET("/upfilter/list", authSvc.Permit("DANMU_POOL_OPERATION"), upFilters)
		group.POST("/upfilter/edit", authSvc.Permit("DANMU_POOL_OPERATION"), editUpFilters)
		// advance dm list
		group.GET("/adv/list", authSvc.Permit("DANMU_POOL_OPERATION"), advList)
		// dm transfer
		group.GET("/transfer/list", authSvc.Permit("DANMU_POOL_OPERATION"), transferList)
		group.POST("/transfer/retry", authSvc.Permit("DANMU_POOL_OPERATION"), reTransferJob)
		// dm subject
		group.POST("/subject/state/edit", authSvc.Permit("DANMAKU_POOL_SWITCH"), uptSubjectsState)
		group.GET("/subject/log", authSvc.Permit("DANMU_LIST"), subjectLog)
		group.GET("/subject/archive", authSvc.Permit("DANMU_LIST"), archiveList)
		group.POST("/subject/maxlimit", authSvc.Permit("DANMU_POOL_OPERATION"), upSubjectMaxLimit)
		// dm mask
		group.GET("/mask/state", authSvc.Permit("DANMU_LIST"), maskState)
		group.POST("/mask/state/update", authSvc.Permit("DANMAKU_MASK_SWITCH"), updateMaskState)
		group.POST("/mask/generate", authSvc.Permit("DANMAKU_MASK_SWITCH"), generateMask)
		group.GET("/mask/up", authSvc.Permit("DANMU_LIST"), maskUps)
		group.POST("/mask/up/open", authSvc.Permit("DANMAKU_MASK_SWITCH"), maskUpOpen)
		// dm task
		group.GET("/task/list", authSvc.Permit("DM_TASK_LIST"), taskList)
		group.POST("/task/new", authSvc.Permit("DM_TASK_OPERATION"), addTask)
		group.POST("/task/review", authSvc.Permit("DM_TASK_REVIEW"), reviewTask)
		group.POST("/task/state/edit", authSvc.Permit("DM_TASK_OPERATION"), editTaskState)
		group.GET("/task/view", authSvc.Permit("DM_TASK_LIST"), taskView)
		group.GET("/task/csv", authSvc.Permit("DM_TASK_LIST"), taskCsv)
		// dm list
		group.GET("/content/list", authSvc.Permit("DANMU_LIST"), contentList)
		group.GET("/report/list/first", authSvc.Permit("DM_REPORT_FIRST_READ"), reportList2)
		group.GET("/report/list/second", authSvc.Permit("DM_REPORT_SECOND_READ"), reportList2)

		// dm bnj shield
		group.POST("/shield/upload", authSvc.Permit("DANMU_LIST"), shieldUpload)
	}
	subtitleG := group.Group("/subtitle")
	{
		subtitleG.GET("/list", authSvc.Permit("DM_SUBTITLE"), subtitleList)
		subtitleG.POST("/edit", authSvc.Permit("DM_SUBTITLE"), subtitleEdit)
		subtitleG.POST("/subject/switch", authSvc.Permit("DM_SUBTITLE"), subtitleSwitch)
	}
}

// innerRouter init inner router.
func innerRouter(engine *bm.Engine) {
	group := engine.Group("/x/internal/dmadmin", verifySvc.Verify)
	{
		group.POST("/trans/add", addTrJob)
		cg := group.Group("/content")
		{
			cg.POST("/edit/state", editDMState)
			cg.POST("/edit/pool", editDMPool)
			cg.POST("/edit/attr", editDMAttr)
			cg.GET("/list", dmSearch)
			cg.POST("/refresh", xmlCacheFlush)
			cg.GET("/index/info", dmIndexInfo)
			cg.GET("/log/query", logList)
		}
		sg := group.Group("/subject")
		{
			sg.GET("/info", dmIndexInfo)
			sg.POST("/fix/count", fixDMCount)
		}
		rg := group.Group("/report")
		{
			rg.GET("/list", reportList)
			rg.GET("/log", reportLog)
			rg.POST("/user/stat/change", changeReportUserStat)
			rg.POST("/stat/change", changeReportStat)
			rg.POST("/judge", transferJudge)
			rg.POST("/judge/result", JudgeResult)
		}
		mg := group.Group("/monitor")
		{
			mg.GET("/list", monitorList)
			mg.POST("/edit", editMonitor)
		}
		subtitleG := group.Group("/subtitle")
		{
			subtitleG.POST("/workflow/callback", subtitleEditCallback)
			subtitleG.GET("/workflow/status/list", subtitleStatusList)
			subtitleG.GET("/workflow/lans/list", subtitleLanList)
		}
	}
}

// ping check server state.
func ping(ctx *bm.Context) {
	if err := dmSvc.Ping(ctx); err != nil {
		log.Error("dm admin ping error(%v)", err)
		ctx.JSON(nil, err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
	ctx.Next()
}
