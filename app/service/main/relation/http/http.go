package http

import (
	"net/http"

	"go-common/app/service/main/relation/conf"
	"go-common/app/service/main/relation/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/rate"
	v "go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/metadata"
)

var (
	anti             *antispam.Antispam
	relationSvc      *service.Service
	verify           *v.Verify
	addFollowingRate *rate.Limiter
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	relationSvc = s
	verify = v.New(c.Verify)
	anti = antispam.New(c.Antispam)
	addFollowingRate = rate.New(c.AddFollowingRate)
	// init inner router
	engine := bm.DefaultServer(c.BM)
	setupInnerEngine(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start() error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router.
func setupInnerEngine(e *bm.Engine) {
	// health check
	e.Ping(ping)
	e.Register(register)
	//new defined api lists
	g := e.Group("/x/internal/relation", verify.Verify)
	// relation
	g.GET("", relation)
	g.GET("/relations", relations)
	// stat
	g.GET("/stat", stat)
	g.GET("/stats", stats)
	// private api
	g.POST("/stat/set", setStat)
	g.POST("/tag/cache/del", delTagCache)
	g.GET("/tag/special", special)

	// following group
	following := g.Group("/following")
	following.GET("/followings", followings)
	following.GET("/same/followings", sameFollowings)
	following.POST("/add", addFollowingRate.Handler(), anti.ServeHTTP, addFollowing)
	following.POST("/del", anti.ServeHTTP, delFollowing)

	// whisper group
	whisper := g.Group("/whisper")
	whisper.GET("/whispers", whispers)
	whisper.POST("/add", anti.ServeHTTP, addWhisper)
	whisper.POST("/del", anti.ServeHTTP, delWhisper)

	// black group
	black := g.Group("/black")
	black.GET("/blacks", blacks)
	black.POST("/add", anti.ServeHTTP, addBlack)
	black.POST("/del", anti.ServeHTTP, delBlack)

	// follower group
	follower := g.Group("/follower")
	follower.GET("/followers", followers)
	follower.POST("/del", anti.ServeHTTP, delFollower)

	// recommend group
	// recommend := g.Group("/recommend")
	// recommend.GET("/global/hot", globalHot)

	// cache group
	cache := g.Group("/cache")
	cache.POST("/following/del", delFollowingCache)
	cache.POST("/following/update", updateFollowingCache)
	cache.POST("/follower/del", delFollowerCache)
	cache.POST("/stat/del", delStatCache)

	// cache group
	admin := g.Group("/admin")
	admin.POST("/monitor/add", addMonitor)
	admin.POST("/monitor/del", delMonitor)
	admin.GET("/monitor/load", loadMonitor)
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := relationSvc.Ping(c); err != nil {
		log.Error("service ping error(%v)", err)
		c.Writer.WriteHeader(http.StatusServiceUnavailable)
	}
}

// register check server ok.
func register(c *bm.Context) {
	c.JSON(map[string]struct{}{}, nil)
}

func infocArg(c *bm.Context) (arg map[string]string) {
	var (
		ua      string
		referer string
		sid     string
		req     = c.Request
	)
	ua = req.Header.Get(service.RelInfocUA)
	referer = req.Header.Get(service.RelInfocReferer)
	sidCookie, err := req.Cookie(service.RelInfocSid)
	if err != nil {
		log.Warn("relation infoc get sid failed error(%v)", err)
	} else {
		sid = sidCookie.Value
	}
	buvid := req.Header.Get(service.RelInfocHeaderBuvid)
	if buvid == "" {
		buvidCookie, _ := req.Cookie(service.RelInfocCookieBuvid)
		if buvidCookie != nil {
			buvid = buvidCookie.Value
		}
	}
	arg = map[string]string{
		service.RelInfocIP:      metadata.String(c, metadata.RemoteIP),
		service.RelInfocReferer: referer,
		service.RelInfocSid:     sid,
		service.RelInfocBuvid:   buvid,
		service.RelInfocUA:      ua,
	}
	return
}
