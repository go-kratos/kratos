package http

import (
	"net/http"

	"go-common/app/admin/main/up-rating/conf"
	"go-common/app/admin/main/up-rating/service"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	svr *service.Service
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	svr = s
	engine := bm.DefaultServer(c.BM)
	initRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func initRouter(r *bm.Engine) {
	r.Ping(ping)
	rating := r.Group("/allowance/api/x/admin/rating")
	statis := rating.Group("/statis")
	{
		statis.GET("/graph", statisGraph)
		statis.GET("/list", statisList)
		statis.GET("/export", statisExport)
	}
	score := rating.Group("/score")
	{
		score.GET("/list", scoreList)
		score.GET("/export", scoreExport)
		score.GET("/up/current", scoreCurrent)
		score.GET("/up/history", scoreHistory)
	}
	param := rating.Group("/param")
	{
		param.POST("/insert", paramInsert)
	}
	trend := rating.Group("trend")
	{
		trend.GET("/ascend/list", ascList)
		trend.GET("/descend/list", descList)
	}
	au := rating.Group("authority")
	{
		au.POST("/add", addAuthority)
		au.POST("/remove", removeAuthority)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	var err error
	if err = svr.Ping(c); err != nil {
		log.Error("up-rating-admin ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
