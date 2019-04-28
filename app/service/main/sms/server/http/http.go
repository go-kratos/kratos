package http

import (
	"net/http"

	"go-common/app/service/main/sms/conf"
	"go-common/app/service/main/sms/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	smsSvc *service.Service
	idfSvc *verify.Verify
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	idfSvc = verify.New(c.Verify)
	smsSvc = s
	engine := bm.DefaultServer(c.HTTPServer)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/x/internal/sms", bm.CORS())
	{
		g.POST("/send", idfSvc.Verify, send)
		g.POST("/sendBatch", idfSvc.Verify, sendBatch)
	}
}

func ping(c *bm.Context) {
	if err := smsSvc.Ping(c); err != nil {
		log.Error("sms-service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
