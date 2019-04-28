package http

import (
	"net/http"

	webucenter_http "go-common/app/interface/live/web-ucenter/api/http"
	"go-common/app/interface/live/web-ucenter/api/http/v1"

	"go-common/app/interface/live/web-ucenter/conf"
	"go-common/app/interface/live/web-ucenter/dao"
	webcenterSvc "go-common/app/interface/live/web-ucenter/service/v1"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	srv     *webcenterSvc.Service
	vfy     *verify.Verify
	midAuth *auth.Auth
	// AnchorTask .
	AnchorTask *webcenterSvc.AnchorTaskService
)

// Init init
func Init(c *conf.Config) {
	dao.InitAPI()
	initService(c)
	initMidWare(c)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}
func initService(c *conf.Config) {
	srv = webcenterSvc.New(c)
	midAuth = auth.New(c.Auth)
}
func initMidWare(c *conf.Config) {
	vfy = verify.New(c.Verify)
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/xlive/web-ucenter")
	{
		g.GET("/Auth", midAuth.User, howToStart)
	}
	v1.RegisterV1HistoryService(e, srv, map[string]bm.HandlerFunc{"auth": midAuth.UserWeb})
	midMap := map[string]bm.HandlerFunc{
		"auth":  midAuth.User,
		"guest": midAuth.Guest,
	}
	v1.RegisterV1CapsuleService(e, webcenterSvc.NewCapsuleService(conf.Conf), midMap)
	v1.RegisterV1AnchorTaskService(e, webcenterSvc.NewAnchorTaskService(conf.Conf), midMap)

	webucenter_http.RegisterUserService(
		e, webcenterSvc.NewUserService(conf.Conf), map[string]bm.HandlerFunc{"auth": midAuth.User})
}

func ping(c *bm.Context) {
	c.AbortWithStatus(http.StatusOK)
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

// example for http request handler
func howToStart(c *bm.Context) {
	c.String(0, "Golang 大法好 !!!")
}
