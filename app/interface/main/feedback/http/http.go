package http

import (
	"go-common/app/interface/main/feedback/conf"
	"go-common/app/interface/main/feedback/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	verifySvc   *verify.Verify
	authSvc     *auth.Auth
	feedbackSvr *service.Service
)

// Init init http
func Init(c *conf.Config) error {
	initService(c)
	engineOut := bm.DefaultServer(c.BM.Outer)
	outerRouter(engineOut)
	internalRouter(engineOut)
	// init external server
	if err := engineOut.Start(); err != nil {
		log.Error("engineOut.Start() error(%v) | config(%v)", err, c)
		panic(err)
	}
	engineLocal := bm.DefaultServer(c.BM.Local)
	localRouter(engineLocal)
	// init local server
	if err := engineLocal.Start(); err != nil {
		log.Error("engineLocal.Start() error(%v) | config(%v)", err, c)
		panic(err)
	}
	return nil
}

// initService init services.
func initService(c *conf.Config) {
	verifySvc = verify.New(nil)
	authSvc = auth.New(nil)
	feedbackSvr = service.New(c)
}

func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	// init api
	feedback := e.Group("/x/feedback")
	feedback.POST("/add", addReply)
	feedback.GET("/reply", verifySvc.Verify, replys)
	feedback.GET("/tag", verifySvc.Verify, replyTag)
	feedback.POST("/upload", upload)
	feedback.POST("/uploadFile", uploadFile)
	feedback.POST("/playerCheck", playerCheck)
	feedback.GET("/h5/tag", replyTag)
	feedback.POST("/h5/add", addReplyH5)
}

func internalRouter(e *bm.Engine) {
	feedbackInner := e.Group("/x/internal/feedback")
	feedbackInner.POST("/ugc/add", addWebReply)
	feedbackInner.GET("/ugc/session", verifySvc.Verify, sessions)
	feedbackInner.POST("/ugc/session/close", verifySvc.Verify, sessionsClose)
	feedbackInner.GET("/ugc/tag", verifySvc.Verify, ugcTag)
	feedbackInner.GET("/ugc/reply", verifySvc.Verify, webReply)
	feedbackInner.POST("/upload", authSvc.User, upload)
}

func localRouter(e *bm.Engine) {
}
