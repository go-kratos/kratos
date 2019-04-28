package http

import (
	"net/http"

	"go-common/app/service/main/tag/conf"
	"go-common/app/service/main/tag/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var svr *service.Service

// Init http init and router init .
func Init(c *conf.Config, s *service.Service) {
	svr = s
	engine := bm.DefaultServer(c.BM)
	router(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start() error(%v)", err)
		panic(err)
	}
}
func router(e *bm.Engine) {
	e.Ping(ping)
}

func ping(c *bm.Context) {
	if svr.Ping(c) != nil {
		log.Error("tag-service service ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
