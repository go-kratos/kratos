package http

import (
	"net/http"

	"go-common/app/job/main/sms/conf"
	"go-common/app/job/main/sms/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var srv *service.Service

// Init init http service
func Init(c *conf.Config, s *service.Service) {
	srv = s
	engine := bm.DefaultServer(c.HTTPServer)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("sms-job ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
