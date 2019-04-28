package http

import (
	"net/http"

	"go-common/app/job/main/app-wall/conf"
	"go-common/app/job/main/app-wall/service/offer"
	"go-common/app/job/main/app-wall/service/unicom"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	offerSvc  *offer.Service
	unicomSvc *unicom.Service
)

// Init init
func Init(c *conf.Config) {
	initService(c)
	// init router
	engineInner := bm.DefaultServer(c.BM.Inner)
	outerRouter(engineInner)
	if err := engineInner.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	offerSvc = offer.New(c)
	unicomSvc = unicom.New(c)
}

// Close
func Close() {
	offerSvc.Close()
	unicomSvc.Close()
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	//init api
	e.Ping(ping)
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := offerSvc.Ping(c); err != nil {
		log.Error("app-wall-job service ping error(%+v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
