package http

import (
	"net/http"

	"go-common/app/service/main/coupon/conf"
	"go-common/app/service/main/coupon/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	verifySvc *verify.Verify
	svc       *service.Service
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	svc = s
	initService(c)
	engine := bm.DefaultServer(c.BM)
	interRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	verifySvc = verify.New(nil)
}

// interRouter init inner router api path.
func interRouter(e *bm.Engine) {
	//init api
	e.Ping(ping)
	e.Register(register)
	group := e.Group("/x/internal/coupon", verifySvc.Verify)
	{
		group.GET("/count", userCoupon)
		group.POST("/use", useCoupon)
		group.GET("/info", couponInfo)
		group.POST("/add", addCoupon)
		group.POST("/change", changeCoupon)
		group.POST("/cartoon/use", useCartoonCoupon)
		group.POST("/grant", salaryCoupon)

		ae := group.Group("/allowance")
		ae.POST("/use", useAllowance)
		ae.POST("/notify", useNotify)
		ae.GET("/count", allowanceCount)
		ae.POST("/receive", receiveAllowance)
		// 元旦活动
		group.GET("/prize/cards", prizeCards)
		group.POST("/prize/draw", prizeDraw)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("coupon http service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// register check server ok.
func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
