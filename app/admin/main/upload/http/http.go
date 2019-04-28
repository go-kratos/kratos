package http

import (
	"go-common/app/admin/main/upload/conf"
	"go-common/app/admin/main/upload/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	authSvc   *permit.Permit
	uaSvc     *service.Service
	verifySvc *verify.Verify
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	initService(c, s)
	// init router
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config, s *service.Service) {
	authSvc = permit.New(c.Auth)
	verifySvc = verify.New(nil)
	uaSvc = s
}

// innerRouter init outer router api path.
func innerRouter(e *bm.Engine) {
	//init api
	e.Ping(ping)

	uploadAdmin := e.Group("/x/admin/upload")
	{
		uploadAdmin.POST("/add", add)
		uploadAdmin.GET("/list", authSvc.Permit(""), list)
		uploadAdmin.DELETE("/delete", authSvc.Permit(""), deleteFile)

		file := uploadAdmin.Group("/file")
		{
			file.POST("/upload", authSvc.Permit(""), InternalUploadAdminImage)
			file.DELETE("/delete", authSvc.Permit(""), deleteRawFile)
		}
	}

	uploadAdminV2 := e.Group("/x/admin/upload/v2")
	{
		uploadAdminV2.GET("/list", authSvc.Permit(""), multiList)
		uploadAdminV2.DELETE("/delete", authSvc.Permit(""), deleteFileV2)
	}

	bucket := uploadAdmin.Group("/bucket")
	{
		bucket.POST("/add", verifySvc.Verify, addBucket)
		bucket.GET("/list", listBucket)
		bucket.GET("/list/public", listPublicBucket)
		bucket.GET("/detail", detailBucket)

		dir := bucket.Group("/dir")
		{
			dir.POST("/add", verifySvc.Verify, addDir)
		}
	}
}
