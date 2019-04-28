package http

import (
	"net/http"

	"go-common/app/job/main/up-rating/conf"
	"go-common/app/job/main/up-rating/service"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	svr *service.Service
)

// Init init http router.
func Init(c *conf.Config) {
	svr = service.New(conf.Conf)
	// bm
	engine := bm.DefaultServer(c.BM)
	initRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	mr := e.Group("/x/internal/job/up-rating")
	mr.POST("/past/score", pastScore)
	mr.POST("/past/record", pastRecord)
	mr.POST("/score", score)
	mr.POST("/score/del", delScore)
	mr.POST("/statistics", statistics)
	mr.POST("/trend", trend)
	mr.POST("/trend/del", delTrends)
	mr.POST("/task/status", taskStatus)
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := svr.Ping(c); err != nil {
		log.Error("service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// Close http close
func Close() {
	svr.Close()
}
