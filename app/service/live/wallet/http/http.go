package http

import (
	"net/http"

	"go-common/app/service/live/wallet/conf"
	"go-common/app/service/live/wallet/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	// depend service
	idfSvc    *verify.Verify
	walletSvr *service.Service
)

func Init(c *conf.Config, s *service.Service) {
	// init service
	walletSvr = s
	// init external router
	idfSvc = verify.New(nil)
	engineOut := bm.DefaultServer(c.BM.Inner)
	innerRouter(engineOut)
	// init Inner serve
	if err := engineOut.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
	engineLocal := bm.DefaultServer(c.BM.Local)
	localRouter(engineLocal)
	// init Local serve
	if err := engineLocal.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

// outerRouter init outer router api path.
func innerRouter(e *bm.Engine) {
	// init api
	e.Ping(ping)
	group := e.Group("/x/internal/livewallet", idfSvc.Verify)
	{
		group.GET("/wallet/get", get)
		group.GET("/wallet/delCache", delCache)
		group.GET("/wallet/getAll", getAll)
		group.POST("/wallet/getTid", getTid)
		group.POST("/wallet/recharge", recharge)
		group.GET("/wallet/query", query)
		group.POST("/wallet/pay", pay)
		group.POST("/wallet/exchange", exchange)
		group.POST("/wallet/modify", modify)
		group.POST("/flowwater/recordCoinStream", recordCoinStream)
	}
}

// localRouter init local router.
func localRouter(e *bm.Engine) {
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := walletSvr.Ping(c); err != nil {
		log.Error("live-userwallet http service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
