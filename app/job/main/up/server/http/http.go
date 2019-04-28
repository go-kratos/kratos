package http

import (
	"net/http"

	"go-common/app/job/main/up/conf"
	"go-common/app/job/main/up/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	svc *service.Service
)

// Init init
func Init(c *conf.Config) {
	svc = service.New(c)
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
	g := e.Group("/x/job/up")
	{
		g.GET("/job/run", runJob)
		g.GET("/start")
		g.GET("/job/warm-up", warmUp)
		g.GET("/job/warm-up-mid", warmUpMid)
		g.GET("/job/add-staff", addStaff)
		g.GET("/job/delete-staff", deleteStaff)
	}
}

func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
