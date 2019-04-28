package http

import (
	"net/http"

	"go-common/app/job/main/tag/conf"
	"go-common/app/job/main/tag/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var svc *service.Service

// Init http server .
func Init(c *conf.Config, s *service.Service) {
	svc = s
	engine := bm.DefaultServer(c.BM)
	router(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

func router(e *bm.Engine) {
	e.Ping(ping)
}

func ping(c *bm.Context) {
	if svc.Ping(c) != nil {
		log.Error("tag-job ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
