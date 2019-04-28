package http

import (
	"net/http"

	"go-common/app/job/main/activity/conf"
	"go-common/app/job/main/activity/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var ajSrv *service.Service

// Init .
func Init(conf *conf.Config, srv *service.Service) {
	ajSrv = srv
	engine := bm.DefaultServer(conf.BM)
	outerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("httpx.Serve(%v) error(%+v)", conf.BM, err)
		panic(err)
	}
}

func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.GET("/match/finish", finishMatch)
}

func ping(c *bm.Context) {
	if err := ajSrv.Ping(c); err != nil {
		log.Error("activity-job ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func finishMatch(c *bm.Context) {
	v := new(struct {
		MoID int64 `form:"mo_id" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, ajSrv.FinishMatch(c, v.MoID))
}
