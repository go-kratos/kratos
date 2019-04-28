package http

import (
	"net/http"

	"go-common/app/job/main/member/conf"
	"go-common/app/job/main/member/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var srv *service.Service

// Init init http service
func Init(c *conf.Config, s *service.Service) {
	srv = s
	// init inner router
	// intM := http.NewServeMux()
	// intR := router.New(intM)
	// innerRouter(intR)

	engine := bm.DefaultServer(c.BM)
	setupInnerEngine(engine)

	if err := engine.Start(); err != nil {
		panic(err)
	}
	// init inner server
	// if err := xhttp.Serve(intM, c.MultiHTTP.Inner); err != nil {
	// 	log.Error("xhttp.Serve error(%v)", err)
	// 	panic(err)
	// }
}

// innerRouter init local router api path.
// func innerRouter(r *router.Router) {
// 	r.MonitorPing(ping)
// }

func setupInnerEngine(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(ctx *bm.Context) {
	ctx.JSON(nil, nil)
}
