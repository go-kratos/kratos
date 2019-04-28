package http

import (
	"net/http"

	"go-common/app/service/main/point/conf"
	"go-common/app/service/main/point/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	svc          *service.Service
	verSvc       *verify.Verify
	authSvc      *auth.Auth
	whiteAppkeys string
)

// Init init.
func Init(s *service.Service) {
	initService(s)
	// init router.
	engineInner := bm.DefaultServer(conf.Conf.BM)
	innerRouter(engineInner)
	outerRouter(engineInner)
	if err := engineInner.Start(); err != nil {
		log.Error("engineInner.Start() error(%v)", err)
		panic(err)
	}
	whiteAppkeys = conf.Conf.Property.PointWhiteAppkeys
}

func initService(s *service.Service) {
	svc = s
	verSvc = verify.New(conf.Conf.Verify)
	authSvc = auth.New(conf.Conf.Auth)
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("point http service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// register check server ok.
func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
