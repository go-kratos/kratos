package http

import (
	"go-common/app/admin/live/live-admin/dao"
	"go-common/library/net/http/blademaster/middleware/auth"
	"net/http"

	v1API "go-common/app/admin/live/live-admin/api/http/v1"
	v2API "go-common/app/admin/live/live-admin/api/http/v2"
	"go-common/app/admin/live/live-admin/conf"
	"go-common/app/admin/live/live-admin/service"
	v1Service "go-common/app/admin/live/live-admin/service/v1"
	v2Service "go-common/app/admin/live/live-admin/service/v2"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	vfy     *verify.Verify
	svc     *service.Service
	midAuth *auth.Auth
	d       *dao.Dao
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	svc = s
	dao.InitAPI()
	d = dao.New(c)
	vfy = verify.New(c.Verify)
	midAuth = auth.New(c.Auth)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/x/live-admin")
	{
		g.GET("/start", vfy.Verify, howToStart)
	}
	midMap := map[string]bm.HandlerFunc{
		"guest": midAuth.Guest,
		"cors":  bm.CORS(),
	}
	v1API.RegisterV1ResourceService(e, v1Service.NewResourceService(conf.Conf), midMap)
	v1API.RegisterV1CapsuleService(e, v1Service.NewCapsuleService(conf.Conf), midMap)
	v1API.RegisterV1GaeaService(e, v1Service.NewGaeaService(conf.Conf), midMap)
	v1API.RegisterV1RoomMngService(e, v1Service.NewRoomMngService(conf.Conf), midMap)
	v2API.RegisterV2UserResourceService(e, v2Service.NewUserResourceService(conf.Conf), midMap)
	v1API.RegisterV1PayGoodsService(e, v1Service.NewPayGoodsService(conf.Conf), midMap)
	v1API.RegisterV1PayLiveService(e, v1Service.NewPayLiveService(conf.Conf), midMap)
	v1API.RegisterV1TokenService(e, v1Service.NewTokenService(conf.Conf, d), midMap)
	v1API.RegisterV1UploadService(e, v1Service.NewUploadService(conf.Conf, d), midMap)
}

func ping(ctx *bm.Context) {
	if err := svc.Ping(ctx); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

// example for http request handler
func howToStart(c *bm.Context) {
	c.String(0, "Golang 大法好 !!!")
}
