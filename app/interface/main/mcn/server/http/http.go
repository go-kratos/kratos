package http

import (
	"net/http"

	"go-common/app/interface/main/mcn/conf"
	"go-common/app/interface/main/mcn/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	srv        *service.Service
	authSvc    *auth.Auth
	uploadAnti *antispam.Antispam
	verifySvc  *verify.Verify
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

func initService(c *conf.Config) {
	srv = service.New(c)
	authSvc = auth.New(nil)
	uploadAnti = antispam.New(c.UploadAntispam)
	verifySvc = verify.New(nil)
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	// 以下接口在 api.bilibili.com，对外使用
	g := e.Group("/x/mcn")
	{
		//g.GET("/start", vfy.Verify, howToStart)
		g.GET("/state", authSvc.User, mcnState)
		g.GET("/exist", authSvc.User, mcnExist)
		g.POST("/file/upload", multipartForm, authSvc.User, uploadAnti.ServeHTTP, upload)
		g.GET("/account/info", authSvc.User, mcnGetAccountInfo)
		g.GET("/base/info", authSvc.User, mcnBaseInfo)

		g.POST("/apply", authSvc.User, mcnApply)
		g.POST("/mcn/bindup", authSvc.User, mcnBindUpApply)
		g.GET("/mcn/get_data_summary", authSvc.User, mcnGetDataSummary)
		g.GET("/mcn/get_data_up_list", authSvc.User, mcnGetDataUpList)
		g.GET("/mcn/get_old_info", authSvc.User, mcnGetOldInfo)
		g.POST("/mcn/permit/change", authSvc.User, mcnGetChangePermit)
		g.POST("/mcn/publication/change-price", authSvc.User, mcnPublicationPriceChange)

		g.POST("/up/confirm", authSvc.User, mcnUpConfirm)
		g.GET("/up/get_bind", authSvc.User, mcnUpGetBind)
		g.POST("/up/permit/confirm-reauth", authSvc.User, mcnUpPermitApplyConfirm)
		g.GET("/up/permit/get-reauth", authSvc.User, mcnPermitApplyGetBind)

		g.GET("/rank/up_fans", authSvc.User, mcnGetRankUpFans)
		g.GET("/rank/archive_likes", authSvc.User, mcnGetRankArchiveLikesOuter)
		g.GET("/recommend/list", authSvc.User, mcnGetRecommendPool)
		g.GET("/recommend/list_tids", authSvc.User, mcnGetRecommendPoolTidList)

		g.GET("/data/index/inc", authSvc.User, mcnGetMcnGetIndexInc)
		g.GET("/data/index/source", authSvc.User, mcnGetMcnGetIndexSource)
		g.GET("/data/play/source", authSvc.User, mcnGetPlaySource)
		g.GET("/data/fans", authSvc.User, mcnGetMcnFans)
		g.GET("/data/fans/inc", authSvc.User, mcnGetMcnFansInc)
		g.GET("/data/fans/dec", authSvc.User, mcnGetMcnFansDec)
		g.GET("/data/fans/attention/way", authSvc.User, mcnGetMcnFansAttentionWay)

		// mcn粉丝和游客的粉丝分析
		g.GET("/data/fans/base/attr", authSvc.User, mcnGetBaseFansAttrReq)
		g.GET("/data/fans/area", authSvc.User, mcnGetFansArea)
		g.GET("/data/fans/type", authSvc.User, mcnGetFansType)
		g.GET("/data/fans/tag", authSvc.User, mcnGetFansTag)

		// mcn创作中心数据分析
		g.GET("/creative/archives", authSvc.User, archives)
		g.GET("/creative/archive/history/list", authSvc.User, archiveHistoryList)
		g.GET("/creative/archive/videos", authSvc.User, archiveVideos)
		g.GET("/creative/data/archive", authSvc.User, dataArchive)
		g.GET("/creative/data/videoquit", authSvc.User, dataVideoQuit)
		g.GET("/creative/danmu/distri", authSvc.User, danmuDistri)
		g.GET("/creative/data/base", authSvc.User, dataBase)
		g.GET("/creative/data/trend", authSvc.User, dataTrend)
		g.GET("/creative/data/action", authSvc.User, dataAction)
		g.GET("/creative/data/fan", authSvc.User, dataFan)
		g.GET("/creative/data/pandect", authSvc.User, dataPandect)
		g.GET("/creative/data/survey", authSvc.User, dataSurvey)
		g.GET("/creative/data/playsource", authSvc.User, dataPlaySource)
		g.GET("/creative/data/playanalysis", authSvc.User, dataPlayAnalysis)
		g.GET("/creative/data/article/rank", authSvc.User, dataArticleRank)
	}

	cmd := e.Group("/cmd")
	{
		cmd.GET("/reload_rank", cmdReloadRank)
	}

	// 以下接口在 api.bilibili.co，内部使用
	internal := e.Group("/x/internal/mcn")
	{
		internal.GET("/rank/archive_likes", verifySvc.Verify, mcnGetRankArchiveLikesAPI)
		// mcn 数据概况
		internal.GET("/data/mcn/summary", verifySvc.Verify, getMcnSummaryAPI)
		internal.GET("/data/index/inc", verifySvc.Verify, getIndexIncAPI)
		internal.GET("/data/index/source", verifySvc.Verify, getIndexSourceAPI)
		internal.GET("/data/play/source", verifySvc.Verify, getPlaySourceAPI)
		internal.GET("/data/fans", verifySvc.Verify, getMcnFansAPI)
		internal.GET("/data/fans/inc", verifySvc.Verify, getMcnFansIncAPI)
		internal.GET("/data/fans/dec", verifySvc.Verify, getMcnFansDecAPI)
		internal.GET("/data/fans/attention/way", verifySvc.Verify, getMcnFansAttentionWayAPI)
		// mcn粉丝和游客的粉丝分析
		internal.GET("/data/fans/base/attr", verifySvc.Verify, getFansBaseFansAttrAPI)
		internal.GET("/data/fans/area", verifySvc.Verify, getFansAreaAPI)
		internal.GET("/data/fans/type", verifySvc.Verify, getFansTypeAPI)
		internal.GET("/data/fans/tag", verifySvc.Verify, getFansTagAPI)
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

func multipartForm(c *bm.Context) {
	c.Request.ParseMultipartForm(maxFileSize)
}
