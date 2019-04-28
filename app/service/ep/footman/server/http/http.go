package http

import (
	"net/http"

	"go-common/app/service/ep/footman/conf"
	"go-common/app/service/ep/footman/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	srv *service.Service
)

// Init Init.
func Init(c *conf.Config, s *service.Service) {
	srv = s
	engine := bm.DefaultServer(c.BM)
	engine.Ping(ping)
	outerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func outerRouter(e *bm.Engine) {
	version := e.Group("/ep/admin/footman/v1")
	{
		version.GET("/version", getVersion)

		bugly := version.Group("/bugly")
		{
			bugly.GET("/issue/query", queryIssue)
			bugly.GET("/issue/save", saveIssue)
			bugly.GET("/saveissue", saveIssue)
			bugly.GET("/saveissues", saveIssues)
			bugly.GET("/updatetoken", updateToken)
		}

		tapd := version.Group("/tapd")
		{
			tapd.GET("/file/save", saveFiles)
			tapd.GET("/file/story/download", downloadStoryFile)
			tapd.GET("/file/change/download", downloadChangeFile)
			tapd.GET("/file/iteration/download", downloadIterationFile)
			tapd.GET("/file/bug/download", downloadBugFile)
		}

		bugly2tapd := version.Group("/bugly2tapd")
		{
			bugly2tapd.GET("/save", saveBugly2Tapd)
			bugly2tapd.GET("/status/update", updateBuglyStatusInTapd)
			bugly2tapd.GET("/title/update", updateTitleInTapd)
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
	v.Version = "v.1.5.8.2"
	c.JSON(v, nil)

}
