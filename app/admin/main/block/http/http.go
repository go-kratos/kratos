package http

import (
	"go-common/app/admin/main/block/conf"
	"go-common/app/admin/main/block/model"
	"go-common/app/admin/main/block/service"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"

	"github.com/pkg/errors"
)

var (
	authSvc *permit.Permit
	svc     *service.Service
)

// Init http server
func Init() {
	initService()
	engine := bm.DefaultServer(conf.Conf.BM)
	innerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start() error(%v)", err)
		panic(err)
	}
}

// initService init biz services
func initService() {
	authSvc = permit.New(conf.Conf.Auth)
	svc = service.New()
}

func innerRouter(e *bm.Engine) {
	e.GET("/monitor/ping", func(c *bm.Context) {})
	cb := e.Group("/x/admin/block", authSvc.Permit(conf.Conf.Perms.Perm["search"]))
	{
		cb.POST("/search", blockSearch)
		cb.GET("/history", blockHistory)
	}
	cb = e.Group("/x/admin/block", authSvc.Permit(conf.Conf.Perms.Perm["block"]))
	{
		cb.POST("", batchBlock)
	}
	cb = e.Group("/x/admin/block", authSvc.Permit(conf.Conf.Perms.Perm["remove"]))
	{
		cb.POST("/remove", batchRemove)
	}
}

func bind(c *bm.Context, v model.ParamValidator) (err error) {
	if err = c.Bind(v); err != nil {
		err = errors.WithStack(err)
		return
	}
	if !v.Validate() {
		err = ecode.RequestErr
		c.JSON(nil, ecode.RequestErr)
		return
	}
	return
}
