package http

import (
	"net/http"

	"go-common/app/interface/main/upload/conf"
	xanti "go-common/app/interface/main/upload/http/antispam"
	"go-common/app/interface/main/upload/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	uploadSvr    *service.Service
	authInterSvr *auth.Auth
	authSvr      *auth.Auth
	verifySvr    *verify.Verify
	anti         *xanti.Antispam
)

// Init init http
func Init(c *conf.Config, s *service.Service) {
	initService(c, s)
	engine := bm.DefaultServer(c.BM)
	initRouter(engine)
	// init Outer serve
	if err := engine.Start(); err != nil {
		log.Error("xhttp.Serve error(%v)", err)
		panic(err)
	}
}

func initService(c *conf.Config, s *service.Service) {
	uploadSvr = s
	authInterSvr = auth.New(c.AuthInter)
	authSvr = auth.New(c.AuthOut)
	verifySvr = verify.New(nil)
	anti = xanti.New(c.Antispam, s.GetRateLimit) //mid+dir 限流
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	uploadInternal := e.Group("/x/internal")
	{
		uploadInternal.POST("/upload", verifySvr.Verify, internalUpload)
		uploadInternal.POST("/upload/image", authInterSvr.User, anti.Handler(), internalUploadImage)
		uploadInternal.POST("/upload/admin/image", verifySvr.Verify, anti.Handler(), internalUploadAdminImage)
		uploadInternal.POST("/image/gen", verifySvr.Verify, anti.Handler(), genImageUpload)
	}

	upload := e.Group("/x/upload")
	{
		upload.POST("/image", uploadImagePublic)
		upload.POST("/app/image", authSvr.UserMobile, anti.Handler(), uploadMobileImage)
		upload.POST("/web/image", authSvr.UserWeb, anti.Handler(), uploadWebImage)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	var err error
	if err = uploadSvr.Ping(c); err != nil {
		log.Error("upload service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
