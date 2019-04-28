package http

import (
	"net/http"

	"go-common/app/admin/main/card/conf"
	"go-common/app/admin/main/card/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

const _maxnamelen = 18

var (
	srv       *service.Service
	cf        *conf.Config
	permitSvc *permit.Permit
)

// Init init
func Init(c *conf.Config, svc *service.Service) {
	srv = svc
	permitSvc = permit.New2(nil)
	cf = c
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
	group := e.Group("/x/admin/card", permitSvc.Permit2("CARD"))
	group.GET("/group/list", groups)
	group.POST("/group/add", addGroup)
	group.POST("/group/modify", updateGroup)
	group.POST("/group/state/change", groupStateChange)
	group.POST("/group/delete", deleteGroup)
	group.POST("/group/order", groupOrderChange)
	group.POST("/order", cardOrderChange)
	group.POST("/add", addCard)
	group.POST("/modify", updateCard)
	group.GET("/list", cards)
	group.POST("/state/change", cardStateChange)
	group.POST("/delete", deleteCard)
}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
