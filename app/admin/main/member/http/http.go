package http

import (
	"go-common/app/admin/main/member/conf"
	"go-common/app/admin/main/member/http/block"
	"go-common/app/admin/main/member/service"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	authSvc *permit.Permit
	vfySvc  *verify.Verify
	svc     *service.Service
)

// Init init
func Init(c *conf.Config) {
	initService(c)
	// init router
	engine := bm.DefaultServer(c.BM)
	initRouter(engine)
	block.Setup(svc.BlockImpl(), engine, authSvc)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	authSvc = permit.New(c.Auth)
	vfySvc = verify.New(nil)
	svc = service.New(c)
}

// outerRouter init outer router api path.
func initRouter(e *bm.Engine) {
	// init api
	e.Ping(ping)
	e.Register(register)
	mg := e.Group("/x/admin/member")
	{
		mg.GET("/list", authSvc.Permit("ACCOUNT_REVIEW"), members)
		mg.POST("/exp/set", authSvc.Permit("ACCOUNT_REVIEW_SET_EXP"), expSet)
		mg.POST("/moral/set", authSvc.Permit("ACCOUNT_REVIEW"), moralSet)
		mg.POST("/rank/set", authSvc.Permit("ACCOUNT_REVIEW_SET_RANK"), rankSet)
		mg.POST("/coin/set", authSvc.Permit("ACCOUNT_REVIEW"), coinSet)
		mg.POST("/addit/remark/set", authSvc.Permit("ACCOUNT_REVIEW"), additRemarkSet)
		mg.GET("/profile", authSvc.Permit("ACCOUNT_REVIEW"), memberProfile)
		mg.GET("/exp/log", authSvc.Permit("ACCOUNT_EXP_LOG"), expLog)
		mg.GET("/face/history", authSvc.Permit("ACCOUNT_FACE_HISTORY"), faceHistory)
		mg.POST("/batch/formal", authSvc.Permit("ACCOUNT_BATCH_FORMAL"), batchFormal)
		mg.GET("/moral/log", moralLog)
		mg.POST("/sign/del", delSign)
		mg.POST("/pub/exp/msg", authSvc.Permit("ACCOUNT_PUB_EXP_MSG"), pubExpMsg)

		//个人信息审核
		mg.GET("/base/review", authSvc.Permit("ACCOUNT_REVIEW_AUDIT"), baseReview)
		mg.POST("/face/clear", authSvc.Permit("ACCOUNT_REVIEW_AUDIT"), clearFace)
		mg.POST("/sign/clear", authSvc.Permit("ACCOUNT_REVIEW_AUDIT"), clearSign)
		mg.POST("/name/clear", authSvc.Permit("ACCOUNT_REVIEW_AUDIT"), clearName)

		og := mg.Group("/official")
		{
			og.GET("/list", authSvc.Permit("OFFICIAL_VIEW"), officials)
			og.GET("/list/excel", authSvc.Permit("OFFICIAL_VIEW"), officialsExcel)
			og.GET("/doc", authSvc.Permit("OFFICIAL_VIEW"), officialDoc)
			og.GET("/docs", authSvc.Permit("OFFICIAL_VIEW"), officialDocs)
			og.GET("/docs/excel", authSvc.Permit("OFFICIAL_VIEW"), officialDocsExcel)
			og.GET("/doc/audit", authSvc.Permit("OFFICIAL_AUDIT"), officialDocAudit)
			og.GET("/doc/edit", authSvc.Permit("OFFICIAL_MNG"), officialDocEdit)

			ex := og.Group("/internal")
			ex.GET("/doc", vfySvc.Verify, officialDoc)
			ex.POST("/doc/audit", vfySvc.Verify, officialDocAudit)
			ex.POST("/doc/submit", vfySvc.Verify, officialDocSubmit)
		}

		rw := mg.Group("/review")
		{
			rw.GET("", authSvc.Permit("ACCOUNT_PROPERTY_REVIEW_AUDIT"), review)
			rw.GET("/list", authSvc.Permit("ACCOUNT_PROPERTY_REVIEW_AUDIT"), reviewList)
			rw.POST("/audit", authSvc.Permit("ACCOUNT_PROPERTY_REVIEW_AUDIT"), reviewAudit)

			rw.GET("/face/list", authSvc.Permit("ACCOUNT_PROPERTY_REVIEW_FACE_AUDIT"), reviewFaceList)
			rw.POST("/face/audit", authSvc.Permit("ACCOUNT_PROPERTY_REVIEW_FACE_AUDIT"), reviewAudit)
		}

		mn := mg.Group("/monitor")
		{
			mn.GET("/list", authSvc.Permit("ACCOUNT_MONITOR_REVIEW"), monitors)
			mn.POST("/add", authSvc.Permit("ACCOUNT_MONITOR_MNG"), addMonitor)
			mn.POST("/del", authSvc.Permit("ACCOUNT_MONITOR_MNG"), delMonitor)
		}

		rn := mg.Group("/realname")
		{
			rn.GET("/list", authSvc.Permit("REALNAME_QUERY"), realnameList)
			rn.GET("/list/pending", authSvc.Permit("REALNAME_AUDIT"), realnamePendingList)
			rn.GET("/image", authSvc.Permit("REALNAME_AUDIT"), realnameImage)
			rn.GET("/image/preview", realnameImagePreview)
			rn.POST("/apply/audit", authSvc.Permit("REALNAME_AUDIT"), realnameAuditApply)
			rn.GET("/reason", authSvc.Permit("REALNAME_AUDIT"), realnameReasonList)
			rn.POST("/reason/set", authSvc.Permit("REALNAME_SET_REASON"), realnameSetReason)
			rn.POST("/search/card", vfySvc.Verify, realnameSearchCard)
			rn.POST("/unbind", authSvc.Permit("ACCOUNT_REVIEW"), realnameUnbind)
			rn.GET("/excel", authSvc.Permit("REALNAME_EXPORT"), realnameExport)
			rn.POST("/file/upload", authSvc.Permit("REALNAME_SUBMIT"), realnameFileUpload)
			rn.POST("/submit", authSvc.Permit("REALNAME_SUBMIT"), realnameSubmit)

			ex := rn.Group("/internal", vfySvc.Verify)
			ex.GET("/list", realnameList)
			ex.GET("/list/pending", realnamePendingList)
			ex.GET("/image", realnameImage)
			ex.GET("/image/preview", realnameImagePreview)
			ex.POST("/apply/audit", realnameAuditApply)
			ex.GET("/reason", realnameReasonList)
			ex.POST("/reason/set", realnameSetReason)
			ex.POST("/search/card", realnameSearchCard)
			ex.POST("/unbind", realnameUnbind)
			ex.GET("/excel", realnameExport)
			ex.POST("/file/upload", realnameFileUpload)
			ex.POST("/submit", realnameSubmit)
		}
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("Failed to ping service: %+v", err)
		c.AbortWithStatus(ecode.Cause(err).Code())
		return
	}
	c.AbortWithStatus(200)
}

func register(c *bm.Context) {
	c.JSON(nil, nil)
}
