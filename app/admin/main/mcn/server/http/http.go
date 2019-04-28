package http

import (
	"net/http"

	"go-common/app/admin/main/mcn/conf"
	"go-common/app/admin/main/mcn/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	srv     *service.Service
	authSvc *permit.Permit
)

// Init init
func Init(c *conf.Config) {
	initService(c)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

// initService init service
func initService(c *conf.Config) {
	srv = service.New(c)
	authSvc = permit.New(c.Auth)
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/allowance/api/x/admin/mcn") // authSvc.Verify() manager use
	{
		// mcn account .
		g.POST("/sign/upload", upload)
		g.POST("/sign/entry", mcnSignEntry)
		g.GET("/sign/list", mcnSignList)
		g.POST("/sign/op", mcnSignOP)
		g.GET("/sign/up/list", mcnUPReviewList)
		g.POST("/sign/up/op", mcnUPOP)
		g.POST("/sign/permit/op", mcnPermitOP)
		g.GET("/sign/up/permit/list", mcnUPPermitList)
		g.POST("/sign/up/permit/op", mcnUPPermitOP)

		// mcn list.
		g.GET("/list", mcnList)
		// g.POST("/pay/add", mcnPayAdd)
		g.POST("/pay/edit", mcnPayEdit)
		g.POST("/pay/state/edit", mcnPayStateEdit)
		g.POST("/state/edit", mcnStateEdit)
		g.POST("/renewal", mcnRenewal)

		// up list
		g.GET("/info", mcnInfo)
		g.GET("/up/list", mcnUPList)
		g.POST("/up/state/edit", mcnUPStatEdit)
		// 二期
		g.GET("/cheat/list", mcnCheatList)
		g.GET("/cheat/up/list", mcnCheatUPList)
		g.GET("/import/up/info", mcnImportUPInfo)
		g.POST("/import/up/reward/sign", mcnImportUPRewardSign)
		g.GET("/increase/list", mcnIncreaseList)

		// up fans rank
		g.GET("/rank/archive/likes", arcTopDataStatistics)

		// up fans analyze
		g.GET("/up/fans/analyze", mcnFansAnalyze)

		// mcn total statistics
		g.GET("/total/statistics", mcnsTotalDatas)

		// up recommend
		g.GET("/recommend/list", recommendList)
		g.POST("/recommend/op", recommendOP)
		g.POST("/recommend/add", recommendAdd)
	}
}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
