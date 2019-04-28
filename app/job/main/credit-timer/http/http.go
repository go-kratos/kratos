package http

import (
	"net/http"
	"strconv"

	"go-common/app/job/main/credit-timer/conf"
	"go-common/app/job/main/credit-timer/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var svc *service.Service

// Init init
func Init(c *conf.Config) {
	initService(c)
	// init  inner router
	engineInner := bm.DefaultServer(c.BM)
	innerRouter(engineInner)
	if err := engineInner.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	svc = service.New(c)
}

// innerRouter init local router api path.
func innerRouter(e *bm.Engine) {
	//init api
	e.GET("/monitor/ping", ping)
	e.POST("/fixkpi", fixkpi)
}

func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("answer interface ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func fixkpi(c *bm.Context) {
	params := c.Request.Form
	yStr := params.Get("year")
	mStr := params.Get("month")
	dStr := params.Get("day")
	midStr := params.Get("mid")
	y, _ := strconv.ParseInt(yStr, 10, 64)
	m, _ := strconv.ParseInt(mStr, 10, 64)
	d, _ := strconv.ParseInt(dStr, 10, 64)
	mid, _ := strconv.ParseInt(midStr, 10, 64)
	res, _ := svc.FixKPI(c, int(y), int(m), int(d), mid)
	c.JSON(res, nil)
}
