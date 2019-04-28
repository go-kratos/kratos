package http

import (
	uphttp "go-common/app/service/main/up/api/v1"
	"go-common/app/service/main/up/conf"
	"go-common/app/service/main/up/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	idfSvc *verify.Verify
	//Svc service.
	Svc     *service.Service
	authSrc *permit.Permit
)

// Init init account service.
func Init(c *conf.Config, s *service.Service) {
	// service
	initService(c, s)
	// init internal router
	innerEngine := bm.DefaultServer(c.BM.Inner)
	setupInnerEngine(innerEngine)
	registerUpEngine(innerEngine)
	uphttp.RegisterUpBMServer(innerEngine, s)
	// init internal server
	if err := innerEngine.Start(); err != nil {
		log.Error("innerEngine.Start error(%v)", err)
		panic(err)
	}
}

func initService(c *conf.Config, s *service.Service) {
	idfSvc = verify.New(nil)
	Svc = s
	authSrc = permit.New(c.Auth)
	Svc.SetAuthServer(authSrc)
}

// registerUpEngine .
func registerUpEngine(e *bm.Engine) {
	e.Inject("^/x/internal/uper/archive", idfSvc.Verify)
}

// innerRouter
func setupInnerEngine(e *bm.Engine) {
	// monitor ping
	e.Ping(ping)
	e.Register(disRegister)
	// base
	var base, admin *bm.RouterGroup
	if conf.Conf.IsTest {
		base = e.Group("/x/internal/uper")
		admin = e.Group("/x/admin/uper")
	} else {
		base = e.Group("/x/internal/uper", idfSvc.Verify)
		// 现在只要登录，默认放过
		admin = e.Group("/x/admin/uper", authSrc.Verify(), authSrc.Permit("PRIORITY_UP"))
	}
	{
		base.POST("/register", register)
		base.GET("/info", info)
		base.GET("/all", all)
		base.GET("/special", specialUps)
		base.GET("/stat/base", baseStat)
		base.GET("/info/active", active)
		base.GET("/info/actives", actives)
		//播放器关注按钮开关
		base.POST("/switch/set", switchSet)
		base.GET("/switch", upSwitch)
		//人物卡片
		base.GET("/card/all", listCardBase)
		base.GET("/card/info", getCardInfo)
		base.GET("/card/info/list", listCardDetail)
		base.GET("/card/info/list_by_mids", listCardByMids)

		// 下面接口在admin中也有
		base.GET("/special/get", specialGet)
		base.GET("/group/get", getGroup)
		base.GET("/special/get_by_mid", specialGetByMid)

		// list_up
		base.GET("/list_up", listUp)
	}

	admin.GET("/special/get", specialGet)
	admin.GET("/special/get_by_mid", specialGetByMid)
	admin.POST("/special/delete", specialDel)
	admin.POST("/special/add", specialAdd)
	admin.POST("/special/edit", specialEdit)
	admin.GET("/group/get", getGroup)
	admin.POST("/group/add", authSrc.Permit("UPGROUP_ADD"), addGroup)
	admin.POST("/group/update", updateGroup)
	admin.POST("/group/delete", authSrc.Permit("UPGROUP_ADD"), removeGroup)
}

// ping check server ok.
func ping(ctx *bm.Context) {
	if err := Svc.Ping(ctx); err != nil {
		ctx.Error = err
		ctx.AbortWithStatus(503)
	}
}

// disRegister check server ok.
func disRegister(ctx *bm.Context) {
	ctx.JSON(map[string]interface{}{}, nil)
}

//bmGetStringOrDefault get string
func bmGetStringOrDefault(c *bm.Context, key string, defaul string) (value string, exist bool) {
	i, exist := c.Get(key)

	if !exist {
		value = defaul
		return
	}

	value, exist = i.(string)
	if !exist {
		value = defaul
	}
	return
}
