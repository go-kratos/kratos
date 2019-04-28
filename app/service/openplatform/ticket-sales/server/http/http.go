package http

import (
	"net/http"

	"go-common/app/service/openplatform/ticket-sales/conf"
	"go-common/app/service/openplatform/ticket-sales/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	svc *service.Service
)

// Init init ticket service.
func Init(c *conf.Config, s *service.Service) {
	svc = s
	// init router
	engine := bm.DefaultServer(c.BM)
	initRouter(engine)
	outerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// initRouter init inner router.
func initRouter(e *bm.Engine) {
	e.Ping(ping)
	// 等依赖方修改后，再删除
	group := e.Group("/x/internal/ticket/sales")
	{
		group.POST("/distrib/syncorder", syncOrder)
		group.GET("/distrib/getorder", getOrder)
	}
	//ticket
	group = e.Group("/openplatform/internal/ticket/sales")
	{
		group.POST("/distrib/syncorder", syncOrder)
		group.GET("/distrib/getorder", getOrder)

		group.POST("/promo/order/check", checkCreatePromoOrder)
		group.POST("/promo/order/create", createPromoOrder)
		group.GET("/promo/order/pay", payNotify)
		group.GET("/promo/order/cancel", cancelOrder)
		group.GET("/promo/order/check/issue", checkIssue)
		group.POST("/promo/order/finish/issue", finishIssue)

		group.GET("/promo/get", getPromo)
		group.POST("promo/create", createPromo)
		group.POST("promo/operate", operatePromo)
		group.POST("promo/edit", editPromo)

		group.POST("/settle/compare", settleCompare)
		group.POST("/settle/repush", settleRepush)
	}
}

// outerRouter init inner router.
func outerRouter(e *bm.Engine) {
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("ticket http service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
