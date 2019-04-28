package http

import (
	"net/http"

	"go-common/app/admin/ep/marthe/conf"
	"go-common/app/admin/ep/marthe/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	srv     *service.Service
	authSvc *permit.Permit
)

// Init Init.
func Init(c *conf.Config, s *service.Service) {
	srv = s
	authSvc = permit.New(c.Auth)
	engine := bm.DefaultServer(c.BM)
	engine.Ping(ping)
	outerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func outerRouter(e *bm.Engine) {
	base := e.Group("/ep/admin/marthe/v1", authSvc.Permit(""))
	{
		base.GET("/version", getVersion)

		bugly := base.Group("/bugly")
		{

			bugly.GET("/test", test)

			project := bugly.Group("/project") //done test
			{
				project.GET("", queryProject)                         //done test
				project.POST("/update", accessToBugly, updateProject) //done test
				project.POST("/query", queryProjects)                 //done test
				project.GET("/all", queryAllProjects)                 //done test

				project.GET("/versions", queryProjectVersions) // done test
			}

			version := bugly.Group("/version") //done test
			{
				version.GET("/run", accessToBugly, runVersions) //done test

				version.POST("/update", accessToBugly, updateVersion) //done test
				version.POST("/query", queryVersions)                 //done test
				version.GET("/list", getVersionAndProjectList)        // done test
				version.POST("/batch/query", queryBatchRun)
			}

			cookie := bugly.Group("/cookie")
			{
				cookie.POST("/add", updateCookie)                // done test
				cookie.POST("/update", updateCookie)             // done test
				cookie.POST("/query", queryCookies)              // done test
				cookie.GET("/status/update", updateCookieStatus) //done test
			}

			issue := bugly.Group("/issue")
			{
				issue.POST("/query", queryBuglyIssue) // done test
			}
		}

		tapd := base.Group("/tapd")
		{
			bug := tapd.Group("/bug")
			{
				bug.GET("/project/insert", bugInsertTapdWithProject) // done test access
				bug.GET("/version/insert", bugInsertTapdWithVersion) // done test access
				bug.POST("/filtersql/check", checkFilterSql)         // done test
				bug.POST("/record/query", queryBugRecord)            // done test

				bug.POST("/template/update", updateTapdBugTpl) // done test access
				bug.POST("/template/query", queryTapdBugTpl)   // done test
				bug.GET("/template/all", queryAllTapdBugTpl)   // done test

				bug.POST("/version/template/update", updateTapdBugVersionTpl) //done test access
				bug.POST("/version/template/query", queryTapdBugVersionTpl)   //done test

				bug.POST("/conf/priority/update", updateTapdBugPriorityConf) //done test
				bug.POST("/conf/priority/query", queryTapdBugPriorityConf)   //done test

				bug.GET("/conf/auth/check", checkAuth)

			}
		}

		user := base.Group("/user")
		{
			user.GET("/query", queryUserInfo)                   // done test
			user.GET("/wechat/contact/sync", syncWechatContact) // done test
			user.GET("/visible/bugly", updateVisibleBugly)      // done test
		}
	}
}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func getVersion(c *bm.Context) {
	v := new(struct {
		Version string `json:"version"`
	})
	v.Version = "v.0.0.0.3"
	log.Info("marthe current version [%s]", v.Version)
	c.JSON(v, nil)
}
