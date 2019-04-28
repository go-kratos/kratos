package http

import (
	"go-common/app/service/main/member/conf"
	"go-common/app/service/main/member/server/http/block"
	"go-common/app/service/main/member/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	v "go-common/library/net/http/blademaster/middleware/verify"
	"net/http"
)

var (
	memberSvc *service.Service
	verify    *v.Verify
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	verify = v.New(c.Verify)
	memberSvc = s
	engine := bm.DefaultServer(c.BM)
	setup(engine)
	block.Setup(memberSvc.BlockImpl(), engine, verify)
	if err := engine.Start(); err != nil {
		log.Error("http.Serve inner error(%v)", err)
		panic(err)
	}
}

func setup(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)

	mb := e.Group("/x/internal/member", verify.Verify)
	mb.POST("/sign/update", setSign)
	mb.POST("/name/update", setName)
	mb.POST("/rank/update", setRank)
	mb.POST("/birthday/update", setBirthday)
	mb.POST("/sex/update", setSex)
	mb.POST("/face/update", setFace)
	mb.POST("/base/update", setBase)
	mb.POST("/morals/update", updateMorals)
	mb.POST("/moral/update", updateMoral)
	mb.POST("/moral/undo", undoMoral)
	mb.GET("/moral", moral)
	mb.GET("/moral/log", moralLog)
	mb.GET("", member)
	mb.GET("/base", base)
	mb.GET("/batchBase", batchBase)
	mb.GET("/exp", exp)
	mb.GET("/level", level)
	mb.GET("/official", official)
	mb.POST("/exp/set", setExp)
	mb.POST("/exp/update", updateExp)
	mb.GET("/exp/log", explog)
	mb.GET("/exp/stat", stat)
	mb.GET("/cache/del", cacheDel)
	mb.POST("/property/review/add", addPropertyReview)
	// realname
	mb.GET("/realname/status", realnameStatus)
	mb.GET("/realname/info", realnameInfo)
	mb.POST("/realname/tel/capture", realnameTelCapture)
	mb.GET("/realname/tel/capture/check", realnameCheckTelCapture)
	mb.GET("/realname/apply/status", realnameApplyStatus)
	mb.POST("/realname/apply", realnameApply)
	mb.GET("/realname/adult", realnameAdult)
	mb.GET("/realname/check", realnameCheck)
	mb.GET("/realname/stripped/info", realnameStrippedInfo)
	mb.GET("/realname/mid/by/card", realnameMidByCard)
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := memberSvc.Ping(c); err != nil {
		log.Error("service ping error(%+v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(nil, nil)
}
