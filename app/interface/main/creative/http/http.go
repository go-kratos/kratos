package http

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/server/grpc"
	"go-common/app/interface/main/creative/service"
	"go-common/app/interface/main/creative/service/academy"
	"go-common/app/interface/main/creative/service/account"
	"go-common/app/interface/main/creative/service/ad"
	"go-common/app/interface/main/creative/service/app"
	"go-common/app/interface/main/creative/service/appeal"
	"go-common/app/interface/main/creative/service/archive"
	"go-common/app/interface/main/creative/service/article"
	"go-common/app/interface/main/creative/service/assist"
	"go-common/app/interface/main/creative/service/danmu"
	"go-common/app/interface/main/creative/service/data"
	"go-common/app/interface/main/creative/service/dynamic"
	"go-common/app/interface/main/creative/service/elec"
	"go-common/app/interface/main/creative/service/faq"
	"go-common/app/interface/main/creative/service/feedback"
	"go-common/app/interface/main/creative/service/geetest"
	"go-common/app/interface/main/creative/service/medal"
	"go-common/app/interface/main/creative/service/music"
	"go-common/app/interface/main/creative/service/newcomer"
	"go-common/app/interface/main/creative/service/operation"
	"go-common/app/interface/main/creative/service/pay"
	"go-common/app/interface/main/creative/service/reply"
	"go-common/app/interface/main/creative/service/resource"
	"go-common/app/interface/main/creative/service/staff"
	"go-common/app/interface/main/creative/service/template"
	"go-common/app/interface/main/creative/service/up"
	"go-common/app/interface/main/creative/service/version"
	"go-common/app/interface/main/creative/service/watermark"
	"go-common/app/interface/main/creative/service/weeklyhonor"
	"go-common/app/interface/main/creative/service/whitelist"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc/warden"
)

var (
	//app service
	apSvc     *appeal.Service
	arcSvc    *archive.Service
	elecSvc   *elec.Service
	dataSvc   *data.Service
	accSvc    *account.Service
	tplSvc    *template.Service
	gtSvc     *geetest.Service
	replySvc  *reply.Service
	fdSvc     *feedback.Service
	operSvc   *operation.Service
	assistSvc *assist.Service
	artSvc    *article.Service
	mdSvc     *medal.Service
	wmSvc     *watermark.Service
	appSvc    *app.Service
	danmuSvc  *danmu.Service
	vsSvc     *version.Service
	whiteSvc  *whitelist.Service
	adSvc     *ad.Service
	musicSvc  *music.Service
	resSvc    *resource.Service
	rpcdaos   *service.RPCDaos
	acaSvc    *academy.Service
	faqSvc    *faq.Service
	dymcSvc   *dynamic.Service
	honorSvc  *weeklyhonor.Service
	paySvc    *pay.Service

	// api middleware
	verifySvc   *verify.Verify
	authSvc     *auth.Auth
	antispamSvc *antispam.Antispam
	dmAnti      *antispam.Antispam
	//up service
	upSvc *up.Service
	// grpc TODO mv out http
	grpcSvr     *warden.Server
	newcomerSvc *newcomer.Service
	pubSvc      *service.Public
	staffSvc    *staff.Service
)

// Init init account service.
func Init(c *conf.Config) {
	// service
	initService(c)
	// init grpc
	grpcSvr = grpc.New(nil, arcSvc, newcomerSvc)
	engineOuter := bm.DefaultServer(c.BM.Outer)
	// init outer router
	outerRouter(engineOuter)
	if err := engineOuter.Start(); err != nil {
		log.Error("engineOuter.Start() error(%v) | config(%v)", err, c)
		panic(err)
	}
}

//Close for close server
func Close() {
	grpcSvr.Shutdown(context.TODO())
}

