package http

import (
	"net/http"

	"go-common/app/job/main/passport-sns/conf"
	"go-common/app/job/main/passport-sns/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	srv *service.Service
)

// Init init http sever instance.
func Init(c *conf.Config) {
	srv = service.New(c)
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start() error(%v)", err)
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
