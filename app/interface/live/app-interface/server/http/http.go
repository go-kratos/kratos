package http

import (
	"net/http"

	apiV1 "go-common/app/interface/live/app-interface/api/http/v1"
	apiV2 "go-common/app/interface/live/app-interface/api/http/v2"
	"go-common/app/interface/live/app-interface/conf"
	"go-common/app/interface/live/app-interface/dao"
	"go-common/app/interface/live/app-interface/service"
	v1index "go-common/app/interface/live/app-interface/service/v1"
	v1appConf "go-common/app/interface/live/app-interface/service/v1/app_conf"
	v2index "go-common/app/interface/live/app-interface/service/v2"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	srv        *service.Service
	indexV2Srv *v2index.IndexService
	vfy        *verify.Verify
	midAuth    *auth.Auth
)

// Init init
func Init(c *conf.Config) {
	srv = service.New(c)
	dao.InitAPI()
	initMiddleware(c)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func initMiddleware(c *conf.Config) {
	vfy = verify.New(c.Verify)
	midAuth = auth.New(nil)
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	// 上线后注释掉，方便调试代码
	e.GET("/test", test)
	g := e.Group("/x/app-interface")
	{
		g.GET("/start", vfy.Verify, howToStart)
	}

	midMap := map[string]bm.HandlerFunc{"auth": midAuth.User, "guest": midAuth.Guest, "verify": vfy.Verify}
	apiV1.RegisterV1IndexService(e, v1index.New(conf.Conf), midMap)
	apiV1.RegisterV1RelationService(e, v1index.NewRelationService(conf.Conf), midMap)
	//移动端获取配置通用接口
	apiV1.RegisterV1ConfigService(e, v1appConf.NewAppConfService(conf.Conf), midMap)
	// v2 首页
	indexV2Srv = v2index.NewIndexService(conf.Conf)
	apiV2.RegisterV2IndexService(e, indexV2Srv, midMap)

}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
	// some frontend rely logic ping
	if indexV2Srv.GetAllModuleInfoMapFromCache(c) == nil {
		log.Error("ping error(AllMInfoMap must not nil)")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func test(c *bm.Context) {
	if err := srv.Test(c); err != nil {
		log.Error("test error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

// example for http request handler
func howToStart(c *bm.Context) {
	c.String(0, "Golang 大法好 !!!")
}
