package http

import (
	"net/http"

	"go-common/app/service/main/workflow/conf"
	"go-common/app/service/main/workflow/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	verifySrv *verify.Verify
	wkfSvc    *service.Service
)

// Init http server
func Init(c *conf.Config, svc *service.Service) {
	wkfSvc = svc
	verifySrv = verify.New(c.Verify)
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start() error(%v)", err)
		panic(err)
	}
}

// innerRouter
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	wkfG := e.Group("/x/internal/workflow")
	{
		wkfG.POST("/add", verifySrv.Verify, addChallenge)
		wkfG.GET("/close", verifySrv.Verify, closeChallenge)
		wkfG.GET("/info", verifySrv.Verify, challengeInfo)
		wkfG.GET("/list", verifySrv.Verify, listChallenge)
		wkfG.POST("/update/state", verifySrv.Verify, upChallengeState)
		wkfG.GET("/extra/info", verifySrv.Verify, businessExtra)
		wkfG.POST("/extra/up", verifySrv.Verify, upBusinessExtra)
		wkfG.POST("/reply/add", verifySrv.Verify, replyAddChallenge)
		wkfG.GET("/untreated", verifySrv.Verify, untreatedChallenge)
		appealG := wkfG.Group("/appeal")
		{
			appealG.POST("/add", verifySrv.Verify, addChallenge)
			appealG.GET("/close", verifySrv.Verify, closeChallenge)
			appealG.GET("/info", verifySrv.Verify, challengeInfo)
			appealG.GET("/list", verifySrv.Verify, listChallenge)
			appealG.POST("/state", verifySrv.Verify, upChallengeState)
			appealG.GET("/extra/info", verifySrv.Verify, businessExtra)
			appealG.POST("/extra/up", verifySrv.Verify, upBusinessExtra)
			appealG.POST("/reply/add", verifySrv.Verify, replyAddChallenge)
		}
		appealG3 := wkfG.Group("/appeal/v3")
		{
			appealG3.POST("/add", verifySrv.Verify, addChallenge3)
			appealG3.GET("/list", verifySrv.Verify, listChallenge3)
			appealG3.GET("/state", verifySrv.Verify, groupState3)
			appealG3.POST("/delete", verifySrv.Verify, deleteGroup)
			appealG3.POST("/public/referee", verifySrv.Verify, pubRefereeGroup)
		}
		tagG3 := wkfG.Group("/tag/v3", verifySrv.Verify)
		{
			tagG3.GET("/list", tagList3)
		}
		tagG := wkfG.Group("/tag", verifySrv.Verify)
		{
			tagG.GET("/list", tagList)
		}
		sobotG := wkfG.Group("/sobot")
		{
			sobotG.GET("/user", sobotSign(sobotFetchUser))
			sobotG.GET("/ticket/info", verifySrv.Verify, sobotInfoTicket)
			sobotG.POST("/ticket/add", verifySrv.Verify, sobotAddTicket)
			sobotG.POST("/ticket/modify", verifySrv.Verify, sobotModifyTicket)
			sobotG.POST("/reply/add", verifySrv.Verify, sobotAddReply)
		}
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := wkfSvc.Ping(c); err != nil {
		log.Error("workflow-service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
