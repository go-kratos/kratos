package http

import (
	"go-common/app/admin/main/growup/conf"
	"go-common/app/admin/main/growup/service"
	"go-common/app/admin/main/growup/service/income"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	svr       *service.Service
	incomeSvr *income.Service
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	svr = s
	incomeSvr = income.New(conf.Conf)
	engine := bm.DefaultServer(c.BM)
	initRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func initRouter(r *bm.Engine) {
	r.Ping(ping)
	// 在up-profit.bilibili.co域名下
	b := r.Group("/allowance/api/x/admin/growup")
	up := b.Group("/up")
	{
		up.GET("/list", queryForUps)
		up.POST("/add", add)
		up.POST("/reject", reject)
		up.POST("/pass", pass)
		up.POST("/dismiss", dismiss)
		up.POST("/forbid", forbid)
		up.POST("/recovery", recovery)
		up.POST("/delete", deleteUp)
		up.POST("/state", updateAccountState)
		up.POST("/account/delete", delUpAccount)
		up.POST("/account/update", updateUpAccount)
		up.GET("/export", exportUps)
		up.POST("/white/add", addWhite)
		up.GET("/account/state", upState)
		credit := up.Group("/credit")
		{
			credit.GET("/list", creditRecords)
			credit.POST("/recover", recoverCredit)
		}
	}
	block := b.Group("/block")
	{
		block.GET("/list", queryFromBlocked)
		block.POST("/add", addToBlocked)
		block.POST("/delete", deleteFromBlocked)
	}
	authority := b.Group("/authority")
	{
		authority.GET("/user/privileges", getAuthorityUserPrivileges)
		authority.GET("/user/groups", getAuthorityUserGroup)

		authority.GET("/user/list", listAuthorityUsers)
		authority.POST("/user/add", addAuthorityUser)
		authority.POST("/user/update/info", updateAuthorityUserInfo)
		authority.POST("/user/update/auth", updateAuthorityUserAuth)
		authority.POST("/user/delete", deleteAuthorityUser)

		authority.GET("/taskgroup/list", listAuthorityTaskGroups)
		authority.POST("/taskgroup/add", addAuthorityTaskGroup)
		authority.POST("/taskgroup/add/user", addAuthorityTaskGroupUser)
		authority.POST("/taskgroup/update/info", updateAuthorityTaskGroupInfo)
		authority.POST("/taskgroup/delete", deleteAuthorityTaskGroup)
		authority.POST("/taskgroup/delete/user", deleteAuthorityTaskGroupUser)
		authority.GET("/taskgroup/list/privilege", listAuthorityGroupPrivilege)
		authority.POST("/taskgroup/update/privilege", updateAuthorityGroupPrivilege)

		authority.GET("/taskrole/list", listAuthorityTaskRoles)
		authority.POST("/taskrole/add", addAuthorityTaskRole)
		authority.POST("/taskrole/add/user", addAuthorityTaskRoleUser)
		authority.POST("/taskrole/update/info", updateAuthorityTaskRoleInfo)
		authority.POST("/taskrole/delete", deleteAuthorityTaskRole)
		authority.POST("/taskrole/delete/user", deleteAuthorityTaskRoleUser)
		authority.GET("/taskrole/list/privilege", listAuthorityRolePrivilege)
		authority.POST("/taskrole/update/privilege", updateAuthorityRolePrivilege)

		authority.GET("/list/groupandrole", listAuthorityGroupAndRole) //list task groups and task roles
		authority.GET("/privilege/list", listPrivilege)
		authority.POST("/privilege/add", addPrivilege)
		authority.POST("/privilege/update", updatePrivilege)

		// merge business and privilege
		authority.GET("/business", busPrivilege)
	}
	tag := b.Group("/tag")
	{
		tag.POST("/add", addTagInfo)
		tag.POST("/update", updateTagInfo)
		tag.POST("/state", modTagState)
		tag.POST("/addups", addTagUps)
		tag.POST("/release", releaseUp)
		tag.GET("/list", listTagInfo)
		tag.GET("/listup", listUps)
		tag.GET("/listav", listAvs)
		tag.GET("/details", tagDetails)
		tag.POST("/update/activity", updateActivity)
	}
	income := b.Group("/income")
	{
		income.GET("/up/list", upIncomeList)
		income.GET("/up/list/export", upIncomeListExport)
		income.GET("/up/statis", upIncomeStatis)
		income.GET("/archive/detail", archiveDetail)
		income.GET("/archive/statis", archiveStatis)
		income.GET("/archive/section", archiveSection)
		income.GET("/archive/top", archiveTop)
		income.GET("/bgm/detail", bgmDetail)
		income.POST("/archive/black", archiveBlack)
		income.POST("/archive/breach", archiveBreach)
		income.GET("/up/withdraw", upWithdraw)
		income.GET("/up/withdraw/export", upWithdrawExport)
		income.GET("/up/withdraw/statis", upWithdrawStatis)
		income.GET("/up/withdraw/detail", upWithdrawDetail)
		income.GET("/up/withdraw/detail/export", upWithdrawDetailExport)
		income.GET("/breach/list", breachList)
		income.GET("/breach/statis", breachStatis)
		income.GET("/breach/export", exportBreach)
		income.GET("/black/list", queryBlacklist)
		income.GET("/black/export", exportBlack)
		income.POST("/black/recover", recoverBlacklist)
	}
	notice := b.Group("/notice")
	{
		notice.GET("/list", notices)
		notice.POST("/add", insertNotice)
		notice.POST("/update", updateNotice)
	}
	cheat := b.Group("/cheat")
	{
		cheat.GET("/up", cheatUps)
		cheat.GET("/av", cheatArchives)
		cheat.GET("/export/up", exportCheatUps)
		cheat.GET("/export/av", exportCheatAvs)
		cheat.GET("/up/fans", queryCheatFans)
		cheat.POST("/up/info", cheatFans)
	}
	charge := b.Group("/charge")
	{
		charge.GET("/archive/statis", archiveChargeStatis)
		charge.GET("/archive/section", archiveChargeSection)
		charge.GET("/archive/detail", archiveChargeDetail)
		charge.GET("/bgm/detail", bgmChargeDetail)
		charge.GET("/up/ratio", upRatio)
	}
	b.POST("/upload", upload)

	budget := b.Group("/budget")
	{
		budget.GET("/day/info", budgetDayStatistics)
		budget.GET("/day/graph", budgetDayGraph)
		budget.GET("/month/info", budgetMonthStatistics)
	}

	banner := b.Group("/banner")
	{
		banner.GET("/list", banners)
		banner.POST("/add", addBanner)
		banner.POST("/edit", editBanner)
		banner.POST("/off", off)
	}

	activity := b.Group("/activity")
	{
		activity.POST("/add", activityAdd)
		activity.GET("/list", activityList)
		activity.POST("/update", activityUpdate)
		activity.GET("/sign_up", activitySignUp)
		activity.GET("/winners", activityWinners)
		activity.POST("/award", activityAward)
	}
	auto := b.Group("/auto")
	{
		auto.POST("/archive/breach", autoBreach)
		auto.POST("/up/dismiss", autoDismiss)
		auto.POST("/up/forbid", autoForbid)
	}
	offlineActivity := b.Group("/offlineactivity")
	{
		offlineActivity.POST("/add", offlineactivityAdd)
		offlineActivity.POST("/pre_add", offlineactivityPreAdd)
		offlineActivity.POST("/upload", uploadLocal)
		offlineActivity.GET("/query/activity", offlineactivityQueryActivity)
		offlineActivity.GET("/query/upbonus_summary", offlineactivityQueryUpBonusSummary)
		offlineActivity.GET("/query/upbonus_activity", offlineactivityQueryUpBonusActivity)
		offlineActivity.GET("/query/month", offlineactivityQueryMonth)
	}

	// 激励金兑换
	trade := b.Group("/trade")
	{
		// goods query
		trade.GET("/goods/list", goodsList)
		// goods mng
		trade.GET("/goods/sync", goodsSync)
		trade.POST("/goods/update", goodsUpdate)
		trade.POST("/goods/display_set", goodsDisplay)
		// order query
		trade.GET("/order/list", orderList)
		trade.GET("/order/export", orderExport)
		trade.GET("/order/statistics", orderStatistics)
	}

	// 专项奖
	award := b.Group("/special_award")
	{
		award.POST("/add", awardAdd)                      //新增专项奖
		award.POST("/update", awardUpdate)                //编辑专项奖
		award.GET("/list", awardList)                     //专项奖列表
		award.GET("/detail", awardDetail)                 //专项奖详情
		award.GET("/winner/list", awardWinnerList)        //获奖者列表
		award.GET("/winner/export", awardWinnerExport)    //获奖者导出
		award.POST("/winner/replace", awardWinnerReplace) //获奖者替换
		award.GET("/result", awardResult)                 //评奖信息查询
		award.POST("/result/save", awardResultSave)       //评奖信息录入
	}

	// 在api.bilibili.co域名下，建议迁移到service
	api := r.Group("/x/internal/growup")
	{
		api.GET("/offlineactivity/callback", offlineactivityShellCallback)
	}
}
