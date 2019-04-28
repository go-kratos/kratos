package http

import (
	"go-common/app/infra/notify/conf"
	mrl "go-common/app/infra/notify/model"
	"go-common/app/infra/notify/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	svc *service.Service
)

// Init init
func Init(c *conf.Config) {
	initService(c)
	// init router
	eng := bm.DefaultServer(c.BM)
	initRouter(eng)
	if err := eng.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	svc = service.New(c)
}

// initRouter init outer router api path.
func initRouter(e *bm.Engine) {
	e.Ping(ping)
	group := e.Group("/x/internal/notify")
	{
		group.POST("/pub", pub)
	}
}

func pub(c *bm.Context) {
	arg := new(mrl.ArgPub)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, svc.Pub(c, arg))
}

// ping check server ok.
func ping(c *bm.Context) {
}
