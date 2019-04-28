package http

import (
	"github.com/pkg/errors"
	"net/http"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/net/http/blademaster/middleware/verify"

	"go-common/app/service/bbq/sys-msg/api/v1"
	"go-common/app/service/bbq/sys-msg/internal/conf"
	"go-common/app/service/bbq/sys-msg/internal/service"
)

var (
	vfy *verify.Verify
	svc *service.Service
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	svc = s
	vfy = verify.New(c.Verify)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/x/sys-msg")
	{
		g.GET("/msg/list", listSysMsg)
		g.POST("/msg/create", createSysMsg)
		g.POST("/msg/update", updateSysMsg)
	}
}

func ping(ctx *bm.Context) {
	if err := svc.Ping(ctx); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

func createSysMsg(c *bm.Context) {
	arg := &v1.SysMsg{}
	if err := c.BindWith(arg, binding.JSON); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	resp, err := svc.CreateSysMsg(c, arg)
	c.JSON(resp, err)
}

func updateSysMsg(c *bm.Context) {
	arg := &v1.SysMsg{}
	if err := c.BindWith(arg, binding.JSON); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	resp, err := svc.UpdateSysMsg(c, arg)
	c.JSON(resp, err)
}

func listSysMsg(c *bm.Context) {
	arg := &v1.ListSysMsgReq{}
	if err := c.BindWith(arg, binding.Query); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	resp, err := svc.ListSysMsg(c, arg)
	c.JSON(resp, err)
}
