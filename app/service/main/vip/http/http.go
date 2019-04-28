package http

import (
	"net/http"

	"go-common/app/service/main/vip/conf"
	"go-common/app/service/main/vip/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	vipSvc *service.Service
	vrfSvr *verify.Verify
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	vrfSvr = verify.New(nil)
	vipSvc = service.New(c)
	// init router
	engineOuter := bm.DefaultServer(c.BM)
	innerRouter(engineOuter)
	if err := engineOuter.Start(); err != nil {
		log.Error("engineOuter.Start() error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)

	//internal api
	big := e.Group("/x/internal/big", vrfSvr.Verify)
	{
		big.GET("/batchInfo", batchInfo)
		big.POST("/useBatchInfo", useBatchInfo)
	}

	vip := e.Group("/x/internal/vip", vrfSvr.Verify)
	{
		// bcoin
		vip.GET("/bcoin/list", bpList)
		// point
		vip.POST("/point/exchange_vip", buyVipWithPoint)
		vip.POST("/point/rule", rule)
		// user
		vip.GET("/user/info", byMid)
		vip.GET("/user/list", vipInfos)
		vip.GET("/user/history", vipHistory)
		vip.GET("/user/history/h5", vipH5History)
		vip.GET("/user/infobo", vipInfo) // for old service.
		// order
		vip.GET("/order/status", status)
		vip.GET("/order/list", orders)
		vip.POST("/order/create", createOrder)
		vip.POST("/order/oldcreate", createOldOrder) // for old service
		vip.GET("/order/mng", orderMng)
		vip.GET("/order/rescision", rescision)

		//panel
		vip.GET("/panel", pannelInfoNew)

		// panel
		vip.GET("/panel/single/info", vipUserMonthPanel)
		vip.GET("/panel/pirce", vipPirce)

		// price
		vip.GET("/price/by_product_id", priceceByProductID)
		vip.GET("/price/by_id", priceceByID)

		// code
		vip.GET("/code/verify", webToken)
		vip.POST("/code/open", openCode)
		vip.GET("/code/info", codeInfo)
		vip.GET("/code/infos", codeInfos)
		vip.POST("/code/belong", belong)
		vip.POST("/active/infos", actives)
		vip.GET("/code/opened", codeOpened)

		// tips
		vip.GET("/tips", tips)
		//coupon
		vip.POST("/coupon/cancel", cancelUseCoupon)
		vip.GET("/coupon/info", allowanceInfo)

		// FIXME: sync user
		vip.POST("/sync/user", syncUser)

		vip.POST("/order/create/qr", createQrCodeOrder)

		//act
		vip.POST("/activity/prize/grant", thirdPrizeGrant)
		vip.POST("/ele/vip/grant", grantAssociateVip)
	}

	vip2 := e.Group("/x/internal/vip/v2", vrfSvr.Verify)
	{
		vip2.POST("/order/create", createOrder2)
	}

	vipNotSign := e.Group("/x/internal/vip")
	{
		// notify
		vipNotSign.GET("/notify", notify)
		vipNotSign.GET("/notify/v2", notify2)
		vipNotSign.GET("/notify/sign", signNotify)
		vipNotSign.GET("/notify/refund", refundOrderNotify)
	}

}

// ping check server ok.
func ping(c *bm.Context) {
	var err error
	if err = vipSvc.Ping(c); err != nil {
		log.Error("service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// register check server ok.
func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
