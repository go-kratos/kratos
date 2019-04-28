package http

import (
	"net/http"

	"go-common/library/net/http/blademaster/middleware/auth"

	"go-common/app/interface/live/app-blink/api/http/v1"
	"go-common/app/interface/live/app-blink/conf"
	"go-common/app/interface/live/app-blink/service"
	spSrv "go-common/app/interface/live/app-blink/service/v1"
	resRpc "go-common/app/service/live/resource/api/grpc/v1"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	srv      *service.Service
	vfy      *verify.Verify
	slSrv    *spSrv.SplashService
	bnSrv    *spSrv.BannerService
	roomSrv  *spSrv.RoomService
	midAuth  *auth.Auth
	topicSrv *spSrv.TopicService
)

// Init init
func Init(c *conf.Config, svc *service.Service) {
	srv = svc
	vfy = verify.New(c.Verify)
	initService(c)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func initService(c *conf.Config) {
	slSrv = spSrv.NewSplashService(c)
	bnSrv = spSrv.NewBannerService(c)
	roomSrv = spSrv.NewRoomService(c)
	midAuth = auth.New(c.Auth)
	topicSrv = spSrv.NewTopicService(c)
}

func getInfo(ctx *bm.Context) {
	p := new(resRpc.GetInfoReq)
	if err := ctx.Bind(p); err != nil {
		return
	}
	resp, err := slSrv.GetInfo(ctx, p)
	ctx.JSON(resp, err)
}

func getBlinkBanner(ctx *bm.Context) {
	p := new(resRpc.GetInfoReq)
	if err := ctx.Bind(p); err != nil {
		return
	}
	resp, err := bnSrv.GetBlinkBanner(ctx, p)
	ctx.JSON(resp, err)
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/xlive/app-blink/v1")
	{
		g.GET("/start", vfy.Verify, howToStart)
	}
	g.GET("/splash/getInfo", getInfo)
	g.GET("/banner/getBlinkBanner", getBlinkBanner)
	v1.RegisterV1RoomService(e, roomSrv, map[string]bm.HandlerFunc{"auth": midAuth.User})
	v1.RegisterV1TopicService(e, topicSrv, map[string]bm.HandlerFunc{"auth": midAuth.User})
}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
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
