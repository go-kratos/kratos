package http

import (
	"net/http"

	"go-common/app/interface/main/history/conf"
	"go-common/app/interface/main/history/model"
	"go-common/app/interface/main/history/service"
	"go-common/library/log"
	"go-common/library/log/anticheat"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	cnf       *conf.Config
	hisSvc    *service.Service
	collector *anticheat.AntiCheat
	authSvc   *auth.Auth
	verifySvc *verify.Verify
)

// Init init http
func Init(c *conf.Config, s *service.Service) {
	cnf = c
	authSvc = auth.New(c.Auth)
	verifySvc = verify.New(c.Verify)
	if c.Collector != nil {
		collector = anticheat.New(c.Collector)
	}
	hisSvc = s
	engine := bm.DefaultServer(c.BM)
	outerRouter(engine)
	interRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("xhttp.Serve error(%v)", err)
		panic(err)
	}
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	e.GET("/monitor/ping", ping)
	group := e.Group("/x/v2/history", authSvc.User)
	{
		group.GET("", history)
		group.POST("/add", addHistory)
		group.POST("/del", delHistory)
		group.POST("/delete", delete)
		group.POST("/clear", clearHistory)
		group.POST("/report", report)   //just mobile use it.
		group.POST("/reports", reports) //just mobile use it.
	}
	toViewGroup := group.Group("/toview", authSvc.User)
	{
		toViewGroup.GET("", toView)
		toViewGroup.GET("/web", webToView)
		toViewGroup.GET("/remaining", remainingToView)
		toViewGroup.POST("/add", addToView)
		toViewGroup.POST("/adds", addMultiToView)
		toViewGroup.POST("/del", delToView)
		toViewGroup.POST("/clear", clearToView)
	}
	shadowGroup := group.Group("/shadow", authSvc.User)
	{
		shadowGroup.GET("", shadow)
		shadowGroup.POST("/set", setShadow)
	}
}

// interRouter init internal router api path using the outer router port
func interRouter(e *bm.Engine) {
	group := e.Group("/x/internal/v2/history")
	{
		group.GET("/manager", verifySvc.Verify, managerHistory)
		group.GET("/aids", verifySvc.VerifyUser, aids) // bangumi use it.
		group.POST("/delete", verifySvc.VerifyUser, delete)
		group.POST("/clear", verifySvc.VerifyUser, clearHistory)
		group.POST("/flush", verifySvc.Verify, flush)                // history-job use it.
		group.POST("/report", verifySvc.VerifyUser, innerReport)     //just mobile use it.
		group.GET("/position", verifySvc.VerifyUser, position)       //just mobile use it.
		group.GET("/resource", verifySvc.VerifyUser, resource)       //just mobile use it.
		group.GET("/resource/type", verifySvc.VerifyUser, resources) //just mobile use it.
	}
	toviewGroup := group.Group("/toview")
	{
		toviewGroup.GET("/manager", verifySvc.Verify, managerToView)
		toviewGroup.POST("/adds", verifySvc.VerifyUser, addMultiToView) // space use it.
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := hisSvc.Ping(c); err != nil {
		log.Error("history service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func business(c *bm.Context) (tp int8, ok bool) {
	ok = true
	business := c.Request.Form.Get("business")
	if business == "" {
		return
	}
	tp, err := model.CheckBusiness(business)
	if err != nil {
		ok = false
		c.JSON(nil, err)
		c.Abort()
	}
	return
}
