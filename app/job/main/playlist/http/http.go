package http

import (
	"net/http"

	"go-common/app/job/main/playlist/conf"
	"go-common/app/job/main/playlist/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var pjSrv *service.Service

// Init .
func Init(c *conf.Config, s *service.Service) {
	pjSrv = s
	engineOut := bm.DefaultServer(c.HTTPServer)
	outerRouter(engineOut)
	// init Outer server
	if err := engineOut.Start(); err != nil {
		log.Error("engineOut.Start error(%v)", err)
		panic(err)
	}
}

func outerRouter(e *bm.Engine) {
	e.Ping(ping)
}

func ping(c *bm.Context) {
	if err := pjSrv.Ping(c); err != nil {
		log.Error("playlist job ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
