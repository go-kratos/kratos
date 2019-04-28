package http

import (
	"net/http"

	"go-common/app/job/main/creative/conf"
	"go-common/app/job/main/creative/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	//Svc service.
	Svc *service.Service
)

// Init init account service.
func Init(c *conf.Config) {
	// service
	initService(c)
	Svc.InitCron()
	// init inner router
	eng := bm.DefaultServer(c.BM.Outer)
	innerRouter(eng)
	if err := eng.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

func initService(c *conf.Config) {
	Svc = service.New(c)
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/x/job/creative")
	{
		g.GET("/test", sendMsg)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := Svc.Ping(c); err != nil {
		log.Error("svr.Ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

func sendMsg(c *bm.Context) {
	arg := new(struct {
		Mids []int64 `form:"mids,split" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(Svc.TestSendMsg(c, arg.Mids), nil)
}
