package http

import (
	v12 "go-common/app/interface/live/web-room/api/http/v1"
	"go-common/app/interface/live/web-room/conf"
	"go-common/app/interface/live/web-room/service"
	"go-common/app/interface/live/web-room/service/v1"
	v1index "go-common/app/interface/live/web-room/service/v1/dm"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"net/http"
)

var (
	authn     *auth.Auth
	srv       *service.Service
	dmservice *v1index.DMService
)

// Init init
func Init(c *conf.Config) {
	srv = service.New(c)
	initService(c)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func initService(c *conf.Config) {
	dmservice = v1index.NewDMService(c)
	authn = auth.New(c.AuthN)
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/xlive/web-room")
	g.POST("/v1/dM/sendmsg", authn.User, sendMsgSendMsg)
	g.POST("/v1/dM/gethistory", getHistory)
	v12.RegisterV1CaptchaService(e, v1.NewCaptchaService(conf.Conf), map[string]bm.HandlerFunc{
		"auth": authn.User,
	})
	v12.RegisterV1RoomAdminService(e, v1.NewRoomAdminService(conf.Conf), map[string]bm.HandlerFunc{})
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
