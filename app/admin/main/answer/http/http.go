package http

import (
	"go-common/app/admin/main/answer/conf"
	"go-common/app/admin/main/answer/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	answerSvc *service.Service
	permitX   *permit.Permit
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	answerSvc = s
	permitX = permit.New(c.Auth)
	// init inner router
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	// init inner server
	if err := engine.Start(); err != nil {
		log.Error("xhttp.Serve inner error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	// health check
	e.Ping(ping)
	//internal api
	que := e.Group("/x/admin/answer/v3")
	{
		que.GET("/types", types)
		que.GET("/questions", permitX.Permit("ANSWER_LIST"), quesList)
		que.POST("/upload", permitX.Permit("ANSWER_UPLOAD"), uploadQsts)
		que.POST("/question/edit", permitX.Permit("ANSWER_EDIT"), queEdit)
		que.POST("/question/disable", permitX.Permit("ANSWER_DISABLE"), queDisable)
		que.GET("/question/history", queHistory)
		que.GET("/history", history)

		que.GET("/load/img", loadImg)
	}
}

// ping check server ok.
func ping(c *bm.Context) {}
