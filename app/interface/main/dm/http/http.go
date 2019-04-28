package http

import (
	"net/http"

	"go-common/app/interface/main/dm/conf"
	"go-common/app/interface/main/dm/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	// service
	dmSvc       *service.Service
	antispamSvc *antispam.Antispam
	authSvc     *auth.Auth
	verifySvc   *verify.Verify
)

// Init http init
func Init(c *conf.Config, s *service.Service) {
	dmSvc = s
	authSvc = auth.New(c.Auth)
	verifySvc = verify.New(c.Verify)
	antispamSvc = antispam.New(c.Antispam)
	engine := bm.DefaultServer(c.HTTPServer)
	outerRouter(engine)
	interRouter(engine)
	// init external server
	if err := engine.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	group := e.Group("/x/dm")
	{
		fltGroup := group.Group("/filter")
		{
			fltGroup.GET("/user", authSvc.User, userRules)
			fltGroup.POST("/user/add", authSvc.User, antispamSvc.Handler(), addUserRule)
			fltGroup.POST("/user/add2", authSvc.User, antispamSvc.Handler(), multiAddUserRule)
			fltGroup.POST("/user/del", authSvc.User, delUserRules)
			fltGroup.GET("/global", authSvc.Guest, globalRuleEmpty)
		}
		group.POST("/protect/apply", authSvc.User, addPa)
		group.POST("/recall", authSvc.User, recall)
		group.GET("/user", authSvc.User, midHash)
		group.GET("/transfer/list", authSvc.User, transferList)
		group.POST("/transfer/retry", authSvc.User, transferRetry)
		group.POST("/adv/buy", authSvc.User, buyAdv)
		group.GET("/adv/state", authSvc.User, advState)
		group.GET("/up/banned/users", authSvc.User, assistBannedUsers)
		group.POST("/up/banned/del", authSvc.User, AssistDelBanned2)
		group.POST("/assist/banned", authSvc.User, assistBanned)
		group.POST("/assist/del", authSvc.User, assistDelete)
		group.POST("/report/add", authSvc.User, addReport)
		group.POST("/report/add2", authSvc.User, addReport2)
	}
}

func interRouter(e *bm.Engine) {
	group := e.Group("/x/internal/dm")
	{
		advGroup := group.Group("/adv")
		{
			advGroup.POST("/pass", verifySvc.VerifyUser, passAdv)
			advGroup.POST("/deny", verifySvc.VerifyUser, denyAdv)
			advGroup.POST("/cancel", verifySvc.VerifyUser, cancelAdv)
			advGroup.GET("/list", verifySvc.VerifyUser, advList)
		}
		fltGroup := group.Group("/filter")
		{
			fltGroup.POST("/global/add", verifySvc.Verify, addGlobalRule)
			fltGroup.POST("/global/del", verifySvc.Verify, delGlobalRules)
			fltGroup.GET("/index/list", verifySvc.VerifyUser, filterList)
			fltGroup.POST("/index/edit", verifySvc.VerifyUser, editFilter)
		}
		rptGroup := group.Group("/report")
		{
			rptGroup.POST("/up/edit", verifySvc.VerifyUser, editReport)
			rptGroup.GET("/up/list", verifySvc.VerifyUser, reportList)
			rptGroup.GET("/up/archives", verifySvc.VerifyUser, rptArchives)
		}
		prtGroup := group.Group("/up/protect/apply")
		{
			prtGroup.POST("/notice/switch", verifySvc.Verify, uptPaSwitch)
			prtGroup.POST("/status", verifySvc.Verify, UptPaStatus)
			prtGroup.GET("/list", verifySvc.Verify, paLs)
			prtGroup.GET("/video/list", verifySvc.Verify, paVideoLs)
		}
		group.POST("/assist/banned/upt", verifySvc.VerifyUser, assistBannedUpt)
		group.POST("/up/transfer", verifySvc.VerifyUser, transfer)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := dmSvc.Ping(c); err != nil {
		log.Error("dm service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
