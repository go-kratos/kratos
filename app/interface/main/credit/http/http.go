package http

import (
	"net/http"

	"go-common/app/interface/main/credit/conf"
	"go-common/app/interface/main/credit/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	creditSvc *service.Service
	verifySvc *verify.Verify
	authSvc   *auth.Auth
	antSvc    *antispam.Antispam
)

// Init init http sever instance.
func Init(c *conf.Config) {
	initService(c)
	engineInner := bm.DefaultServer((c.BM.Inner))
	// init inner router
	innerRouter(engineInner)
	if err := engineInner.Start(); err != nil {
		log.Error("engineInner.Start() error(%v)", err)
		panic(err)
	}
}

// initService init service
func initService(c *conf.Config) {
	verifySvc = verify.New(c.Verify)
	authSvc = auth.New(c.AuthN)
	creditSvc = service.New(c)
	antSvc = antispam.New(c.Antispam)
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	// ping monitor
	e.Ping(ping)
	// internal api
	ig := e.Group("/x/internal/credit", verifySvc.Verify)
	{
		ig.POST("/labour/addQs", addQs)
		ig.POST("/labour/setQs", setQs)
		ig.POST("/labour/delQs", delQs)
		ig.GET("/labour/isanswered", isAnswered)
		// blocked
		ig.POST("/blocked/case/add", addBlockedCase)
		ig.POST("/blocked/info/add", addBlockedInfo)
		ig.POST("/blocked/info/batch/add", addBatchBlockedInfo)
		ig.GET("/blocked/user/num", blockedNumUser)
		ig.GET("/blocked/infos", batchBLKInfos)
		ig.GET("/blocked/cases", batchBLKCases)
		ig.GET("/blocked/historys", blkHistorys)
		// jury
		ig.GET("/jury/infos", batchJuryInfos)
		// other
		ig.POST("/opinion/del", delOpinion)
		ig.GET("/publish/infos", batchPublishs)
	}
	og := e.Group("/x/credit", bm.CORS())
	{
		// jury
		og.GET("/jury/requirement", authSvc.User, requirement)
		og.POST("/jury/apply", authSvc.User, apply)
		og.GET("/jury/jury", authSvc.User, jury)
		og.POST("/jury/caseObtain", authSvc.User, caseObtain)
		og.POST("/jury/vote", authSvc.User, vote)
		og.GET("/jury/voteInfo", authSvc.User, voteInfo)
		og.GET("/jury/caseInfo", authSvc.Guest, caseInfo)
		og.GET("/jury/juryCase", authSvc.User, juryCase)
		og.GET("/jury/caseList", authSvc.User, caseList)
		og.GET("/jury/kpi", authSvc.User, kpiList)
		og.GET("/jury/vote/opinion", authSvc.Guest, voteOpinion)
		og.GET("/jury/case/opinion", authSvc.Guest, caseOpinion)
		og.GET("/jury/notice", notice)
		og.GET("/jury/reasonList", reasonList)
		og.POST("/jury/caseObtain/open", authSvc.User, caseObtainByID)
		og.GET("/jury/juryCase/open", authSvc.Guest, spJuryCase)
		// appeal
		og.POST("/jury/appeal/add", authSvc.User, addAppeal)
		og.GET("/jury/appeal/status", authSvc.User, appealStatus)
		// labour
		og.GET("/labour/getQs", authSvc.User, getQs)
		og.POST("/labour/commitQs", authSvc.User, antSvc.Handler(), commitQs)
		// blocked
		og.GET("/blocked/user", authSvc.User, blockedUserCard)
		og.GET("/blocked/user/list", authSvc.User, blockedUserList)
		og.GET("/blocked/info", blockedInfo)
		og.GET("/blocked/info/appeal", authSvc.User, blockedAppeal)
		og.GET("/blocked/list", blockedList)
		// announcement
		og.GET("/publish/info", announcementInfo)
		og.GET("/publish/list", announcementList)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	err := creditSvc.Ping(c)
	if err != nil {
		log.Error("credit interface ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
