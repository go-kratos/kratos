package http

import (
	"net/http"

	"go-common/app/service/openplatform/abtest/conf"
	"go-common/app/service/openplatform/abtest/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	abSvr *service.Service
)

// Init init account service.
func Init(c *conf.Config, s *service.Service) {
	abSvr = s
	engineInner := bm.DefaultServer(c.BM.Inner)
	innerRouter(engineInner)
	// init Inner serve
	if err := engineInner.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
	engineOuter := bm.DefaultServer(c.BM.Outer)
	outerRouter(engineOuter)
	// init Outer serve
	if err := engineOuter.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router api path.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	group := e.Group("/openplatform/internal/abtest")
	{
		//common
		group.GET("/versionid", versionID)
		group.GET("/version", version)
	}

	group = e.Group("/openplatform/admin/abtest")
	{
		//mng abtest
		group.GET("/list", listAb)
		group.GET("/add", addAb)
		group.GET("/update", updateAb)
		group.GET("/status", updateStatus)
		group.GET("/delete", deleteAb)
	}

	group = e.Group("/openplatform/admin/abtest/group")
	{
		//mng group
		group.GET("/add", addGroup)
		group.GET("/list", listGroup)
		group.GET("/update", updateGroup)
		group.GET("/delete", deleteGroup)
	}

	group = e.Group("/openplatform/admin/abtest/stat")
	{
		group.GET("/total", total)
	}
}

// outerRouter init outer router.
func outerRouter(e *bm.Engine) {
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := abSvr.Ping(c); err != nil {
		log.Error("open-abtest http service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
