package http

import (
	"go-common/app/interface/main/reply/conf"
	"go-common/app/interface/main/reply/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	cnf       *conf.Config
	rpSvr     *service.Service
	authSvc   *auth.Auth
	verifySvc *verify.Verify
)

// Init init http
func Init(c *conf.Config) {
	initService(c)
	engine := bm.DefaultServer(c.BM)
	outerRouter(engine)
	interRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("xhttp.Serve error(%v)", err)
		panic(err)
	}
}

func initService(c *conf.Config) {
	cnf = c
	authSvc = auth.New(c.Auth)
	verifySvc = verify.New(c.Verify)
	rpSvr = service.New(conf.Conf)
}

func outerRouter(e *bm.Engine) {
	// init api
	e.GET("/monitor/ping", ping)
	// reply
	group := e.Group("/x/v2/reply")
	{
		group.GET("", authSvc.Guest, reply)
		group.GET("/cursor", authSvc.Guest, replyByCursor)
		group.GET("/reply/cursor", authSvc.Guest, subReplyByCursor)
		group.GET("/hot", replyHots)
		group.GET("/emojis", emojis)
		group.GET("/web/emojis", emojis)
		group.GET("/info", authSvc.Guest, replyInfo)
		group.GET("/minfo", authSvc.Guest, replyMultiInfo)
		group.GET("/reply", authSvc.Guest, replyReply)
		group.GET("/jump", authSvc.Guest, jumpReply)
		group.GET("/count", authSvc.Guest, replyCount)
		group.GET("/mcount", authSvc.Guest, replyMultiCount)
		group.GET("/log", authSvc.Guest, replyAdminLog)
		group.POST("/add", authSvc.User, addReply)
		group.POST("/action", authSvc.User, likeReply)
		group.POST("/hate", authSvc.User, hateReply)
		group.POST("/show", authSvc.User, showReply)
		group.POST("/hide", authSvc.User, hideReply)
		group.POST("/del", authSvc.User, delReply)
		group.POST("/top", authSvc.User, AddTopReply)
		group.POST("/report", authSvc.User, reportReply)
		group.GET("/topics", authSvc.Guest, getTopics)
		group.GET("/report/related", authSvc.Guest, reportRelated)
		group.GET("/report/reply", authSvc.Guest, reportSndReply)
		group.GET("/dialog", authSvc.Guest, dialog)
		group.GET("/dialog/cursor", authSvc.Guest, dialogByCursor)
		// 5.37需求
		group.GET("/main", authSvc.Guest, xreply)
		group.GET("/folded", authSvc.Guest, subFolder)
		group.GET("/reply/folded", authSvc.Guest, rootFolder)
	}
}

func interRouter(e *bm.Engine) {
	// internal admin
	group := e.Group("/x/internal/v2/reply")
	{
		group.GET("/subject", verifySvc.Verify, adminSubject)
		group.POST("/subject/mid", verifySvc.Verify, adminSubjectMid)
		group.GET("/hot", verifySvc.Verify, replyHots)
		group.POST("/subject/state", verifySvc.Verify, adminSubjectState)
		group.POST("/subject/regist", verifySvc.Verify, adminSubRegist)
		group.POST("/audit", verifySvc.Verify, adminAuditSub)
		group.POST("/pass", verifySvc.Verify, adminPassReply)
		group.POST("/recover", verifySvc.Verify, adminRecoverReply)
		group.POST("/edit", verifySvc.Verify, adminEditReply)
		group.POST("/del", verifySvc.Verify, adminDelReply)
		group.POST("/top", verifySvc.Verify, adminAddTopReply)
		group.POST("/report/del", verifySvc.Verify, adminDelReplyByReport)
		group.POST("/report/ignore", verifySvc.Verify, adminIgnoreReport)
		group.POST("/report/recover", verifySvc.Verify, adminReportRecover)
		group.POST("/report/transfer", verifySvc.Verify, adminTransferReport)
		group.POST("/report/state", verifySvc.Verify, adminReportStateSet)
		group.GET("/info", verifySvc.Verify, replyInfo)
		group.GET("/count", verifySvc.Verify, replyCount)
		group.GET("/counts", verifySvc.Verify, replyCounts)
		group.GET("/minfo", verifySvc.Verify, replyMultiInfo)
		group.GET("/mcount", verifySvc.Verify, replyMultiCount)
		group.GET("/record", verifySvc.Verify, replyRecord)
		group.GET("/hots", verifySvc.Verify, hotsBatch)
		group.GET("/ishot", isHotReply)
	}
}
