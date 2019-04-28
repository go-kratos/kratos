package http

import (
	"errors"
	"net/http"

	"go-common/app/interface/live/push-live/conf"
	"go-common/app/interface/live/push-live/service"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	pushLiveSvr      *service.Service
	errInvalidParams = errors.New("invalid params")
)

// Init http
func Init(c *conf.Config, srv *service.Service) {
	pushLiveSvr = srv
	// init router
	engineInner := bm.DefaultServer(c.BM.Inner)
	innerRouter(engineInner)
	if err := engineInner.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router api path.
func innerRouter(e *bm.Engine) {
	// init api
	e.Ping(ping)

	// http internal api
	group := e.Group("/xlive/internal/push-live/")
	{
		// 用户每日推送额度占用
		group.POST("/limit/decrease", decrease)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := pushLiveSvr.Ping(c); err != nil {
		log.Error("[http.http|ping] push-live ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// decrease decrease user daily push limit
func decrease(c *bm.Context) {
	params := c.Request.Form
	business := params.Get("business")
	uuid := params.Get("uuid")
	targetID := params.Get("target_id")
	mids := params.Get("mids")

	// params check
	if business == "" || uuid == "" || targetID == "" || mids == "" {
		log.Error("[http.http|decrease] request params(%v) error(%v)", params, errInvalidParams)
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if err := pushLiveSvr.LimitDecrease(c, business, targetID, uuid, mids); err != nil {
		log.Error("[http.http|decrease] pushLiveSvr.LimitDecrease error(%v), params(%v)", err, params)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	log.Info("[http.http|decrease] decrease success, params(%v)", params)
	c.JSON(nil, nil)
}
