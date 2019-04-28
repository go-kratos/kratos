package http

import (
	"go-common/app/service/main/archive/conf"
	"go-common/app/service/main/archive/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	idfSvc *verify.Verify
	arcSvc *service.Service
)

// Init init http router.
func Init(c *conf.Config, s *service.Service) {
	arcSvc = s
	idfSvc = verify.New(nil)
	// init internal router
	en := bm.DefaultServer(c.BM.Inner)
	innerRouter(en)
	// init internal server
	if err := en.Start(); err != nil {
		log.Error("xhttp.Serve error(%v)", err)
		panic(err)
	}
	// init external router
	enlocal := bm.DefaultServer(c.BM.Local)
	localRouter(enlocal)
	// init external server
	if err := enlocal.Start(); err != nil {
		log.Error("xhttp.Serve error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	archive := e.Group("/x/internal/v2/archive", bm.CORS())
	{
		archive.GET("", idfSvc.Verify, arcInfo)
		archive.GET("/view", idfSvc.Verify, arcView)
		archive.GET("/views", idfSvc.Verify, arcViews)
		archive.GET("/page", idfSvc.Verify, arcPage)
		archive.GET("/video", idfSvc.Verify, video)
		archive.GET("/archives", idfSvc.Verify, archives)
		archive.GET("/archives/playurl", idfSvc.Verify, archivesWithPlayer)
		archive.GET("/typelist", idfSvc.Verify, typelist)
		archive.GET("/description", idfSvc.Verify, description)
		archive.GET("/maxAid", idfSvc.Verify, maxAID)
		regionGp := archive.Group("/region")
		{
			regionGp.GET("", idfSvc.Verify, regionArcs)
		}
		videoshotGp := archive.Group("/videoshot")
		{
			videoshotGp.GET("", idfSvc.Verify, videoshot)
		}
		statGp := archive.Group("/stat")
		{
			statGp.GET("", idfSvc.Verify, arcStat)
			statGp.GET("/stats", idfSvc.Verify, arcStats)
		}
		upGp := archive.Group("/up")
		{
			upGp.GET("/count/batch", idfSvc.Verify, uppersCount)
			upGp.GET("/count", idfSvc.Verify, upperCount)
			upGp.GET("/passed", idfSvc.Verify, upperPassed)
			upGp.GET("/cache", idfSvc.Verify, upperCache)
		}
	}
}

// localRouter init local router.
func localRouter(e *bm.Engine) {
	e.GET("/archive-service/rank/init", addRegionArc)
}
