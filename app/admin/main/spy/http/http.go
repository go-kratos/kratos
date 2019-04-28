package http

import (
	"net/http"

	"go-common/app/admin/main/spy/conf"
	"go-common/app/admin/main/spy/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	spySrv *service.Service
	vfySvc *verify.Verify
)

// Init init http sever instance.
func Init(c *conf.Config) {
	initService(c)
	// init inner router
	engine := bm.DefaultServer(c.HTTPServer)
	innerRouter(engine)
	// init inner server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	// health check
	e.Ping(ping)
	spy := e.Group("/x/admin/spy", vfySvc.Verify)

	score := spy.Group("/score")
	score.GET("/query", userInfo)
	score.GET("/record", historyPage)
	score.POST("/base/reset", resetBase)
	score.POST("/base/refresh", refreshBase)
	score.POST("/event/reset", resetEvent)

	spy.POST("/stat/clear", clearCount)

	factor := spy.Group("/factor")
	factor.GET("/list", factors)
	factor.POST("/add", addFactor)
	factor.POST("/modify", updateFactor)

	spy.POST("/event/add", addEvent)
	spy.POST("/event/modify", updateEventName)
	spy.POST("/service/add", addService)
	spy.POST("/group/add", addGroup)

	setting := spy.Group("/setting")
	setting.GET("/list", settingList)
	setting.POST("/add", updateSetting)

	sin := spy.Group("/sin")
	sin.POST("/state/update", updateStatState)
	sin.POST("/quantity/update", updateStatQuantity)
	sin.POST("/delete", deleteStat)
	sin.POST("/remark/add", addRemark)
	sin.GET("/remark/list", remarkList)
	sin.GET("/page", statPage)

	spy.GET("/report", report)

	spyWithoutAuth := e.Group("/x/admin/spy/internal")
	spyWithoutAuth.POST("/event/add", addEvent)
	spyWithoutAuth.POST("/event/modify", updateEventName)

	scoreWithoutAuth := spyWithoutAuth.Group("/score")
	scoreWithoutAuth.GET("/query", userInfo)
	scoreWithoutAuth.GET("/record", historyPage)
	scoreWithoutAuth.POST("/base/reset", resetBase)
	scoreWithoutAuth.POST("/base/refresh", refreshBase)
	scoreWithoutAuth.POST("/event/reset", resetEvent)

	factorWithoutAuth := spyWithoutAuth.Group("/factor")
	factorWithoutAuth.GET("/list", factors)
	factorWithoutAuth.POST("/add", addFactor)
	factorWithoutAuth.POST("/modify", updateFactor)
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := spySrv.Ping(c); err != nil {
		log.Error("spy admin ping error(%v)", err)
		c.JSON(nil, err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func initService(c *conf.Config) {
	spySrv = service.New(c)
	vfySvc = verify.New(nil)
}
