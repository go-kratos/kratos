package http

import (
	"net/http"

	"go-common/app/service/main/coin/conf"
	"go-common/app/service/main/coin/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	verifySrv *verify.Verify
	coinSvc   *service.Service
	antispamM *antispam.Antispam
)

// Init init http linstening.
func Init(c *conf.Config, s *service.Service) {
	coinSvc = s
	antispamM = antispam.New(c.Antispam)
	verifySrv = verify.New(c.Verify)
	// init outer router
	engine := bm.DefaultServer(c.BM)
	outerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router.
func outerRouter(r *bm.Engine) {
	r.Ping(ping)
	r.Register(register)
	cr := r.Group("/x/coin", verifySrv.Verify, Business)
	{
		cr.POST("/add", antispamM.ServeHTTP, addCoin)
		cr.POST("/settle", updateSettle)
		cr.GET("/v2/list", list)
		cr.GET("/today/exp", todayexp)
	}
	cr1 := r.Group("/x/internal/v1/coin", verifySrv.Verify, Business)
	{
		cr1.POST("/add", antispamM.ServeHTTP, internalAddCoin)
		cr1.GET("/list", list)
		cr1.GET("/coins", coins)
		cr1.GET("/item/coins", itemCoins)
		cr1.POST("/amend", amend)
		cr1.GET("/user/count", userCoins)
		cr1.GET("/user/log", coinLog)
		cr1.POST("/user/modify", modify)
		cr1.GET("/creation/counts", ccounts)
	}
}

func ping(c *bm.Context) {
	if err := coinSvc.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

// Business set business
func Business(c *bm.Context) {
	business := c.Request.Form.Get("business")
	if business == "" {
		return
	}
	tp, err := coinSvc.CheckBusiness(business)
	if err != nil {
		c.JSON(nil, err)
		c.Abort()
		return
	}
	c.Set("business", tp)
}
