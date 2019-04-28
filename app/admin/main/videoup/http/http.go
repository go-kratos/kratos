package http

import (
	"go-common/app/admin/main/videoup/conf"
	"go-common/app/admin/main/videoup/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	vfySvc  *verify.Verify
	authSrc *permit.Permit
	vdaSvc  *service.Service
)

// Init http server
func Init(c *conf.Config, s *service.Service) {
	vdaSvc = s
	vfySvc = verify.New(nil)
	authSrc = permit.New(c.Auth)
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	// init internal server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// innerRouter
func innerRouter(e *bm.Engine) {
	e.Ping(ping)

	group2 := e.Group("/x/admin/videoup", authSrc.Permit(""))
	{
		task := group2.Group("/task")
		{
			task.GET("/weightconfig/maxweight", authSrc.Permit("TASKWEIGHT"), maxweight)
			task.POST("/weightconfig/add", authSrc.Permit("TASKWEIGHT"), addwtconf)
			task.POST("/weightconfig/del", authSrc.Permit("TASKWEIGHT"), delwtconf)
			task.GET("/weightconfig/list", authSrc.Permit("TASKWEIGHT"), listwtconf)
			task.GET("/weightlog/list", authSrc.Permit("TASKWEIGHT"), listwtlog)
			task.GET("/wcv/show", authSrc.Permit("TASKWEIGHT"), show)
			task.POST("/wcv/set", authSrc.Permit("TASKWEIGHT"), set)

			// 用户在线
			task.GET("/consumer/on", checkgroup(), on)
			task.GET("/consumer/off", checkgroup(), off) //自己退出
			task.POST("/consumer/forceoff", forceoff)    //强制踢出
			task.GET("/online", authSrc.Permit("ONLINE"), online)
			task.GET("/inoutlist", inoutlist)
			// 任务状态
			task.POST("/delay", checkowner(), delay)
			task.POST("/free", taskfree)
		}
		oversea := group2.Group("/oversea")
		{
			oversea.GET("/policy/groups", policyGroups)                                                      //策略组列表
			oversea.POST("/policy/group/add", authSrc.Permit("POLICY_GROUP_EDIT"), addPolicyGroup)           //新增策略组
			oversea.POST("/policy/group/edit", authSrc.Permit("POLICY_GROUP_EDIT"), editPolicyGroup)         //编辑策略组
			oversea.POST("/policy/groups/del", authSrc.Permit("POLICY_GROUP_EDIT"), delPolicyGroups)         //删除策略组
			oversea.POST("/policy/groups/restore", authSrc.Permit("POLICY_GROUP_EDIT"), restorePolicyGroups) //恢复被删除的策略组
			oversea.GET("/policies", policies)                                                               //获取策略列表
			oversea.GET("/archive/groups", archiveGroups)                                                    //稿件策略组列表
			oversea.POST("/policies/add", authSrc.Permit("POLICY_ITEM_EDIT"), addPolicies)                   //给组添加策略
			oversea.POST("/policies/del", authSrc.Permit("POLICY_ITEM_EDIT"), delPolicies)                   //删除策略
		}
		search := group2.Group("/search")
		{
			search.GET("video", searchVideo)
			search.GET("archive", searchArchive)
			search.POST("copyright", searchCopyright)
		}
		staff := group2.Group("/staff")
		{
			staff.GET("", staffs)
			staff.POST("/batch/submit", authSrc.Permit("ARC_STAFF"), batchStaff)
		}

		//监控规则
		monitor := group2.Group("/monitor")
		{
			monitor.GET("/rule/result", authSrc.Permit("MONITOR_RULE_READ"), monitorRuleResult) //查看监控规则结果
			monitor.GET("/rule/result/oids", monitorRuleResultOids)
			monitor.POST("/rule/update", authSrc.Permit("MONITOR_RULE_EDIT"), monitorRuleUpdate) //修改监控规则
			monitor.GET("/notify", monitorNotify)
		}

		group2.POST("/archive/uptag", upArcTag)
		group2.POST("/archive/batch/tag", authSrc.Permit("BATCH_ARC_CHANNEL"), batchTag)
		group2.GET("/archive/channel/info", channelInfo)
	}
	group := e.Group("/va", vfySvc.Verify)
	{
		group.POST("/video/audit", videoAudit)
		group.POST("/video/batch/video", batchVideo)
		group.POST("/video/add", upVideo)
		group.POST("/video/change/index", changeIndex)
		group.POST("/video/del", delVideo)
		group.POST("/video/weblink/up", upWebLink)

		group.POST("/archive/submit", submit)
		group.POST("/archive/batch/archive", batchArchive)
		group.POST("/archive/batch/databus", batchArchiveSecondRound)
		group.POST("/archive/batch/attrs", batchAttrs)
		group.POST("/archive/batch/types", batchTypeIDs)
		group.POST("/archive/auther/up", upAuther)
		group.POST("/archive/access/up", upAccess)
		group.GET("/archive/flow/hit", hitFlows)
		group.GET("/archive/aitrack", aiTrack)

		group.POST("/cm/attr/up", upCMArr)
		group.POST("/cm/dtime/up", upCMArcDelay)

		group.POST("/pgc/pass", passByPGC)
		group.POST("/pgc/modify", modifyByPGC)
		group.POST("/pgc/lock", lockByPGC)

		group.GET("/track/archive", trackArchive)
		group.GET("/track/video", trackVideo)
		group.GET("/track/detail", trackDetail)

		group.GET("/task/tooks", taskTooks)
		group.GET("/stats/points", statsPoints)

		group.GET("/task/next", next)
		group.GET("/task/list", list)
		//同步音乐库  add,edit,delete ,log full sync tools
		group.POST("/music/sync", syncMusic)
		group.GET("/task/info", info)

		monitor := group.Group("/monitor")
		{
			monitor.GET("/notify", monitorNotify)
		}
	}
}
