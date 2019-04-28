package http

import (
	newbieService "go-common/app/interface/main/growup/service/newbie"
	"net/http"

	"go-common/app/interface/main/growup/conf"
	"go-common/app/interface/main/growup/service"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
)

var (
	svc       *service.Service
	authSvr   *auth.Auth
	newbieSvr *newbieService.Service
)

// Init http server
func Init(c *conf.Config) {
	initService(c)

	engine := bm.DefaultServer(c.BM)
	setupInnerEngine(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func initService(c *conf.Config) {
	svc = service.New(c)
	authSvr = auth.New(nil)
	newbieSvr = newbieService.New(c)
}

func setupInnerEngine(e *bm.Engine) {
	e.Ping(ping)

	allowance := e.Group("/allowance/api/x/internal/growup")
	allowanceUp := allowance.Group("/up")
	{
		//		allowanceUp.GET("/status", getUpStatus)
		allowanceUp.POST("/add", join)
		allowanceUp.POST("/quit", quit)
		allowanceUp.GET("/withdraw", getWithdraw)
		allowanceUp.POST("/withdraw/success", withdrawSuccess)
	}

	studio := e.Group("/studio/growup/web", authSvr.User)
	studioUp := studio.Group("/up")
	{
		studioUp.GET("/income/stat", upIncomeStat)
		studioUp.GET("/summary", upSummary)
		studioUp.GET("/archive/summary", archiveSummary)
		studioUp.GET("/charge", upCharge)
		studioUp.GET("/archive/income", archiveIncome)
		studioUp.GET("/archive/detail", archiveDetail)
		studioUp.GET("/archive/breach", archiveBreach)
		studioUp.GET("/withdraw/detail", withdrawDetail)
		studioUp.POST("/quit", quit1)
		studioUp.GET("/status", getUpStatus)
		studioUp.POST("/av/join", joinAv)
		studioUp.POST("/bgm/join", joinBgm)
		studioUp.POST("/column/join", joinColumn)
		studioUp.GET("/bill", upBill)
		studioUp.GET("/year", upYear)
		// 新手信
		studioUp.GET("/newbie/letter", upNewbieLetter)

		exchange := studioUp.Group("/exchange")
		{
			exchange.GET("/state", goodsState)
			exchange.GET("/show", goodsShow)
			exchange.GET("/record", goodsRecord)
			exchange.POST("/buy", goodsBuy)
		}

	}
	studioActivity := studio.Group("/activity")
	{
		studioActivity.GET("/show", showActivity)
		studioActivity.POST("/sign_up", signUpActivity)
	}

	specialAward := e.Group("/studio/growup/web/special/award", authSvr.Guest)
	{
		specialAward.GET("/info", sepcialAwardInfo)
		specialAward.GET("/detail", specialAwardDetail)
		specialAward.GET("/list", listSpecialAward)
		specialAward.GET("/winner", specialAwardWinners)
	}

	specialAwardUser := e.Group("/studio/growup/web/special/award", authSvr.User)
	{
		specialAwardUser.GET("/record", specialAwardRecord)
		specialAwardUser.GET("/record/poster", specialAwardPoster)
		specialAwardUser.GET("/up/status", specialAwardUpStatus)
		specialAwardUser.POST("/join", joinSpecialAward)
	}

	studio.GET("/notice/latest", latestNotice)
	studio.GET("/notices", notices)

	studio.GET("/banner", banner)
}

func ping(c *bm.Context) {
	var err error
	if err = svc.Ping(c); err != nil {
		log.Error("service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
