package http

import (
	"net/http"

	"go-common/app/job/main/passport-encrypt/conf"
	"go-common/app/job/main/passport-encrypt/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	srv *service.Service
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	engine := bm.DefaultServer(c.BM.Inner)
	innerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
}

// ping service
func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