func initService(c *conf.Config) {
	verifySvc = verify.New(nil)
	authSvc = auth.New(nil)
	antispamSvc = antispam.New(c.RouterAntispam)
	dmAnti = antispam.New(c.DmAntispam)
	// public for injection
	rpcdaos = service.NewRPCDaos(c)
	pubSvc = service.New(c, rpcdaos)
	// services
	apSvc = appeal.New(c, rpcdaos)
	arcSvc = archive.New(c, rpcdaos, pubSvc)
	elecSvc = elec.New(c, rpcdaos)
	dataSvc = data.New(c, rpcdaos, pubSvc)
	accSvc = account.New(c, rpcdaos)
	tplSvc = template.New(c, rpcdaos)
	operSvc = operation.New(c, rpcdaos)
	wmSvc = watermark.New(c, rpcdaos, pubSvc)
	gtSvc = geetest.New(c, rpcdaos)
	replySvc = reply.New(c, rpcdaos)
	fdSvc = feedback.New(c, rpcdaos)
	assistSvc = assist.New(c, rpcdaos)
	artSvc = article.New(c, rpcdaos)
	mdSvc = medal.New(c, rpcdaos, pubSvc)
	appSvc = app.New(c, rpcdaos, pubSvc)
	danmuSvc = danmu.New(c, rpcdaos)
	vsSvc = version.New(c, rpcdaos)
	whiteSvc = whitelist.New(c, rpcdaos)
	adSvc = ad.New(c, rpcdaos)
	musicSvc = music.New(c, rpcdaos, pubSvc)
	resSvc = resource.New(c, rpcdaos)
	acaSvc = academy.New(c, rpcdaos, pubSvc)
	upSvc = up.New(c, rpcdaos)
	faqSvc = faq.New(c, rpcdaos)
	dymcSvc = dynamic.New(c, rpcdaos)
	honorSvc = weeklyhonor.New(c, rpcdaos)
	paySvc = pay.New(c, rpcdaos)
	newcomerSvc = newcomer.New(c, rpcdaos)
	staffSvc = staff.New(c, rpcdaos)

}

func webDanmuRouter(g *bm.RouterGroup) {
	// manager
	g.GET("/danmu/list", webDmList)
	g.GET("/danmu/distri", webDmDistri)
	g.POST("/danmu/edit", dmAnti.ServeHTTP, webDmEdit)
	g.POST("/danmu/transfer", dmAnti.ServeHTTP, webDmTransfer)
	g.POST("/danmu/pool", dmAnti.ServeHTTP, webDmUpPool)
	// purchase
	g.GET("/danmu/purchases", webListDmPurchases)
	g.POST("/danmu/purchase/pass", dmAnti.ServeHTTP, webPassDmPurchase)
	g.POST("/danmu/purchase/deny", dmAnti.ServeHTTP, webDenyDmPurchase)
	g.POST("/danmu/purchase/cancel", dmAnti.ServeHTTP, webCancelDmPurchase)
	// report
	g.POST("/danmu/report/check", dmAnti.ServeHTTP, webDmReportCheck)
	g.GET("/danmu/report", webDmReport)
	// report
	g.GET("/danmu/protect/archive", webDmProtectArchive)
	g.GET("/danmu/protect/list", webDmProtectList)
	g.POST("/danmu/protect/operation", dmAnti.ServeHTTP, webDmProtectOper)
}

func appDanmuRouter(g *bm.RouterGroup) {
	g.GET("/danmu/list", authSvc.UserMobile, appDmList)
	g.GET("/danmu/recent", authSvc.UserMobile, appDmRecent)
	g.GET("/danmu/edit", authSvc.UserMobile, appDmEdit)
	g.POST("/danmu/edit/batch", authSvc.UserMobile, appDmEditBatch)
}

func academyRouter(g *bm.RouterGroup) {
	g.GET("/academy/archive/tags", webAcademyTags)
	g.GET("/academy/archive/list", webAcademyArchives)
	g.POST("/academy/feedback/add", webAddFeedBack)
}

