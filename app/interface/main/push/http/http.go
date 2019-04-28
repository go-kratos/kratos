package http

import (
	"net/http"

	"go-common/app/interface/main/push/conf"
	"go-common/app/interface/main/push/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
)

var (
	pushSrv *service.Service
	authSrv *auth.Auth
)

// Init init http.
func Init(c *conf.Config, srv *service.Service) {
	pushSrv = srv
	authSrv = auth.New(c.Auth)
	engine := bm.DefaultServer(c.HTTPServer)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start() error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/x/push", bm.CORS())
	{
		// for APP client
		g.POST("/report", authSrv.GuestMobile, report)
		g.GET("/report_old", reportOld)
		g.POST("/setting/set", authSrv.UserMobile, setSetting)
		g.GET("/setting/get", authSrv.UserMobile, setting)
		// for test
		g.POST("/test/token", testToken)
		// for callback
		cg := g.Group("/callback")
		{
			cg.POST("/huawei", huaweiCallback)        // 华为送达回执
			cg.POST("/xiaomi", miCallback)            // 小米送达回执
			cg.POST("/xiaomi/regid", miRegidCallback) // 小米token注册回执
			cg.POST("/oppo", oppoCallback)            // oppo送达回执
			cg.POST("/jpush", jpushCallback)          // 极光送达回执
			cg.POST("/ios", iOSCallback)              // iOS送达回执
			cg.POST("/android", androidCallback)      // Android送达回执
			cg.POST("/click", clickCallback)          // 所有平台的点击回执
		}
	}
}

func ping(c *bm.Context) {
	if err := pushSrv.Ping(c); err != nil {
		log.Error("push-interface ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
