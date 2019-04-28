package http

import (
	"go-common/app/job/main/videoup-report/conf"
	"go-common/app/job/main/videoup-report/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"net/http"
)

var (
	vdaSvc *service.Service
)

// Init http server
func Init(c *conf.Config, s *service.Service) {
	vdaSvc = s
	// init internal router
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	// init internal server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// innerRouter
func innerRouter(e *bm.Engine) {
	e.GET("/monitor/ping", ping)

	g := e.Group("/x/videoup/report")
	{
		task := g.Group("/task")
		{
			task.GET("/tooks", taskTooks)
		}
		video := g.Group("/video")
		{
			video.GET("/audit", videoAudit)
			video.GET("/xcode", videoXcode)
		}
		archive := g.Group("/archive")
		{
			archive.GET("/movetype", moveType)
			archive.GET("/roundflow", roundFlow)
		}
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	var err error
	if err = vdaSvc.Ping(c); err != nil {
		log.Error("videoup-report-job ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
