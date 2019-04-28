package http

import (
	"net/http"

	"go-common/app/job/main/relation/conf"
	"go-common/app/job/main/relation/service"
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
	// // init inner server
	// if err := xhttp.Serve(intM, c.MultiHTTP.Inner); err != nil {
	// 	log.Error("xhttp.Serve error(%v)", err)
	// 	panic(err)
	// }

	innerEngine := bm.DefaultServer(c.BM)
	setupInnerEngine(innerEngine)

	if err := innerEngine.Start(); err != nil {
		panic(err)
	}
}

// innerRouter init local router api path.
// func innerRouter(r *router.Router) {
// 	r.Get("/monitor/ping", ping)
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
