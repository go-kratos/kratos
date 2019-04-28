package http

import (
	"net/http"

	"go-common/app/admin/main/filter/conf"
	"go-common/app/admin/main/filter/service"
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

// Init init http service
func Init(s *service.Service) {
	svc = s

	authSvc = permit.New(conf.Conf.Auth)
	vfySvc = verify.New(nil)

	engine := bm.DefaultServer(conf.Conf.BM)
	innerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/x/admin/filter", vfySvc.Verify)
	{
		g.POST("/add", filterAdd)
		g.POST("/del", filterDel)
		g.GET("/list", filterList)
		g.GET("/search", filterSearch) // 搜索过滤内容
		g.GET("/log", filterLog)
		g.POST("/edit", filterEdit)
		g.GET("/get", filterRuleByID)
		g.GET("/origin", filterOrigin)
		g.GET("/origins", filterOrigins)
		gArea := g.Group("/area")
		{
			gArea.GET("/list", areaList)
			gArea.POST("/add", areaAdd)
			gArea.POST("/edit", areaEdit)
			gArea.GET("/log", areaLog)
			gArea.GET("/group/list", areaGroupList)
			gArea.POST("/group/add", areaGroupAdd)
		}
		gKey := g.Group("/key")
		{
			gKey.POST("/add", keyAdd)          // 添加规则
			gKey.POST("/del", keyDelFid)       // 列表页删除
			gKey.GET("/editinfo", keyEditInfo) // 编辑信息
			gKey.POST("/edit", keyEdit)        // 提交编辑
			gKey.GET("/search", keySearch)     // 列表
			gKey.GET("/log", keyLog)           // 日志
		}
		gWhite := g.Group("/white")
		{
			gWhite.POST("/add", whiteAddArea)      // 增加业务白名单
			gWhite.POST("/del", whiteDel)          // 删除白名单
			gWhite.GET("/search", whiteSearch)     // 搜索白名单
			gWhite.GET("/editinfo", whiteEditInfo) // 编辑时获取白名单信息
			gWhite.POST("/edit", whiteEdit)        // 编辑白名单
			gWhite.GET("/log", whiteEditLog)       // 白名单编辑日志
		}
	}
	gAI := e.Group("/x/admin/filter/ai")
	{
		gAI.GET("/config", authSvc.Permit("FIlTER_AI_CONFIG"), aiConfig) // AI配置查看
		// cr.VerifyPost("/config/edit", aiConfigEdit) // AI配置修改
		gAI.GET("/white", authSvc.Permit("FIlTER_AI_WHITE"), aiWhite)                // AI白名单列表
		gAI.POST("/white/add", authSvc.Permit("FIlTER_AI_WHITE_ADD"), aiWhiteAdd)    // AI白名单添加
		gAI.POST("/white/edit", authSvc.Permit("FIlTER_AI_WHITE_EDIT"), aiWhiteEdit) // AI白名单编辑
		gAI.GET("/case/score", authSvc.Permit("FIlTER_AI_CASE_SCORE"), aiCaseScore)  // AI查询案例文本得分
		gAI.POST("/case/add", authSvc.Permit("FIlTER_AI_CASE_ADD"), aiCaseAdd)       // AI添加案例（badcase）
	}
}

func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("filter admin ping error(%+v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// register check server ok.
func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
