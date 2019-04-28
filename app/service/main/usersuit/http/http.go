package http

import (
	"net/http"

	"go-common/app/service/main/usersuit/conf"
	"go-common/app/service/main/usersuit/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	usersuitSvc *service.Service
	vfySvc      *verify.Verify
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	usersuitSvc = s
	vfySvc = verify.New(c.Verify)
	// init inner router
	innerEngine := bm.DefaultServer(c.BM)
	innerRouter(innerEngine)
	// init inner server
	if err := innerEngine.Start(); err != nil {
		log.Error("xhttp.Serve inner error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	// health check
	e.Ping(ping)
	e.Register(register)
	group := e.Group("/x/internal/pendant", vfySvc.Verify)
	{
		group.GET("/groupInfo", groupInfo)
		group.GET("/groupInfoByID", groupInfoByID)
		group.GET("/pendantByID", pendantByID)
		group.GET("/vipGroup", vipGroup)
		group.GET("/entryGroup", entryGroup)
		group.GET("/pointRecommend", pointRecommend)
		group.GET("/package", packageInfo)
		group.GET("/equipment", equipment)
		group.GET("/orderHistory", orderHistory)
		group.POST("/order", order)
		group.POST("/equip", equip)
		group.POST("/multiGrantByMid", multiGrantByMid)
		group.POST("/multiGrantByPid", multiGrantByPid)
	}
	e.POST("/x/internal/pendant/callback", pendantCallback)
	medal := e.Group("/x/internal/medal", vfySvc.Verify)
	{
		medal.GET("/my", medalMy)
		medal.GET("/all", medalAllInfo)
		medal.GET("/info", medalInfo)
		medal.GET("/popup", medalPopup)
		medal.GET("/user", medalUser)
		medal.GET("/check", medalCheck)
		medal.GET("/activated", medalActivated)
		medal.POST("/install", medalInstall)
		medal.POST("/grant", medalGet)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := usersuitSvc.Ping(c); err != nil {
		log.Error("usersuit-service service ping error (%+v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// register check server ok.
func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
