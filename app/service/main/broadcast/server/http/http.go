package http

import (
	"net/http"

	"go-common/app/service/main/broadcast/service"
	"go-common/library/conf/paladin"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	c         paladin.Map
	srv       *service.Service
	verifySvr *verify.Verify
)

// Init init.
func Init(s *service.Service) {
	var hc struct {
		Server *bm.ServerConfig
	}
	if err := paladin.Get("http.toml").UnmarshalTOML(&hc); err != nil {
		if err != paladin.ErrNotExist {
			panic(err)
		}
	}
	if err := paladin.Watch("application.toml", &c); err != nil {
		panic(err)
	}
	srv = s
	verifySvr = verify.New(nil)
	engine := bm.DefaultServer(hc.Server)
	outerRouter(engine)
	interRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("xhttp.Serve2 error(%v)", err)
		panic(err)
	}
}

func outerRouter(e *bm.Engine) {
	g := e.Group("/x/broadcast")
	{
		g.POST("/conn/connect", connect)
		g.POST("/conn/disconnect", disconnect)
		g.POST("/conn/heartbeat", heartbeat)
		g.POST("/online/renew", renewOnline)
	}
}

func interRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/x/internal/broadcast")
	{
		g.POST("/push/keys", verifySvr.Verify, pushKeys)
		g.POST("/push/mids", verifySvr.Verify, pushMids)
		g.POST("/push/room", verifySvr.Verify, pushRoom)
		g.POST("/push/all", verifySvr.Verify, pushAll)
		g.GET("/online/top", onlineTop)
		g.GET("/online/room", onlineRoom)
		g.GET("/online/total", onlineTotal)
		g.GET("/server/list", serverList)
		g.GET("/server/infos", serverInfos)
		g.GET("/server/weight", serverWeight)
	}
}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("broadcast-service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
