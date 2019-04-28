package http

import (
	"net/http"

	"github.com/pkg/errors"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"

	"go-common/app/service/bbq/notice-service/api/v1"
	"go-common/app/service/bbq/notice-service/internal/conf"
	"go-common/app/service/bbq/notice-service/internal/service"
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
	g := e.Group("/x/notice-service")
	{
		g.GET("/notice/list", listNotices)
		g.GET("/notice/unread", unreadInfo)
		g.POST("/notice/create", createNotice)

		g.POST("/push/login", login)
		g.POST("/push/logout", logout)
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

func createNotice(c *bm.Context) {
	arg := &v1.NoticeBase{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	resp, err := svc.CreateNotice(c, arg)
	c.JSON(resp, err)
}

func listNotices(c *bm.Context) {
	arg := &v1.ListNoticesReq{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	resp, err := svc.ListNotices(c, arg)
	c.JSON(resp, err)
}

func unreadInfo(c *bm.Context) {
	arg := &v1.GetUnreadInfoRequest{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	resp, err := svc.GetUnreadInfo(c, arg)
	c.JSON(resp, err)
}

func login(c *bm.Context) {
	arg := &v1.UserPushDev{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	resp, err := svc.PushLogin(c, arg)
	c.JSON(resp, err)
}

func logout(c *bm.Context) {
	arg := &v1.UserPushDev{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	resp, err := svc.PushLogout(c, arg)
	c.JSON(resp, err)
}
