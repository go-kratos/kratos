package http

import (
	"net/http"

	"go-common/app/admin/ep/melloi/conf"
	"go-common/app/admin/ep/melloi/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	srv *service.Service
	// AuthSrv Auth Service
	AuthSrv *permit.Permit
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	srv = s
	AuthSrv = permit.New2(c.PermitGRPC)
	engine := bm.NewServer(c.BM)
	engine.Ping(ping)
	router(engine)
	if err := engine.Start(); err != nil {
		log.Error("xhttp.Serve2 error(%v)", err)
		panic(err)
	}
}

// router init outer router api path
func router(e *bm.Engine) {
	e.Register(register)
	version := e.Group("/ep/admin/melloi/v1", AuthSrv.Permit2(""))
	{
		home := version.Group("/rank")
		{
			home.GET("/trees/query", treesQuery)
			home.GET("/top/http/query", topHttpQuery)
			home.GET("/top/dept/query", topDeptQuery)
			home.GET("/treenum/query", treeNumQuery)
			home.GET("/buildline/query", buildLineQuery)
			home.GET("/stateline/query", stateLineQuery)
			home.GET("/top/grpc/query", topGrpcQuery)
			home.GET("/top/scene/query", topSceneQuery)
		}
		user := version.Group("/user")
		{
			user.GET("/query", queryUser)
			user.POST("/update", updateUser)
		}
		//手工压测
		order := version.Group("/order")
		{
			order.GET("/query", queryOrder)
			order.POST("/add", addOrder)
			order.POST("/update", updateOrder)
			order.GET("/del", delOrder)
			order.POST("/report/add", addOrderReport)
			order.GET("/report/query", queryOrderReport)
		}
		//非压测时间段压测申请
		apply := version.Group("/apply")
		{
			apply.GET("/query", queryApply)
			apply.POST("/add", addApply)
			apply.POST("/update", updateApply)
			apply.GET("/del", delApply)
		}
		orderAdmin := version.Group("/admin")
		{
			orderAdmin.GET("/query", queryOrderAdmin)
			orderAdmin.GET("/add", addOrderAdmin)
		}
		script := version.Group("/script")
		{
			script.POST("/gettgroup", getThreadGroup)
			script.POST("/add", addAndExecuteScript)
			script.POST("/grpc/add", grpcAddScript)
			script.POST("/update", updateScript)
			script.POST("/updateall", updateScriptAll)
			script.GET("/delete", delScript)
			script.GET("/query", queryScripts)
			script.GET("/query/free", queryScripysFree)
			script.GET("/querysnap", queryScriptSnap)
			script.POST("/url/check", urlCheck)
			script.POST("/addtimer", addTimer)
		}

		scene := version.Group("/scene")
		{
			scene.POST("/addscene/auto", AddAndExecuScene)
			scene.GET("/draft/query", queryDraft)
			scene.POST("/update", updateScene)
			scene.POST("/add", addScene)
			scene.POST("/save", saveScene)
			scene.POST("/saveorder", saveOrder)
			scene.GET("/relation/query", queryRelation)
			scene.GET("/api/query", queryAPI)
			scene.POST("/api/delete", deleteAPI)
			scene.POST("/group/config/add", addConfig)
			scene.GET("/tree/query", queryTree)
			scene.GET("/query", queryScenes)
			scene.POST("/doptest", doScenePtest)
			scene.POST("/doptest/batch", doScenePtestBatch)
			scene.GET("/existapi/query", queryExistAPI)
			scene.GET("/preview/query", queryPreview)
			scene.GET("/params/query", queryParams)
			scene.POST("/bindScene/update", updateBindScene)
			scene.GET("/drawrelation/query", queryDrawRelation)
			scene.POST("/script/add", addSceneScript)
			scene.POST("/draft/delete", deleteDraft)
			scene.GET("/group/config/query", queryConfig)
			scene.GET("/delete", deleteScene)
			scene.POST("/copy", copyScene)
			scene.GET("/fusing/query", queryFusing)
		}

		run := version.Group("/run")
		{
			run.GET("/time/check", runTimeCheck)
			run.POST("/grpc/run", runGrpc)
			run.GET("/grpc/query", queryGrpc)
			run.POST("/grpc/update", updateGrpc)
			run.GET("/grpc/delete", deleteGrpc)
			run.GET("/grpc/snap/query", queryGrpcSnap)
			run.GET("/proto/query", getProto)
		}

		file := version.Group("/file")
		{
			file.POST("/upload", upload)
			file.POST("/proto/refer/upload", uploadDependProto)
			file.GET("/read", readFile)
			file.GET("/download", downloadFile)
			file.GET("/exists", isFileExists)
			file.POST("/img", uploadImg)
			file.POST("/proto/compile", compileProtoFile)
		}
		tree := version.Group("/tree")
		{
			tree.GET("/query", queryUserTree)
			tree.GET("/admin", queryTreeAdmin)
		}
		report := version.Group("report")
		{
			report.POST("/update/summary", updateReportSummary)
			report.GET("/query", queryReportSummarys)
			report.POST("/query/regraph", queryReGraph)
			report.POST("/query/regraphavg", queryReGraphAvg)
			report.GET("/update/status", updateReportStatus)
			report.GET("/del", delReport)
			report.GET("/query/id", queryReportByID)
		}
		moni := version.Group("moni")
		{
			moni.GET("/client", queryClientMoni)
		}

		ptest := version.Group("ptest")
		{
			ptest.GET("/add", addPtest)
			ptest.GET("/reduce", reducePtest)
			ptest.GET("/queryalljob", queryAllJob)
			ptest.GET("/queryalljob/free", queryAllJobFree)
			ptest.GET("/stop", stopPtest)
			ptest.GET("/stopall", stopAllPtest)
			ptest.GET("/doptest", doPtest)
			ptest.GET("/doptest/scriptid", doPtestByScriptId)
			ptest.POST("/grpc/quickstart", grpcQuickStart)
			ptest.POST("/grpc/qksave", saveGrpc)
			ptest.POST("/grpc/run", runGrpcByScriptId)
			ptest.POST("/grpc/createpath", createDependencyPath)
			ptest.POST("/doptestfile", doPtestByFile)
			ptest.POST("/dobatchptest", doPtestBatch)
			ptest.POST("/dodebug", doDebug)
		}
		cluster := version.Group("/cluster")
		{
			cluster.GET("/info", ClusterInfo)
			job := version.Group("/job")
			{
				job.POST("/add", addJob)
				job.POST("/delete/batch", deleteJobBatch)
				job.DELETE("/delete", deleteJob)
				job.GET("/get", Job)
				job.GET("/clean/query", queryClearnableDocker)
				job.GET("/clean", cleanNotRunningJob)
			}
		}
		comment := version.Group("/comment")
		{
			comment.GET("/query", queryComment)
			comment.POST("/add", addComment)
			comment.POST("/update", updateComment)
			comment.GET("/delete", deleteComment)
		}

		label := version.Group("/label")
		{
			label.GET("/query", queryLabels)
			label.POST("/add", addLabel)
			label.GET("/del", delLabel)
			label.POST("/relation/add", addLabelRelation)
			label.GET("/relation/del", delLabelRelation)
		}
	}
	versiond := e.Group("/ep/admin/melloi/v2")
	{
		script := versiond.Group("script")
		{
			script.POST("/addsample", addJmeterSample)
			script.POST("/addtgroup", addThreadGroup)
		}
		file := versiond.Group("/file")
		{
			file.GET("/download", downloadFile)
		}
		job := versiond.Group("/job")
		{
			job.GET("/force/delete", forceDelete)
		}
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
