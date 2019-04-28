package http

import (
	"net/http"

	"go-common/app/service/main/thumbup/conf"
	"go-common/app/service/main/thumbup/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/rate"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	likeSrv   *service.Service
	verifySrv *verify.Verify
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	verifySrv = verify.New(c.Verify)
	rateLimit := rate.New(c.Rate)
	likeSrv = s
	// init outer router
	engineOuter := bm.DefaultServer(c.BM)
	engineOuter.Use(rateLimit)
	outerRouter(engineOuter)
	if err := engineOuter.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// outerRouter init outer router
func outerRouter(r *bm.Engine) {
	r.Ping(ping)
	r.Register(register)
	cr := r.Group("/x/internal/thumbup", verifySrv.Verify)
	{
		cr.GET("/stats", stats)
		cr.POST("/multi_stats", multiStats)
		cr.GET("/user_likes", userLikes)
		cr.GET("/item_likes", itemLikes)
		cr.POST("/like", like)
		cr.GET("/has_like", hasLike)
		cr.POST("/update_count", updateCount)
		cr.GET("/raw_stats", rawStats)
		cr.POST("/update_upmids", updateUpMids)
		cr.POST("/item_has_like", itemHasLike)
	}
}

func ping(c *bm.Context) {
	if err := likeSrv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
