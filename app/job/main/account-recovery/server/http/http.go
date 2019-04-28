package http

import (
	"go-common/app/job/main/account-recovery/conf"
	"go-common/app/job/main/account-recovery/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	srv *service.Service
)

// Init init
func Init(c *conf.Config) {
	srv = service.New(c)
	engine := bm.DefaultServer(c.BM)
	router(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func router(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/x/account-recovery")
	{
		g.GET("/test", ping)
	}
}

func ping(c *bm.Context) {
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
