package http

import (
	"go-common/app/job/main/ugcpay/conf"
	"go-common/app/job/main/ugcpay/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	srv *service.Service
)

// Init init
func Init(s *service.Service) {
	srv = s
	engine := bm.DefaultServer(conf.Conf.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%+v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
}

func ping(c *bm.Context) {
}