//工单
func staffRouter(g *bm.RouterGroup) {
	//申请单交互
	g.POST("/staff/apply/submit", webApplySubmit)
	//staff 申请解除
	g.POST("/staff/apply/create", webApplyCreate)
}

func switchRouter(g *bm.RouterGroup) {
	g.POST("/switch/set", setUpSwitch)
	g.GET("/switch", upSwitch)
}

func webElecRouter(g *bm.RouterGroup) {
	g.GET("/elec/user", webUserElec)
	g.GET("/elec/notify", webElecNotify)
	g.GET("/elec/status", webElecStatus)
	g.GET("/elec/rank/recent", webElecRecentRank)
	g.GET("/elec/rank/current", webElecCurrentRank)
	g.GET("/elec/rank/toltal", webElecTotalRank)
	g.GET("/elec/dailybill", webElecDailyBill)
	g.GET("/elec/balance", webElecBalance)
	g.POST("/elec/status/set", webElecUpStatus)
	g.POST("/elec/user/update", webUserElecUpdate)
	g.POST("/elec/arc/update", webArcElecUpdate)
	g.GET("/elec/remark/list", webRemarkList)
	g.GET("/elec/remark/detail", webRemarkDetail)
	g.POST("/elec/remark/reply", webRemark)
	g.GET("/elec/recent", webRecentElec)
}

func webAssistRouter(g *bm.RouterGroup) {
	g.GET("/assist", webAssists)
	g.GET("/assist/status", webAssistStatus)
	g.GET("/assist/logs", webAssistLogs)
	g.POST("/assist/add", webAssistAdd)
	g.POST("/assist/del", webAssistDel)
	g.POST("/assist/set", webAssistSet)
	g.POST("/assist/log/revoc", webAssistLogRevoc)
}

func newcomerRouter(g *bm.RouterGroup) {
	g.GET("/newcomer/task/list", webTaskList)
	g.POST("/newcomer/reward/receive/add", webRewardReceive)
	g.POST("/newcomer/reward/receive/activate", webRewardActivate)
	g.GET("/newcomer/reward/receive/list", webRewardReceiveList)
	g.POST("/newcomer/task/bind", webTaskBind)
	g.GET("/newcomer/task/makeup", webTaskMakeup)
}

