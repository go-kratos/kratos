package http

import (
	"net/http"

	"go-common/app/admin/main/apm/conf"
	per "go-common/app/admin/main/apm/model/user"
	"go-common/app/admin/main/apm/service"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	mpermit "go-common/library/net/http/blademaster/middleware/permit"
)

var (
	apmSvc  *service.Service
	authSrv *mpermit.Permit
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	apmSvc = s
	authSrv = mpermit.New(c.Auth)
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	b := e.Group("x/admin/apm", authSrv.Verify())
	{
		code := b.Group("/ecode")
		{
			code.POST("/add", permit(per.EcodeEdit), ecodeAdd)
			code.POST("/edit", permit(per.EcodeEdit), ecodeEdit)
			code.POST("/delete", permit(per.EcodeEdit), ecodeDelete)
			code.POST("/sync/ecodes", permit(per.EcodeEdit), syncEcodes)
			code.POST("/langs/add", permit(per.EcodeEdit), ecodeLangsAdd)
			code.POST("/langs/edit", permit(per.EcodeEdit), ecodeLangsEdit)
			code.POST("/langs/delete", permit(per.EcodeEdit), ecodeLangsDelete)
			code.POST("/langs/save", permit(per.EcodeEdit), ecodeLangsSave)
		}
		user := b.Group("/user")
		{
			user.GET("/auth", userAuth)
			user.GET("/rule/states", userRuleStates)
			user.POST("/apply", userApply)
			user.POST("/apply/edit", permit(per.UserEdit), userApplyEdit)
			user.GET("/users", userList)
			user.GET("/info", permit(per.UserView), userInfo)
			user.POST("/edit", permit(per.UserEdit), userEdit)
			user.GET("/modules", permit(per.UserEdit), userModules)
			user.GET("/rules", permit(per.UserEdit), userRules)
			user.POST("/module/edit", permit(per.UserEdit), userModulesEdit)
			user.POST("/rule/edit", permit(per.UserEdit), userRulesEdit)
			user.GET("/applies", userApplies)
			user.POST("/audit", userAudit)
			// sync
			user.GET("/tree/sync", userSyncTree)
			user.GET("/tree/appids", userTreeAppids)
			user.GET("/tree/discovery", userTreeDiscovery)
		}

		databus := b.Group("/databus")
		{
			// project
			databus.GET("/projects", permit(per.DatabusKeyView), databusProjects)
			// app
			databus.GET("/apps", permit(per.DatabusKeyView), databusApps)
			databus.POST("/app/add", permit(per.DatabusKeyEdit), databusAppAdd)
			databus.POST("/app/edit", permit(per.DatabusKeyEdit), databusAppEdit)
			// group
			databus.GET("/groups", permit(per.DatabusGroupView), databusGroups)
			databus.GET("/group/projects", permit(per.DatabusGroupView), databusGroupProjects)
			databus.GET("/group/consumer/addrs", permit(per.DatabusGroupView), databusConsumerAddrs)
			databus.POST("/group/sub/add", permit(per.DatabusGroupEdit), databusGroupSubAdd)
			databus.POST("/group/pub/add", permit(per.DatabusGroupEdit), databusGroupPubAdd)
			databus.POST("/group/delete", permit(per.DatabusGroupEdit), databusGroupDelete)
			databus.POST("/group/rename", permit(per.DatabusGroupEdit), databusGroupRename)
			databus.GET("/group/offset", permit(per.DatabusGroupView), databusGroupOffset)
			databus.POST("/group/marked", permit(per.DatabusGroupEdit), databusGroupMarked)
			databus.POST("/group/begin", permit(per.DatabusGroupEdit), databusGroupBegin)
			databus.POST("/group/new/offset", permit(per.DatabusGroupEdit), databusGroupNewOffset)
			databus.POST("/group/time", permit(per.DatabusGroupEdit), databusGroupTime)
			// topic
			databus.GET("/topics", permit(per.DatabusTopicView), databusTopics)
			databus.GET("/topic/names", permit(per.DatabusGroupView), databusTopicNames)
			databus.POST("/topic/add", permit(per.DatabusTopicEdit), databusTopicAdd)
			databus.POST("/topic/edit", permit(per.DatabusTopicEdit), databusTopicEdit)
			databus.GET("/topic/all", permit(per.DatabusGroupView), databusTopicAll)
			// alarm
			databus.POST("/alarm/edit", permit(per.DatabusGroupView), databusAlarmEdit)
			databus.POST("/alarm/init", permit(per.DatabusGroupView), databusAlarmInit)
			databus.POST("/alarm/all/edit", permit(per.DatabusGroupView), databusAlarmAllEdit)
			// apply
			databus.GET("/apply/list", permit(per.DatabusGroupView), databusApplyList)
			databus.POST("/apply/pub/add", permit(per.DatabusGroupView), databusApplyPubAdd)
			databus.POST("/apply/sub/add", permit(per.DatabusGroupView), databusApplySubAdd)
			databus.POST("/apply/edit", permit(per.DatabusGroupView), databusApplyEdit)
			databus.POST("/apply/approval/process", permit(per.DatabusGroupApply), databusApplyApprovalProcess)
			// notify
			databus.GET("/notify/apply/list", permit(per.DatabusGroupView), databusNotifyList)
			databus.POST("/notify/edit", permit(per.DatabusNotifyEdit), databusNotifyEdit)
			databus.POST("/notify/apply/add", permit(per.DatabusGroupView), databusNotifyApplyAdd)
			// databus.POST("/notify/filter/add", databusNotifyFilterAdd)
			// databus.POST("/notify/filter/edit", databusNotifyFilterEdit)
			// message
			databus.GET("/message/fetch", permit(per.DatabusGroupView), databusMsgFetch)
		}

		canal := b.Group("/canal")
		{
			canal.GET("", permit(per.CanalView), canalList)
			canal.POST("/add", permit(per.CanalEdit), canalAdd)
			canal.POST("/edit", permit(per.CanalEdit), canalEdit)
			canal.POST("/delete", permit(per.CanalEdit), canalDelete)
			canal.GET("/scan", permit(per.CanalView), canalScanByAddrFromConfig)
			e.POST("x/admin/apm/canal/apply/config", canalApplyDetailToConfig)
			canal.GET("/addrs", permit(per.CanalView), canalAddrAll)
			canal.POST("/apply/edit", permit(per.CanalView), canalApplyConfigEdit)
			canal.GET("/apply", permit(per.CanalView), canalApplyList)
			canal.POST("/apply/approval/process", permit(per.CanalEdit), canalApplyApprovalProcess)
		}

		need := b.Group("/need")
		{
			need.GET("/list", needList)
			need.POST("/add", needAdd)
			need.POST("/edit", needEdit)
			need.POST("/verify", permit(per.NeedVerify), needVerify)
			need.POST("/thumbsup", needThumbsUp)
			need.GET("/vote/list", permit(per.NeedVerify), needVoteList)
		}

		discovery := b.Group("/discovery")
		{
			discovery.GET("/", discoveryProxy)
		}

		app := b.Group("/app")
		{
			app.GET("/list", permit(per.AppView), appList)
			app.POST("/add", permit(per.AppEdit), appAdd)
			app.POST("/edit", permit(per.AppEdit), appEdit)
			app.POST("/delete", permit(per.AppEdit), appDelete)
			app.GET("/auth/list", permit(per.AppAuthView), appAuthList)
			app.POST("/auth/add", permit(per.AppEdit), appAuthAdd)
			app.POST("/auth/edit", permit(per.AppEdit), appAuthEdit)
			app.POST("/auth/delete", permit(per.AppEdit), appAuthDelete)
			app.GET("/caller/search", permit(per.AppEdit), appCallerSearch)
			//tree
			app.GET("/tree", permit(per.AppEdit), appTree)
		}
		platform := b.Group("/platform")
		{
			platform.GET("/search/get/", permit(per.PlatformSearchView), searchProxyGet)
			platform.POST("/search/post/", permit(per.PlatformSearchView), searchProxyPost)
			platform.GET("/reply/get/", permit(per.PlatformReplyView), replyProxyGet)
			platform.POST("/reply/post/", permit(per.PlatformReplyView), replyProxyPost)
			platform.GET("/tag/get/", permit(per.PlatformTagView), tagProxyGet)
			platform.POST("/tag/post/", permit(per.PlatformTagView), tagProxyPost)
			platform.GET("/bfs/get/", permit(per.BFSView), bfsProxyGet)
			platform.POST("/bfs/post/", permit(per.BFSEdit), bfsProxyPost)
			platform.GET("/reply/feed/get/", permit(per.PlatformReplyView), replyFeedProxyGet)
			platform.POST("/reply/feed/post/", permit(per.PlatformReplyView), replyFeedProxyPost)
		}
		p := b.Group("/pprof")
		{
			p.GET("/profile", permit(per.PerformanceManager), buildSvg)
			p.GET("/svg", permit(per.PerformanceManager), readSvg)
			p.GET("/heap", permit(per.PerformanceManager), heap)
			p.GET("/flame", permit(per.PerformanceManager), flame)
			p.GET("/", permit(per.PerformanceManager), pprof)
		}
		ut := b.Group("/ut")
		{
			ut.GET("/merge/list", utList)
			ut.GET("/detail/list", utDetail)
			ut.GET("/history/commit", utHistoryCommit)
			ut.GET("/rank/list", utRank)
			ut.GET("/rank/user", userRank)
			ut.GET("/quality/trend", utQATrend)
			ut.GET("/commits", utGeneralCommit)
			ut.GET("/dashboard/pkgs", utDashPkgsTree)
			ut.GET("/dashboard/curve", utDashCurve)
			ut.GET("/dashboard/histogram", utDashHistogram)
			ut.GET("/dashboard/histogram/user", utDashUserDetail)
			ut.GET("/app/list", utApps)
		}
		open := b.Group("/open")
		{
			open.GET("/get/", permit(per.OpenView), openProxyGet)
			open.POST("/post/", permit(per.OpenView), openProxyPost)
		}

	}
	// no auth
	d := e.Group("x/admin/apm")
	notAuth := d.Group("/")
	{
		notAuth.GET("canal/alarm", canalList) // 运维canal报警接口
		notAuth.GET("databus/clusters", databusClusters)
		notAuth.GET("databus/alarm", databusAlarm)   //	运维databus告警用
		notAuth.GET("databus/alarms", databusAlarms) // 有所有group diff，查询时间巨长2分钟起板
		notAuth.GET("ut/baseline", utBaseline)
		notAuth.GET("noauth/discovery/fetch", discoveryProxyNoAuth) // 提供给自动告警查询对应实例
		notAuth.GET("ecode", ecodeList)
		notAuth.POST("databus/opsmind", databusOpsmind)
		notAuth.POST("databus/opsmind/remove", databusOpsmindRemove)
		notAuth.GET("databus/query", databusQuery)
		notAuth.GET("ecode/get/ecodes", getEcodes)
		notAuth.GET("ecode/get/prod/ecodes", getProdEcodes)
		notAuth.GET("ecode/langs", codeLangsList)
	}
	dapper := d.Group("/dapper")
	{
		dapper.GET("/", dapperProxy)
	}
	up := d.Group("/ut")
	{
		up.POST("/upload", upload)
		up.POST("/upload/app", uploadApp)
		//up.GET("/tyrant", check)
		up.GET("/check", check)
		up.POST("/merge/set", utSetMerged)
		up.GET("/dashboard/history/commit", dashHistoryCommit)
		up.GET("/git/report", utGitReport)

	}
	m := d.Group("/monitor")
	{
		m.GET("/apps", appNameList)
		m.GET("/prometheus", prometheusList)
		m.GET("/broadcast", broadcastList)
		m.GET("/databus", databusList)
		m.GET("/online", onlineList)
	}
	warn := d.Group("/warn")
	{
		warn.POST("/active", activeWarning)
	}
}

// Paper paper.
type Paper struct {
	Total int         `json:"total"`
	Pn    int         `json:"pn"`
	Ps    int         `json:"ps"`
	Items interface{} `json:"items"`
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := apmSvc.Ping(c); err != nil {
		c.AbortWithStatus(http.StatusServiceUnavailable)
		log.Error("apm-admin service ping error(%v)", err)
	}
}

func permit(rule string) bm.HandlerFunc {
	return func(ctx *bm.Context) {
		usernameI, _ := ctx.Get("username")
		username, ok := usernameI.(string)
		if !ok || username == "" {
			ctx.JSON(nil, ecode.NoLogin)
			ctx.Abort()
			return
		}
		if rule == "" {
			return
		}
		if err := apmSvc.Permit(ctx, username, rule); err != nil {
			log.Error("apmSvc.Permit(%s, %s) error(%v)", username, rule, err)
			ctx.JSON(nil, err)
			ctx.Abort()
		}
	}
}
