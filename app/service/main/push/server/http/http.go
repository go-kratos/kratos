package http

import (
	"net/http"

	"go-common/app/service/main/push/conf"
	"go-common/app/service/main/push/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	pushSrv *service.Service
	idfSrv  *verify.Verify
)

// Init init http.
func Init(c *conf.Config, srv *service.Service) {
	idfSrv = verify.New(c.Verify)
	pushSrv = srv
	engine := bm.DefaultServer(c.HTTPServer)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start() error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/x/internal/push-service", bm.CORS())
	{
		g.POST("/single", idfSrv.Verify, singlePush)
		// for 直播
		g.POST("/setting/set", idfSrv.Verify, setSettingInternal)
		// for 管理后台测试推送
		g.POST("/push", idfSrv.Verify, push)
		// for test
		g.POST("/test/token", idfSrv.Verify, testToken)
		// upload image
		g.POST("/upimg", upimg)
	}
}

func ping(c *bm.Context) {
	if err := pushSrv.Ping(c); err != nil {
		log.Error("push-service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
