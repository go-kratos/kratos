package http

import (
	"go-common/app/interface/bbq/wechat/internal/conf"
	"go-common/app/interface/bbq/wechat/internal/model"
	"go-common/app/interface/bbq/wechat/internal/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"net/http"

	"github.com/pkg/errors"
)

var (
	svc *service.Service
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	svc = service.New(c)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/bbq/wechat")
	{
		g.GET("/token", getToken)
	}
}

func getToken(c *bm.Context) {
	arg := new(model.TokenReq)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	url := c.Request.Referer()

	c.JSON(svc.TokenGet(c, arg, url))
}
