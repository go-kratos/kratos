package http

import (
	"net/http"

	"go-common/app/job/openplatform/article/conf"
	"go-common/app/job/openplatform/article/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var ajSrv *service.Service

// Init .
func Init(conf *conf.Config, srv *service.Service) {
	ajSrv = srv
	// init outer router
	engineOuter := bm.DefaultServer(conf.BM)
	outerRouter(engineOuter)
	if err := engineOuter.Start(); err != nil {
		log.Error("xhttp.Serve error(%v)", err)
		panic(err)
	}
}

// outerRouter init outer router
func outerRouter(r *bm.Engine) {
	r.Ping(ping)
	r.Register(register)
	cr := r.Group("/sitemap")
	{
		cr.GET("/read/detail.xml", sitemap)
	}
}

func ping(c *bm.Context) {
	if err := ajSrv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

func sitemap(c *bm.Context) {
	res, l := ajSrv.SitemapXML(c)
	c.Writer.Header().Set("Content-Type", "text/xml")
	c.Writer.Header().Set("Content-Length", l)
	c.Status(http.StatusOK)
	c.Writer.Write([]byte(res))
}
