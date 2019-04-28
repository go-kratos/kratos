package http

import (
	"go-common/library/net/http/blademaster/middleware/auth"
	"net/http"

	"go-common/app/interface/live/app-ucenter/api/http/v1"
	"go-common/app/interface/live/app-ucenter/conf"
	"go-common/app/interface/live/app-ucenter/service"
	spSrv "go-common/app/interface/live/app-ucenter/service/v1"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	srv      *service.Service
	vfy      *verify.Verify
	roomSrv  *spSrv.RoomService
	midAuth  *auth.Auth
	topicSrv *spSrv.TopicService
	raSrv    *spSrv.RoomAdminService
)

// Init init
func Init(c *conf.Config) {
	srv = service.New(c)
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
	roomSrv = spSrv.NewRoomService(c)
	raSrv = spSrv.NewRoomAdminService(c)
	midAuth = auth.New(c.Auth)
	topicSrv = spSrv.NewTopicService(c)
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/xlive/app-ucenter/v1")
	{
		g.GET("/start", vfy.Verify, howToStart)
	}
	v1.RegisterV1RoomService(e, roomSrv, map[string]bm.HandlerFunc{"auth": midAuth.User})
	v1.RegisterV1TopicService(e, topicSrv, map[string]bm.HandlerFunc{"auth": midAuth.User})
	v1.RegisterV1RoomAdminService(e, raSrv, map[string]bm.HandlerFunc{"auth": midAuth.User})
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