// outerRouter init inner router.
func outerRouter(e *bm.Engine) {
	e.GET("/monitor/ping", ping)
	g := e.Group("/x/web", authSvc.UserWeb)
	{
		webDanmuRouter(g)
		academyRouter(g)
		staffRouter(g)
		switchRouter(g)
		webElecRouter(g)
		webAssistRouter(g)
		newcomerRouter(g)
		g.GET("/ugcpay/protocol", webUgcPayProtocol)
		// mission
		g.GET("/mission/protocol", webMissionProtocol)
		// netsafe
		g.POST("/ns/md5", webNsMd5)
		//white
		g.GET("/white", webWhite)
		// archive.
		g.GET("/archive/parts", webArchVideos)
		g.GET("/archive/view", webViewArc)
		g.GET("/archives", webArchives)
		g.GET("/archive/staff/applies", webStaffApplies)
		g.GET("/archive/pre", webViewPre)
		g.GET("/archive/videos", webVideos)
		g.POST("/archive/delete", webDelArc)
		g.GET("/archive/tags", webTags)
		g.GET("/archive/desc/format", webDescFormat)
		// history
		g.GET("/archive/history/list", webHistoryList)
		g.GET("/archive/history/view", webHistoryView)
		// ad
		g.GET("/ad/game/list", webAdGameList)
		// appeal.
		g.GET("/appeal/list", webAppealList)
		g.GET("/appeal/detail", webAppealDetail)
		g.GET("/appeal/contact", webAppealContact)
		g.POST("/appeal/add", webAppealAdd)
		g.POST("/appeal/reply", antispamSvc.ServeHTTP, webAppealReply)
		g.POST("/appeal/down", webAppealDown)
		g.POST("/appeal/star", webAppealStar)
		// cover list.
		g.GET("/archive/covers", coverList)
		g.GET("/archive/recovers", webRecommandCover)
		// index.
		g.GET("/index/stat", webIndexStat)
		g.GET("/index/tool", webIndexTool)
		g.GET("/index/full", webIndexFull) //collect_arc
		g.GET("/index/notify", webIndexNotify)
		g.GET("/index/operation", webIndexOper)
		g.GET("/index/version", webIndexVersion)
		g.GET("/index/newcomer", webIndexNewcomer)
		// data
		g.GET("/data/videoquit", webVideoQuitPoints)
		g.GET("/data/archive", webArchive)
		g.GET("/data/article", webArticleData)
		g.GET("/data/base", base)
		g.GET("/data/trend", trend)
		g.GET("/data/action", action)
		g.GET("/data/survey", survey)
		g.GET("/data/pandect", pandect)
		g.GET("/data/fan", webFan)
		g.GET("/data/playsource", webPlaySource)
		g.GET("/data/playanalysis", webArcPlayAnalysis)
		g.GET("/data/article/thirty", webArtThirtyDay)
		g.GET("/data/article/rank", webArtRank)
		g.GET("/data/article/source", webArtReadAnalysis)
		// water mark
		g.GET("/watermark", waterMark)
		g.POST("/watermark/set", waterMarkSet)
		// feedback
		g.GET("/feedbacks", webFeedbacks)
		g.GET("/feedback/detail", webFeedbackDetail)
		g.GET("/feedback/tags", webFeedbackTags)
		g.GET("/feedback/newtags", webFeedbackNewTags)
		g.POST("/feedback/add", webFeedbackAdd)
		g.POST("/feedback/close", webFeedbackClose)
		// reply
		g.GET("/replies", replyList)
		// template.
		g.GET("/tpls", webTemplates)
		g.POST("/tpl/add", webAddTpl)
		g.POST("/tpl/update", webUpdateTpl)
		g.POST("/tpl/delete", webDelTpl)
		// fans medal
		g.GET("/medal/status", webMedalStatus)
		g.GET("/medal/recent", webRecentFans)
		g.POST("/medal/open", webMedalOpen)
		g.POST("/medal/check", webMedalCheck)
		g.GET("/medal/rank", webMedalRank)
		g.POST("/medal/rename", webMedalRename)
		g.GET("/medal/fans", webFansMedal)
		// article.
		g.GET("/article/author", webAuthor)
		g.GET("/article/view", webArticle)
		g.GET("/article/list", webArticleList)
		g.GET("/article/pre", webArticlePre)
		g.POST("/article/submit", webSubArticle)
		g.POST("/article/update", webUpdateArticle)
		g.POST("/article/delete", webDelArticle)
		g.POST("/article/withdraw", webWithDrawArticle)
		g.POST("/article/upcover", antispamSvc.ServeHTTP, webArticleUpCover)
		g.GET("/draft/view", webDraft)
		g.GET("/draft/list", webDraftList)
		g.POST("/draft/addupdate", webSubmitDraft)
		g.POST("/draft/delete", webDeleteDraft)
		g.POST("/article/capture", antispamSvc.ServeHTTP, webArticleCapture)
		// cm
		g.GET("/cm/oasis/stat", webCmOasisStat)
		// common
		g.GET("/user/mid", webUserMid)
		g.GET("/user/search", webUserSearch)
		//viewpoint
		g.GET("/viewpoints", webViewPoints)
		//g.POST("/viewpoints/edit", webViewPointsEdit)
	}
	h5 := e.Group("/x/h5")
	{
		// app h5 cooperate pager
		h5.GET("/cooperate/pre", authSvc.User, appCooperatePre)
		// bgm
		h5.GET("/bgm/ext", authSvc.User, appBgmExt)
		// faq
		h5.GET("/faq/editor", authSvc.User, appH5FaqEditor)
		h5.POST("/bgm/feedback", authSvc.User, appH5BgmFeedback)
		h5.GET("/elec/bill", authSvc.User, appElecBill)
		h5.GET("/elec/rank/recent", authSvc.User, appElecRecentRank)
		h5.GET("/medal/status", authSvc.User, appMedalStatus)
		h5.POST("/medal/check", authSvc.User, appMedalCheck)
		h5.POST("/medal/open", authSvc.User, appMedalOpen)
		h5.POST("/medal/rename", authSvc.User, appMedalRename)
		//academy
		h5.POST("/academy/play/add", authSvc.Guest, h5AddPlay)        //添加播放
		h5.POST("/academy/play/del", authSvc.Guest, h5DelPlay)        //删除播放
		h5.GET("/academy/play/list", authSvc.User, h5PlayList)        //我的课程
		h5.GET("/academy/play/view", authSvc.User, h5ViewPlay)        //查看我的课程
		h5.GET("/academy/theme/dir", h5ThemeDir)                      //主题课程目录 对应职业列表
		h5.GET("/academy/newb/course", h5NewbCourse)                  //新人课程
		h5.GET("/academy/tag", h5Tags)                                //标签目录
		h5.GET("/academy/archive", h5Archive)                         //课程列表
		h5.GET("/academy/feature", h5Feature)                         //精选课程
		h5.GET("/academy/recommend/v2", authSvc.Guest, h5RecommendV2) //推荐课程v2
		h5.GET("/academy/theme/course/v2", h5ThemeCousreV2)           //技能树（主题课程）v2
		h5.GET("/academy/keywords", h5Keywords)                       //搜索关键词提示
		// data center
		h5.GET("/data/archive", authSvc.User, appDataArc)
		h5.GET("/data/videoquit", authSvc.User, appDataVideoQuit)
		h5.GET("/data/fan", authSvc.User, appFan)                        //粉丝用户信息分析总览
		h5.GET("/data/fan/rank", authSvc.User, appFanRank)               //新粉丝排行榜
		h5.GET("/data/overview", authSvc.User, appOverView)              //新数据概览
		h5.GET("/data/archive/analyze", authSvc.User, appArchiveAnalyze) //新稿件数据分析
		h5.GET("/data/video/retention", authSvc.User, appVideoRetention) //新视频播放完成度
		h5.GET("/data/article", authSvc.User, appDataArticle)
		h5.GET("/archives/simple", authSvc.User, appSimpleArcVideos)
		// watermark
		h5.GET("/watermark", authSvc.User, waterMark)
		h5.POST("/watermark/set", authSvc.User, waterMarkSet)
		// up weekly honor
		h5.GET("/weeklyhonor", authSvc.Guest, weeklyHonor)
		// switch weekly honor subscribe
		h5.POST("/weeklyhonor/subscribe", authSvc.User, weeklyHonorSubSwitch)
		// task system
		h5.POST("/task/bind", authSvc.User, h5TaskBind)
		h5.GET("/task/list", authSvc.User, h5TaskList)
		h5.POST("/task/reward/receive", authSvc.User, h5RewardReceive)
		h5.POST("/task/reward/activate", authSvc.User, h5RewardActivate)
		h5.GET("/task/reward/list", authSvc.User, h5RewardReceiveList)
		h5.GET("/task/pub/list", authSvc.User, taskPubList) //其他业务方查看任务列表
	}
	app := e.Group("/x/app")
	{
		appDanmuRouter(app)
		// h5
		app.GET("/h5/pre", authSvc.User, appH5Pre)
		app.GET("/h5/mission/type", authSvc.User, appH5MissionByType)
		app.GET("/h5/archive/tags", authSvc.User, appH5ArcTags)
		app.GET("/h5/archive/tag/info", authSvc.User, appH5ArcTagInfo)
		app.GET("/banner", authSvc.User, appBanner)
		// archive
		app.GET("/mission/type", authSvc.UserMobile, appMissionByType)
		app.GET("/index", authSvc.User, appIndex)
		app.GET("/archives", authSvc.UserMobile, appArchives)
		app.GET("/archives/simple", authSvc.UserMobile, appSimpleArcVideos)
		app.GET("/up/info", authSvc.UserMobile, appUpInfo)
		// main app features
		app.GET("/pre", authSvc.User, appPre)
		app.GET("/archive/pre", authSvc.User, appArchivePre)
		app.GET("/archive/desc/format", authSvc.UserMobile, appArcDescFormat)
		app.GET("/archive/view", authSvc.UserMobile, appArcView)
		app.POST("/archive/delete", authSvc.UserMobile, appArcDel)
		// reply.
		app.GET("/replies", authSvc.UserMobile, appReplyList)
		// data
		app.GET("/data/archive", authSvc.UserMobile, appDataArc)
		app.GET("/data/videoquit", authSvc.UserMobile, appDataVideoQuit)
		app.GET("/data/fan", authSvc.UserMobile, appFan)
		app.GET("/data/fan/rank", authSvc.UserMobile, appFanRank)               //新粉丝排行榜
		app.GET("/data/overview", authSvc.UserMobile, appOverView)              //新数据概览
		app.GET("/data/archive/analyze", authSvc.UserMobile, appArchiveAnalyze) //新稿件数据分析
		app.GET("/data/video/retention", authSvc.UserMobile, appVideoRetention) //新视频播放完成度
		app.GET("/data/article", authSvc.UserMobile, appDataArticle)
		// elec
		app.GET("/elec/bill", authSvc.UserMobile, appElecBill)
		app.GET("/elec/rank/recent", authSvc.UserMobile, appElecRecentRank)
		// fans medal
		app.GET("/medal/status", authSvc.UserMobile, appMedalStatus)
		app.POST("/medal/check", authSvc.UserMobile, appMedalCheck)
		app.POST("/medal/open", authSvc.UserMobile, appMedalOpen)
		app.POST("/medal/rename", authSvc.UserMobile, appMedalRename)
		// article
		app.GET("/article/list", authSvc.UserMobile, appArticleList)
		// material
		app.GET("/material/pre", authSvc.UserMobile, appMaterialPre)
		app.GET("/material/view", authSvc.UserMobile, appMaterial)
		// bgm
		app.GET("/bgm/pre", authSvc.UserMobile, appBgmPre)
		app.GET("/bgm/list", authSvc.UserMobile, appBgmList)
		app.GET("/bgm/view", authSvc.UserMobile, appBgmView)
		app.GET("/bgm/search", authSvc.UserMobile, appBgmSearch)
		app.GET("/cooperate/view", authSvc.User, appCooperate)
		// task
		app.POST("/newcomer/task/bind", authSvc.UserMobile, appTaskBind)
	}
	cli := e.Group("/x/client", authSvc.User)
	{
		// archive.
		cli.GET("/archives", clientArchives)
		cli.GET("/archive/search", clientArchiveSearch)
		cli.GET("/archive/view", clientViewArc)
		cli.POST("/archive/delete", clientDelArc)
		cli.GET("/archive/pre", clientPre)
		cli.GET("/archive/tags", clientTags)
		// template.
		cli.GET("/tpls", clientTemplates)
		cli.POST("/tpl/add", clientAddTpl)
		cli.POST("/tpl/update", clientUpdateTpl)
		cli.POST("/tpl/delete", clientDelTpl)
		// cover list.
		cli.GET("/archive/covers", coverList)
	}
	geeg := e.Group("/x/geetest", authSvc.UserWeb)
	{
		// geetest.
		geeg.GET("/pre", gtPreProcess)
		geeg.POST("/validate", gtValidate)
		geeg.GET("/pre/add", gtPreProcessAdd)
	}
	creator := e.Group("/x/creator", authSvc.UserMobile)
	{
		// index
		creator.GET("/my", creatorMy)
		creator.GET("/index", creatorIndex)
		creator.GET("/earnings", creatorEarnings)
		creator.GET("/banner", creatorBanner)
		creator.GET("/replies", creatorReplyList)
		//archive
		creator.GET("/archives", creatorArchives)
		creator.GET("/archive/tag/info", creatorArcTagInfo)
		creator.GET("/archive/view", creatorViewArc)
		creator.GET("/archive/videoquit", creatorVideoQuit)
		creator.GET("/archive/data", creatorArchiveData)
		creator.POST("/archive/delete", creatorDelArc)
		creator.GET("/archive/pre", creatorPre)
		creator.GET("/archive/tags", creatorPredictTag)
		creator.GET("/archive/desc/format", creatorDescFormat)
		// article
		creator.GET("/article/pre", creatorArticlePre)
		creator.GET("/article/list", creatorArticleList)
		creator.GET("/article/view", creatorArticle)
		creator.POST("/article/delete", creatorDelArticle)
		creator.POST("/article/withdraw", creatorWithDrawArticle)
		creator.GET("/draft/list", creatorDraftList)
		// danmu
		creator.GET("/danmu/list", creatorDmList)
		creator.GET("/danmu/recent", creatorDmRecent)
		creator.POST("/danmu/edit", creatorDmEdit)
		creator.POST("/danmu/edit/batch", creatorDmEditBatch)
		//data
		creator.GET("/data/archive", creatorDataArchive)
		creator.GET("/data/article", creatorDataArticle)
	}

	i := e.Group("/x/internal/creative", verifySvc.Verify)
	{
		// TODO deprecated
		i.GET("/porder", upPorder)
		// for main app
		i.GET("/app/pre", appNewPre)
		// get order game info for app
		i.GET("/arc/commercial", arcCommercial)
		i.POST("/watermark/set", waterMarkSetInternal)
		i.GET("/order/game", arcOrderGameInfo)
		i.POST("/upload/material", uploadMaterial)
		i.POST("/join/growup/account", growAccountStateInternal)
		i.GET("/video/viewpoints", videoViewPoints)
		i.GET("/archive/bgm", arcBgmList)
		i.GET("/archive/staff", arcStaff)
		i.GET("/archive/vote", voteAcsByTime)

		//联合投稿配置
		i.GET("/staff/config", staffConfig)

		// data
		i.GET("/data/videoquit", setContextMid, webVideoQuitPoints)
		i.GET("/data/archive", setContextMid, webArchive)
		i.GET("/data/article", setContextMid, webArticleData)
		i.GET("/data/base", setContextMid, base)
		i.GET("/data/trend", setContextMid, trend)
		i.GET("/data/action", setContextMid, action)
		i.GET("/data/survey", setContextMid, survey)
		i.GET("/data/pandect", setContextMid, pandect)
		i.GET("/data/fan", setContextMid, webFan)
		i.GET("/data/playsource", setContextMid, webPlaySource)
		i.GET("/data/playanalysis", setContextMid, webArcPlayAnalysis)
		i.GET("/data/article/thirty", setContextMid, webArtThirtyDay)
		i.GET("/data/article/rank", setContextMid, webArtRank)
		i.GET("/data/article/source", setContextMid, webArtReadAnalysis)

		// archive
		i.GET("/archives", setContextMid, webArchives)
		// videos
		i.GET("/archive/videos", setContextMid, webVideos)

		// history
		i.GET("/archive/history/list", setContextMid, webHistoryList)

		// danmu
		i.GET("/danmu/distri", setContextMid, webDmDistri)

		// up weekly honor
		i.GET("/task/pub/list", setContextMid, taskPubList) //其他业务方查看任务列表
	}
}
