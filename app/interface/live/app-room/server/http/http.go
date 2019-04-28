package http

import (
	"go-common/app/interface/live/app-room/api/http/v1"
	"go-common/app/interface/live/app-room/conf"
	"go-common/app/interface/live/app-room/service"
	resSrv "go-common/app/interface/live/app-room/service/v1"
	v1Svc "go-common/app/interface/live/app-room/service/v1"
	dm "go-common/app/interface/live/app-room/service/v1/dm"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
)

var (
	srv       *service.Service
	midAuth   *auth.Auth
	dmservice *dm.DMService
)

// Init init
func Init(c *conf.Config, svc *service.Service) {
	srv = svc
	resSrv.Init(c)
	initService(c)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func initService(c *conf.Config) {
	dmservice = dm.NewDMService(c)
	midAuth = auth.New(c.Auth)
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/xlive/app-room")
	{
		g.GET("/v1/banner/getBanner", resSrv.GetBanner)
	}
	midMap := map[string]bm.HandlerFunc{
		"guest": midAuth.Guest,
		"auth":  midAuth.User,
	}
	v1.RegisterV1GiftService(e, v1Svc.NewGiftService(conf.Conf), midMap)
	v1.RegisterV1RoomNoticeService(e, resSrv.NewRoomNoticeService(conf.Conf), midMap)

	g.POST("/v1/dM/sendmsg", midAuth.User, sendMsgSendMsg)
	g.GET("/v1/dM/gethistory", getHistory)
}

func ping(c *bm.Context) {
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
