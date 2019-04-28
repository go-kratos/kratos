package http

import (
	"go-common/app/interface/main/broadcast/conf"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// Init init http.
func Init(c *conf.Config) {
	engine := bm.DefaultServer(c.HTTP)
	outerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

func outerRouter(e *bm.Engine) {
	e.Ping(ping)
}

func ping(c *bm.Context) {

}
