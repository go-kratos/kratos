package http

import (
	"net/http"

	"go-common/app/admin/main/laser/conf"
	"go-common/app/admin/main/laser/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	svc     *service.Service
	authSrc *permit.Permit
)

// Init http server
func Init(c *conf.Config) {
	svc = service.New(c)
	authSrc = permit.New(c.Auth)
	engine := bm.DefaultServer(c.BM)
	initRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	laser := e.Group("/x/admin/laser", authSrc.Verify())
	{
		task := laser.Group("/task")
		{
			task.POST("/add", addTask)
			task.GET("/list", queryTask)
			task.POST("/edit", updateTask)
			task.POST("/delete", deleteTask)
		}
		taskLog := laser.Group("/task_log")
		{
			taskLog.GET("/list", queryTaskLog)
		}
		recheck := laser.Group("/archive/stat")
		{
			recheck.GET("/panel", recheckPanel)
			recheck.GET("/user", recheckUser)
			recheck.GET("/123_recheck", recheck123)
		}
		cargo := laser.Group("/archive/cargo")
		{
			cargo.GET("/audit/csv", auditCargoCsv)
			cargo.GET("/auditors", auditorCargo)
		}
		tag := laser.Group("/archive/tag")
		{
			tag.GET("/recheck", tagRecheck)
		}
		video := laser.Group("/video/stat")
		{
			video.GET("/random_video", randomVideo)
			video.GET("/random_video/csv", csvRandomVideo)
			video.GET("/fixed_video", fixedVideo)
			video.GET("/fixed_video/csv", csvFixedVideo)
		}
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("laser-admin service ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
