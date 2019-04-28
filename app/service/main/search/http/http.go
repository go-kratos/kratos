package http

import (
	"go-common/app/service/main/search/conf"
	"go-common/app/service/main/search/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	svr *service.Service
)

// Init init http
func Init(c *conf.Config, s *service.Service) error {
	svr = s
	// init internal router
	engine := bm.DefaultServer(c.HTTPServer)
	route(engine)
	// init internal server
	if err := engine.Start(); err != nil {
		log.Error("bm.Start error(%v)", err)
		return err
	}
	return nil
}

func route(e *bm.Engine) {
	e.Ping(ping)
	searchG := e.Group("/x/internal/search")
	{
		//search
		searchG.GET("/reply", replySearch)
		searchG.GET("/dmhistory", dmHistorySearch)
		searchG.GET("/dmhistory/test", dmHistorySearch)
		searchG.GET("/dm", dmSearch)
		searchG.GET("/dm/date", dmDate)
		searchG.GET("/pgc", pgcSearch)
		//update
		searchG.POST("/reply/update", replyUpdate)
		searchG.POST("/pgc/update", pgcUpdate)
		searchG.POST("/dm/update", dmUpdate)
	}
}

// ping check health
func ping(ctx *bm.Context) {
	if err := svr.Ping(ctx); err != nil {
		log.Error("ping error(%v)", err)
		ctx.Error = err
		ctx.AbortWithStatus(503)
	}
}
