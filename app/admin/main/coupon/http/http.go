package http

import (
	"net/http"

	"go-common/app/admin/main/coupon/conf"
	"go-common/app/admin/main/coupon/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	svc     *service.Service
	authSvc *permit.Permit
)

// Init init
func Init(c *conf.Config) {
	initService(c)
	// init router
	engine := bm.DefaultServer(c.BM)
	initRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	authSvc = permit.New(c.Auth)
	svc = service.New(c)
}

// initRouter init outer router api path.
func initRouter(e *bm.Engine) {
	//init api
	e.Ping(ping)
	group := e.Group("/x/admin/coupon", authSvc.Permit("CONPON"))
	group.GET("/batch/list", batchlist)
	group.POST("/batch/add", batchadd)
	group.GET("/appinfo/list", allAppInfo)
	group.POST("/salary", salaryCoupon)

	allowance := group.Group("/allowance", authSvc.Permit("CONPON_ALLOWANCE"))
	allowance.GET("/batch/list", batchlist)
	allowance.GET("/batch/info", batchInfo)
	allowance.POST("/batch/block", allowanceBatchBlock)
	allowance.POST("/batch/unblock", allowanceBatchUnBlock)
	allowance.POST("/batch/add", allowanceBatchadd)
	allowance.POST("/batch/modify", allowanceBatchModify)
	allowance.GET("/list", allowanceList)
	allowance.POST("/block", allowanceBlock)
	allowance.POST("/unblock", allowanceUnBlock)
	allowance.POST("/salary", allowanceSalary)
	allowance.POST("/activity/salary", batchSalaryCoupon)
	allowance.POST("/uploadfile", uploadFile)

	view := group.Group("/view", authSvc.Permit("COUPON_VIEW"))
	{
		view.GET("/batch/list", batchlist)
		view.GET("/batch/info", batchInfo)
		view.POST("/batch/block", batchBlock)
		view.POST("/batch/unblock", batchUnBlock)
		view.POST("/batch/add", viewBatchAdd)
		view.POST("/batch/save", viewBatchSave)
		view.GET("/list", viewList)
		view.POST("/block", viewBlock)
		view.POST("/unblock", viewUnblock)
		view.POST("/salary", salaryView)
	}

	// code
	group.GET("/batch_code/list", codeBatchList)
	group.POST("/batch_code/block", codeBatchBlock)
	group.POST("/batch_code/unblock", codeBatchUnBlock)
	group.POST("/batch_code/add", codeAddBatch)
	group.POST("/batch_code/modify", codeBatchModify)

	group.GET("/code/list", codePage)
	group.POST("/code/block", codeBlock)
	group.POST("/code/unblock", codeUnBlock)
	group.GET("/code/export", exportCode)
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("coupon http admin ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
