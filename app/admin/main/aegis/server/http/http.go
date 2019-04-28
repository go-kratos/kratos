package http

import (
	"net/http"

	"go-common/app/admin/main/aegis/conf"
	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	srv     *service.Service
	vfy     *verify.Verify
	authSvr *permit.Permit
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	srv = s
	vfy = verify.New(c.Verify)
	authSvr = permit.New(c.Auth)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.GET("/x/admin/aegis/debug/cache", perm(""), cache)
	gb := e.Group("/x/admin/aegis/business", perm(""), preHandlerUser())
	{
		gb.GET("/view", checkBizAdmin(), getBusiness)
		gb.GET("/list", checkBizAdmin(), getBusinessList)
		gb.GET("/enable", checkBizBID(), getBusinessEnable)
		gb.POST("/add", authSvr.Permit("AEGIS_REGISTER"), addBusiness)
		gb.POST("/update", checkBizAdmin(), updateBusiness)
		gb.POST("/set", checkBizAdmin(), setBusinessState)

		gb.GET("config/list", listBizCFGs)
		gb.GET("config/reserve", reserveCFG)
		gb.POST("config/add", addBizCFG)
		gb.POST("config/update", updateBizCFG)
	}

	gt := e.Group("/x/admin/aegis/task", perm(""), preHandlerUser(), checkTaskRole())
	{
		gt.POST("/delay", taskDelay)
		gt.POST("/release", taskRelease)

		gt.POST("/consumer/on", consumerOn)
		gt.POST("/consumer/off", consumerOff)
		gt.POST("/consumer/kickout", kickOut)

		gt.GET("/undostat", taskUnDo)
		gt.GET("/stat", taskStat)
	}

	gtc := e.Group("/x/admin/aegis/task/config", perm(""), preHandlerUser(), checkBizID())
	{
		gtc.GET("/maxweight", maxWeight)
		gtc.GET("/list", configList)
		gtc.POST("/delete", configDelete)
		gtc.POST("/add", configAdd)
		gtc.POST("/edit", configEdit)
		gtc.POST("/set", configSet)
		gtc.GET("/weightlog", weightlog)
	}

	engineTask := e.Group("/x/admin/aegis/engine", perm(""), preHandlerUser(), checkTaskRole())
	{
		// 任务列表操作
		engineTask.GET("/task/next", checkon(), next)
		engineTask.GET("/task/info", infoByTask)
		engineTask.GET("/task/list", listByTask)
	}
	engineRsc := e.Group("/x/admin/aegis/engine", perm(""), preHandlerUser())
	{
		engineRsc.POST("/submit", checkBizBID(), submit)
		// 任务列表下面的flow列表
		engineRsc.GET("/task/listbizflow", checkAccessTask(), listBizFlow)

		// 资源列表操作
		engineRsc.GET("/resource/info", checkBizBID(), infoByResource)
		engineRsc.GET("/resource/list", checkBizBID(), listByResource)
		engineRsc.POST("/resource/batchsubmit", checkBizLeader(), batchSubmit)
		//跳流程
		engineRsc.GET("/resource/listforjump", checkBizBID(), listforjump)
		engineRsc.POST("/resource/jump", checkBizBID(), jump)

		// 封禁相关
		engineRsc.POST("/forbid/img/upload", upload)

		//信息追踪
		engineRsc.GET("/track", checkBizBID(), track)
		//操作日志
		engineRsc.GET("/auditlog", checkBizBID(), auditLog)
		engineRsc.GET("/auditlog/csv", checkBizBID(), auditLogCSV)
		//权限查询
		engineRsc.GET("/auth", auth)
		//手动取消, 运营人员权限
		engineRsc.GET("/cancel", checkBizOper(), cancelByOper)
	}

	gn := e.Group("/x/admin/aegis/net", perm(""), checkBizID())
	{
		gn.GET("", listNet)
		gn.GET("/svg", svg)
		gn.GET("/nets", getNetByBusiness)
		gn.GET("/show", showNet)
		gn.POST("/add", addNet)
		gn.POST("/update", updateNet)
		gn.POST("/switch", switchNet)

		gto := gn.Group("/token")
		{
			gto.GET("", listToken)
			gto.GET("/groups", tokenGroupByType)
			gto.GET("/config", configToken)
			gto.GET("/show", showToken)
			gto.POST("/add", addToken)
		}

		gf := gn.Group("/flow")
		{
			gf.GET("", listFlow)
			gf.GET("/flows", getFlowByNet)
			gf.GET("/show", showFlow)
			gf.POST("/add", addFlow)
			gf.POST("/update", updateFlow)
			gf.POST("/switch", switchFlow)
		}

		gtr := gn.Group("/transition")
		{
			gtr.GET("", listTransition)
			gtr.GET("/trans", getTranByNet)
			gtr.GET("/show", showTransition)
			gtr.POST("/add", addTransition)
			gtr.POST("/update", updateTransition)
			gtr.POST("/switch", switchTransition)
		}

		gd := gn.Group("/direction")
		{
			gd.GET("", listDirection)
			gd.GET("/show", showDirection)
			gd.POST("/add", addDirection)
			gd.POST("/update", updateDirection)
			gd.POST("/switch", switchDirection)
		}
	}

	rt := e.Group("/x/admin/aegis/report", perm(""), preHandlerUser(), checkBizLeader())
	{
		rt.GET("/task/flow", taskflow)
		rt.GET("/task/flow/csv", taskflowCSV)
		rt.GET("/task/submit", taskSubmit)
		rt.GET("/task/submit/csv", taskSubmitCSV)
		rt.GET("/business/flows", getBizFlow)
	}
	e.GET("/x/admin/aegis/net/token/byname", perm(""), checkBizBIDBiz(), tokenByName)

	e.GET("/x/admin/aegis/task/role", perm(""), preHandlerUser(), role)
	e.POST("/x/admin/aegis/task/role/flush", roleFlush)
	e.GET("/x/admin/aegis/task/consumer/watch", consumerWatcher)

	e.POST("/x/internal/aegis/add", gray(), add)
	e.POST("/x/internal/aegis/cancel", cancel)
	e.POST("/x/internal/aegis/update", update)

	//监控平台
	gm := e.Group("/x/admin/aegis/monitor", perm(""))
	{
		gm.GET("/rule/result", perm("MONITOR_RULE_READ"), monitorRuleResult)  //查看监控规则结果
		gm.POST("/rule/update", perm("MONITOR_RULE_EDIT"), monitorRuleUpdate) //修改监控规则
	}
	inter := e.Group("/x/internal/aegis", vfy.Verify)
	{
		inter.GET("/monitor/result/oids", monitorResultOids) //获取满足监控时间的对象id
	}
}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func uid(c *bm.Context) int64 {
	if v, exist := c.Get("uid"); exist {
		return v.(int64)
	}
	return 0
}

func uname(c *bm.Context) string {
	if v, exist := c.Get("username"); exist {
		return v.(string)
	}
	return ""
}

func perm(per string) bm.HandlerFunc {
	if srv.Debug() == "local" {
		return func(ctx *bm.Context) {}
	}
	return authSvr.Permit(per)
}

func gray() bm.HandlerFunc {
	return func(ctx *bm.Context) {
		opt := &model.AddOption{}
		ctx.Bind(opt)

		if !srv.Gray(opt) {
			log.Info("opt(%+v) not hit gray", opt)
			ctx.JSON(nil, nil)
			ctx.Abort()
			return
		}
		log.Info("opt(%+v) hit gray", opt)
	}
}
