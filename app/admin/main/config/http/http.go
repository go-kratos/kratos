package http

import (
	"net/http"

	"go-common/app/admin/main/config/conf"
	"go-common/app/admin/main/config/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	svr     *service.Service
	authSrv *permit.Permit
)

// Init init config.
func Init(c *conf.Config, s *service.Service) {
	svr = s
	authSrv = permit.New(c.Auth)
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
	b := e.Group("x/admin/config", authSrv.Verify())
	d := e.Group("x/admin/config")
	{
		notAuth := d.Group("/")
		{
			notAuth.POST("home/config/update", configUpdate)
			notAuth.POST("home/tag/update", tagUpdate)
			notAuth.POST("canal/tag/update", canalTagUpdate)
			notAuth.POST("canal/name/configs", canalNameConfigs)
			notAuth.POST("canal/config/create", canalConfigCreate)
			notAuth.POST("caster/envs", casterEnvs)
			notAuth.GET("get/apps", getApps)
		}
		service := b.Group("/service")
		{
			service.POST("/token/set", setToken)
			service.GET("/host/infos", hosts)
			service.POST("/delete", clearhost)
		}
		app := b.Group("/app")
		{
			app.POST("/token/update", updateToken)
			app.POST("/create", create)
			app.GET("/apps", appList)
			app.GET("/envs", envs)
			app.GET("/nodeTree", nodeTree)
			app.POST("/zone/copy", zoneCopy)
			app.POST("/rename", rename)
			app.POST("/status", upAppStatus)
		}
		bu := b.Group("/build")
		{
			bu.POST("/create", createBuild)
			bu.POST("/tag/update", updateTag)
			bu.POST("/tagid/update", updateTagID)
			bu.GET("/builds", builds)
			bu.GET("/build", build)
			bu.POST("/delete", buildDel)
			bu.POST("/hosts/force", hostsForce)
			bu.POST("/clear/force", clearForce)
		}
		config := b.Group("/config")
		{
			config.POST("/create", createConfig)
			config.POST("/lint", lintConfig)
			config.POST("/value/update", updateConfValue)
			config.GET("/app/configs", configsByAppName)
			config.GET("/build/configs", configsByBuildID)
			config.GET("/tag/configs", configsByTagID)
			config.GET("/name/configs", configsByName)
			config.GET("/names", namesByAppName)
			config.GET("/configs", configs)
			config.GET("/all/search", configSearchAll)
			config.GET("/app/search", configSearchApp)
			config.GET("/refs", configRefs)
			config.GET("/value", value)
			config.GET("/diff", diff)
			config.POST("/delete", configDel)
			config.GET("/build/infos", configBuildInfos)
		}
		common := b.Group("/common")
		{
			common.POST("/create", createComConfig)
			common.POST("/value/update", updateComConfValue)
			common.GET("/name/configs", comConfigsByName)
			common.GET("/configs", configsByTeam)
			common.GET("/names", namesByTeam)
			common.GET("/value", comValue)
			common.GET("/envs", envsByTeam)
			common.GET("/app", appByTeam)
			common.GET("/tag/push", tagPush)
		}
		tags := b.Group("/tag")
		{
			tags.POST("/create", createTag)
			tags.GET("/last/tags", lastTags)
			tags.GET("/build/tags", tagsByBuild)
			tags.GET("/tag", tag)
			tags.GET("/config/diff", tagConfigDiff)
		}
		apm := b.Group("/apm")
		{
			apm.GET("/copy", apmCopy)
		}
		tree := b.Group("/tree")
		{
			tree.GET("/update", syncTree)
		}
	}
}

// ping check server ok.
func ping(ctx *bm.Context) {
	if err := svr.Ping(); err != nil {
		log.Error("service ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func user(c *bm.Context) (username string) {
	usernameI, _ := c.Get("username")
	username, _ = usernameI.(string)
	return
}
