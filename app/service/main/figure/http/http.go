package http

import (
	"net/http"

	"go-common/app/service/main/figure/conf"
	"go-common/app/service/main/figure/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	svc    *service.Service
	verSvc *verify.Verify
)

// Init new a http server
func Init(s *service.Service) {
	initService(s)
	// init router.
	e := bm.DefaultServer(conf.Conf.BM)
	internalRouter(e)
	if err := e.Start(); err != nil {
		log.Error("e.Start() error(%v)", err)
		panic(err)
	}
}

func initService(s *service.Service) {
	svc = s
	verSvc = verify.New(conf.Conf.Verify)
}

// internalRouter init inner router.
func internalRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)

	group := e.Group("/x/internal/figure", verSvc.Verify)
	{
		group.GET("/info", figureInfo)
		group.POST("/infos", figureInfos)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("figure-service service ping error (%+v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// register check server ok.
func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
