package http

import (
	"net/http"

	"go-common/app/admin/ep/saga/conf"
	"go-common/app/admin/ep/saga/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	srv     *service.Service
	authSvc *permit.Permit
)

// Init init
func Init(s *service.Service) {
	srv = s
	authSvc = permit.New2(nil)

	engine := bm.DefaultServer(conf.Conf.BM)
	engine.Ping(ping)
	initRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// initRouter init outer router api path.
func initRouter(e *bm.Engine) {
	version := e.Group("/ep/admin/saga/v1", authSvc.Permit2(""))
	{
		project := version.Group("/projects")
		{
			project.GET("/favorite", favoriteProjects)
			project.POST("/favorite/edit", editFavorite)
			project.GET("/common", queryCommonProjects)
		}

		tasks := version.Group("/tasks")
		{
			tasks.GET("/project", projectTasks)
		}

		user := version.Group("/user")
		{
			user.GET("/query", queryUserInfo)
		}

		data := version.Group("/data")
		{
			data.GET("/teams", queryTeams)
			data.GET("/project", queryProjectInfo)
			data.GET("/project/commit", queryProjectCommit)
			data.GET("/project/mr", queryProjectMr)
			data.GET("/commit", queryCommit) // ignore
			data.GET("/commit/report", queryTeamCommit)
			data.GET("/mr/report", queryTeamMr)
			data.GET("/pipeline/report", queryTeamPipeline)

			data.GET("/project/pipelines", queryProjectPipelineLists)
			data.GET("/project/branch", queryProjectBranchList)
			data.GET("/project/members", queryProjectMembers)
			data.GET("/project/status", queryProjectStatus)
			data.GET("/project/query/types", queryProjectTypes)
			data.GET("/project/runners", queryProjectRunners)
			data.GET("/job/report", queryProjectJob)
			data.GET("/project/mr/report", queryProjectMrReport)
			data.GET("/branch/report", queryBranchDiffWith)
		}

		config := version.Group("/config")
		{
			config.GET("/whitelist", sagaUserList)

			//get runner sven all config files
			config.GET("", runnerConfig)

			//get saga sven all config files
			config.GET("/saga", sagaConfig)
			config.GET("/exist/saga", existConfigSaga)
			//public saga config
			config.POST("/tag/update", publicSagaConfig)
			//update and public saga config
			config.POST("/update/now/saga", releaseSagaConfig)
			//get current saga config
			config.GET("/option/saga", optionSaga)

		}

		// V1 wechat will carry cookie
		wechat := version.Group("/wechat")
		{
			wechat.GET("", queryContacts)
			contactLog := wechat.Group("/log")
			{
				contactLog.GET("/query", queryContactLogs)
			}
			redisdata := version.Group("/redisdata")
			{
				redisdata.GET("/query", queryRedisdata)
			}

			wechat.GET("/analysis/contacts", syncWechatContacts)
			wechat.POST("/appchat/create", createWechat)
			wechat.GET("/appchat/create/log", queryWechatCreateLog)
			wechat.GET("/appchat/get", getWechat)
			wechat.POST("/appchat/send", sendGroupWechat)
			wechat.POST("/message/send", sendWechat)
			wechat.POST("/appchat/update", updateWechat)
		}
	}

	version1 := e.Group("/ep/admin/saga/v2")
	{
		// V2 wechat will not carry cookie
		wechat := version1.Group("/wechat")
		{
			wechat.POST("/appchat/create", createWechat)
			wechat.GET("/appchat/create/log", queryWechatCreateLog)
			wechat.GET("/appchat/get", getWechat)
			wechat.POST("/appchat/send", sendGroupWechat)
			wechat.POST("/message/send", sendWechat)
			wechat.POST("/appchat/update", updateWechat)
		}
	}
}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
