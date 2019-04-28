package http

import (
	"net/http"
	"strconv"

	"go-common/app/interface/main/push-archive/conf"
	"go-common/app/interface/main/push-archive/model"
	"go-common/app/interface/main/push-archive/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	pushSrv *service.Service
	authSrv *auth.Auth
	veriSrv *verify.Verify
	anti    *antispam.Antispam
)

// Init init http.
func Init(c *conf.Config, srv *service.Service) {
	pushSrv = srv
	eng := bm.DefaultServer(c.Bm)
	authSrv = auth.New(c.Auth)
	veriSrv = verify.New(c.Verify)
	anti = antispam.New(c.Anti)
	addRoutes(eng)

	if err := eng.Start(); err != nil {
		log.Error("eng.Start error(%v)", err)
		panic(err)
	}
}

func addRoutes(e *bm.Engine) {
	e.Ping(ping)

	// TODO delete, for test
	e.GET("/x/push-archive/test", func(ctx *bm.Context) {
		aidStr := ctx.Request.Form.Get("aid")
		aid, _ := strconv.ParseInt(aidStr, 10, 64)
		midStr := ctx.Request.Form.Get("mid")
		mid, _ := strconv.ParseInt(midStr, 10, 64)
		pushSrv.Test(&model.Archive{
			ID:  aid,
			Mid: mid,
		})
	}) // TODO delete, for test

	set := e.Group("/x/push-archive/setting", veriSrv.Verify, authSrv.UserMobile)
	{
		set.GET("/get", anti.ServeHTTP, setting)
		set.POST("/set", setSetting)
	}
}

func ping(c *bm.Context) {
	if err := pushSrv.Ping(c); err != nil {
		log.Error("push-archive ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
