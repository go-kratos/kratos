package http

import (
	"net/http"

	"go-common/app/admin/main/reply/conf"
	"go-common/app/admin/main/reply/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	rpSvc     *service.Service
	verifySvc *verify.Verify
)

// Init init http.
func Init(c *conf.Config) {
	// init services
	rpSvc = service.New(c)
	verifySvc = verify.New(c.Verify)
	auther := permit.New(c.ManagerAuth)
	engine := bm.DefaultServer(c.HTTPServer)
	authRouter(engine, auther)
	interRouter(engine)
	verifyRouter(engine)
	// serve port
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

type managerAuther interface {
	Permit(permit string) bm.HandlerFunc
}

type fakeAuth struct {
}

func (a *fakeAuth) Permit(permit string) bm.HandlerFunc {
	return func(ctx *bm.Context) {
		ctx.Next()
	}
}

func authRouter(engine *bm.Engine, auther managerAuther) {
	engine.GET("/monitor/ping", ping)
	group := engine.Group("/x/admin/reply")
	{
		// reply
		group.GET("/search", auther.Permit("REPLY_READONLY"), replySearch)      // 评论列表
		group.GET("/export", auther.Permit("REPLY_READONLY"), replyExport)      // 评论列表
		group.POST("/pass", auther.Permit("REPLY_MGR"), adminPassReply)         // 通过评论
		group.POST("/recover", auther.Permit("REPLY_MGR"), adminRecoverReply)   //  恢复评论
		group.POST("/edit", auther.Permit("REPLY_MGR"), adminEditReply)         //  编辑评论
		group.POST("/del", auther.Permit("REPLY_MGR"), adminDelReply)           //  删除评论
		group.POST("/top", auther.Permit("REPLY_TOP_MGR"), adminTopReply)       //  置顶评论
		group.GET("/top/log", auther.Permit("REPLY_TOP_MGR"), adminTopReplyLog) //  置顶评论日志搜索

		group.GET("/reply/top", auther.Permit("REPLY_MGR"), topChildReply)

		// report
		group.GET("/report/search", auther.Permit("REPLY_REPORT_READONLY"), reportSearch) // 举报列表
		group.POST("/report/del", auther.Permit("REPLY_MGR"), reportDel)                  // 举报删除
		group.POST("/report/ignore", auther.Permit("REPLY_MGR"), reportIgnore)            // 举报忽略
		group.POST("/report/recover", auther.Permit("REPLY_MGR"), reportRecover)          // 举报恢复
		group.POST("/report/transfer", auther.Permit("REPLY_MGR"), reportTransfer)        // 举报转审
		group.POST("/report/state", auther.Permit("REPLY_MGR"), reportStateSet)           // 设置举报状态
		// monitor
		group.GET("/monitor/search", auther.Permit("REPLY_MONITOR_READONLY"), monitorSearch) // 监控列表
		group.GET("/monitor/stats", auther.Permit("REPLY_MONITOR_READONLY"), monitorStats)   // 监控统计
		group.POST("/monitor/state", auther.Permit("REPLY_MONITOR_MGR"), monitorState)       // 监控状态
		group.GET("/mointor/log", auther.Permit("REPLY_MONITOR_READONLY"), monitorLog)       // 监控操作日志
		// config
		group.POST("/config/update", auther.Permit("REPLY_MGR"), updateReplyConfig)     // 配置更新
		group.POST("/config/renew", auther.Permit("REPLY_MGR"), renewReplyConfig)       // 配置添加
		group.GET("/config/info", auther.Permit("REPLY_READONLY"), loadReplyConfig)     // 配置信息
		group.GET("/config/list", auther.Permit("REPLY_READONLY"), paginateReplyConfig) // 配置列表
		// notice
		group.GET("/notice/detail", auther.Permit("MANAGER_NOTICE"), getNotice)       // 通知信息
		group.GET("/notice/list", auther.Permit("MANAGER_NOTICE"), listNotice2)       // 通知列表
		group.POST("/notice/edit", auther.Permit("MANAGER_NOTICE"), editNotice)       // 通知编辑
		group.POST("/notice/delete", auther.Permit("MANAGER_NOTICE"), deleteNotice)   // 通知删除
		group.POST("/notice/offline", auther.Permit("MANAGER_NOTICE"), offlineNotice) // 下线通知
		group.POST("/notice/online", auther.Permit("MANAGER_NOTICE"), onlineNotice)   // 上线通知
		// subject
		group.GET("/subject/info", auther.Permit("REPLY_READONLY"), adminSubject)              // 主题信息
		group.POST("/subject/state", auther.Permit("REPLY_SUBJECT_FREEZE"), adminSubjectState) // 主题状态设置
		group.POST("/subject/freeze", auther.Permit("REPLY_SUBJECT_FREEZE"), SubFreeze)        // 主题冻结或解冻评论
		group.GET("/subject/log", auther.Permit("REPLY_READONLY"), SubLogSearch)               // 主题操作日志搜索
		// action
		group.GET("/action/count", auther.Permit("REPLY_MGR"), actionCount)
		group.POST("/action/update", auther.Permit("REPLY_MGR"), actionUpdate)
		// log
		group.GET("/log", auther.Permit("REPLY_MGR"), logByRpID)

		// emoji
		group.GET("/emoji/list", auther.Permit("REPLY_EMOJI"), listEmoji)
		group.POST("/emoji/add", auther.Permit("REPLY_EMOJI"), createEmoji)
		group.POST("/emoji/state", auther.Permit("REPLY_EMOJI"), upEmojiState)
		group.POST("/emoji/sort", auther.Permit("REPLY_EMOJI"), upEmojiSort)
		group.POST("/emoji/edit", auther.Permit("REPLY_EMOJI"), upEmoji)
		group.GET("/emoji/package/list", auther.Permit("REPLY_EMOJI"), listEmojiPacks)
		group.POST("/emoji/package/add", auther.Permit("REPLY_EMOJI"), createEmojiPackage)
		group.POST("/emoji/package/edit", auther.Permit("REPLY_EMOJI"), editEmojiPack)
		group.POST("/emoji/package/sort", auther.Permit("REPLY_EMOJI"), upEmojiPackageSort)
		// business
		group.GET("/business/list", listBusiness)
		group.GET("/business/get", getBusiness)
		group.POST("/business/add", addBusiness)
		group.POST("/business/update", upBusiness)
		group.POST("/business/state", upBusiState)

		// fold reply
		group.POST("/fold", foldReply)
	}
}

func interRouter(engine *bm.Engine) {
	group := engine.Group("/x/internal/replyadmin", verifySvc.Verify)
	{
		//log
		group.GET("/log", logByRpID)
		//search
		group.GET("/reply/search", replySearch)
		group.GET("/report/search", reportSearch)
		// moniter
		group.GET("/monitor/stats", monitorStats)
		group.GET("/monitor/search", monitorSearch)
		group.POST("/monitor/state", monitorState)
		// config
		group.POST("/config/update", updateReplyConfig)
		group.POST("/config/rewnew", renewReplyConfig)
		group.GET("/config/info", loadReplyConfig)
		group.GET("/config/list", paginateReplyConfig)
		// notice
		group.GET("/notice/detail", getNotice)
		group.GET("/notice/list", listNotice)
		group.POST("/notice/edit", editNotice)
		group.POST("/notice/delete", deleteNotice)
		group.POST("/notice/offline", offlineNotice)
		group.POST("/notice/online", onlineNotice)
	}
}

func verifyRouter(engine *bm.Engine) {
	group := engine.Group("/x/admin/reply/internal")
	{
		group.POST("/spam", verifySvc.Verify, adminMarkAsSpam)          // 标记评论为垃圾
		group.POST("/del", verifySvc.Verify, adminDelReply)             // 删除评论
		group.POST("/callback/del", verifySvc.Verify, callbackDelReply) // 回调删除评论
		group.GET("/reply", adminReplyList)                             // 获得评论列表接口
	}
}

func ping(ctx *bm.Context) {
	if err := rpSvc.Ping(ctx); err != nil {
		log.Error("reply admin ping error(%v)", err)
		ctx.JSON(nil, err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
