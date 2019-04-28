package http

import (
	"net/http"

	"go-common/app/interface/main/web/conf"
	"go-common/app/interface/main/web/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/cache"
	"go-common/library/net/http/blademaster/middleware/cache/store"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	webSvc  *service.Service
	authSvr *auth.Auth
	vfySvr  *verify.Verify

	// cache components
	cacheSvr *cache.Cache
	deg      *cache.Degrader
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	authSvr = auth.New(c.Auth)
	vfySvr = verify.New(c.Verify)
	webSvc = s
	cacheSvr = cache.New(store.NewMemcache(c.DegradeConfig.Memcache))
	deg = cache.NewDegrader(c.DegradeConfig.Expire)
	// init outer router
	engine := bm.NewServer(c.HTTPServer)
	engine.Use(bm.Recovery(), bm.Logger(), bm.Trace(), bm.Mobile())
	outerRouter(engine)
	internalRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.GET("/x/web-interface/view", authSvr.Guest, view)
	group := e.Group("/x/web-interface", bm.CSRF(), bm.CORS())
	{
		arcGroup := group.Group("/archive")
		{
			arcGroup.GET("/coins", authSvr.User, coins)
			arcGroup.GET("/stat", archiveStat)
			arcGroup.GET("/desc", description)
			arcGroup.POST("/report", authSvr.User, arcReport)
			arcGroup.POST("/appeal", authSvr.User, arcAppeal)
			arcGroup.GET("/appeal/tags", appealTags)
			arcGroup.GET("/author/recommend", authorRecommend)
			arcGroup.GET("/related", relatedArcs)
			arcGroup.POST("/like", authSvr.User, like)
			arcGroup.POST("/like/triple", authSvr.User, likeTriple)
			arcGroup.GET("/has/like", authSvr.User, hasLike)
			arcGroup.GET("/ugc/pay", authSvr.User, arcUGCPay)
			arcGroup.GET("/relation", authSvr.User, arcRelation)
		}
		dyGroup := group.Group("/dynamic")
		{
			dyGroup.GET("/region", dynamicRegion)
			dyGroup.GET("/index", dynamicRegions)
			dyGroup.GET("/tag", dynamicRegionTag)
			dyGroup.GET("/total", dynamicRegionTotal)
		}
		rankGroup := group.Group("/ranking")
		{
			rankGroup.GET("", ranking)
			rankGroup.GET("/index", rankingIndex)
			rankGroup.GET("/region", rankingRegion)
			rankGroup.GET("/recommend", rankingRecommend)
			rankGroup.GET("/tag", rankingTag)
		}
		tagGroup := group.Group("/tag")
		{
			tagGroup.GET("/top", tagAids)
		}
		artGroup := group.Group("/article")
		{
			artGroup.GET("/list", authSvr.Guest, articleList)
			artGroup.GET("/up/list", authSvr.Guest, articleUpList)
			artGroup.GET("/categories", categories)
			artGroup.GET("/newcount", newCount)
			artGroup.GET("/early", upMoreArts)
		}
		coinGroup := group.Group("/coin")
		{
			coinGroup.POST("/add", authSvr.User, addCoin)
			coinGroup.GET("/today/exp", authSvr.User, coinExp)
		}
		onlineGroup := group.Group("/online")
		{
			onlineGroup.GET("", onlineInfo)
			onlineGroup.GET("/list", onlineList)
		}
		helpGroup := group.Group("/help")
		{
			helpGroup.GET("/list", cacheSvr.Cache(deg.Args("parentTypeId"), nil), helpList)
			helpGroup.GET("/detail", cacheSvr.Cache(deg.Args("pn", "ps", "fId", "questionTypeId"), nil), helpDetail)
			helpGroup.GET("/search", helpSearch)
		}
		viewGroup := group.Group("/view")
		{
			viewGroup.GET("/detail", authSvr.Guest, detail)
		}
		searchGroup := group.Group("/search")
		{
			searchGroup.GET("/all", authSvr.Guest, searchAll)
			searchGroup.GET("/type", authSvr.Guest, searchByType)
			searchGroup.GET("/recommend", authSvr.Guest, searchRec)
			searchGroup.GET("/default", authSvr.Guest, searchDefault)
			searchGroup.GET("/egg", searchEgg)
		}
		wxGroup := group.Group("/wx")
		{
			wxGroup.GET("/hot", wxHot)
			wxGroup.GET("/search/all", authSvr.Guest, wxSearchAll)
		}
		bnjGroup := group.Group("/bnj2019")
		{
			bnjGroup.GET("", authSvr.Guest, bnj2019)
			bnjGroup.GET("/timeline", authSvr.Guest, timeline)
		}
		group.GET("/region/custom", regionCustom)
		group.GET("/attentions", authSvr.User, attentions)
		group.GET("/card", authSvr.Guest, card)
		group.GET("/nav", authSvr.Guest, nav)
		group.GET("/newlist", newList)
		group.POST("/feedback", authSvr.Guest, feedback)
		group.GET("/zone", ipZone)
		group.POST("/share/add", authSvr.Guest, addShare)
		group.GET("/elec/show", authSvr.Guest, elecShow)
		group.GET("/index/icon", indexIcon)
		group.GET("/baidu/kv", kv)
		group.GET("/cmtbox", cmtbox)
		group.GET("/abserver", authSvr.Guest, abServer)
		group.GET("/up/rec", authSvr.User, upRec)
		group.GET("/broadcast/servers", broadServer)
	}
	e.GET("/x/coin/list", coinList)
}

func internalRouter(e *bm.Engine) {
	group := e.Group("/x/internal/web-interface")
	{
		dyGroup := group.Group("/dynamic")
		{
			dyGroup.GET("/region", vfySvr.Verify, dynamicRegion)
			dyGroup.GET("/index", vfySvr.Verify, dynamicRegions)
			dyGroup.GET("/tag", vfySvr.Verify, dynamicRegionTag)
			dyGroup.GET("/total", vfySvr.Verify, dynamicRegionTotal)
		}
		rankGroup := group.Group("/ranking")
		{
			rankGroup.GET("", vfySvr.Verify, ranking)
			rankGroup.GET("/index", vfySvr.Verify, rankingIndex)
			rankGroup.GET("/region", vfySvr.Verify, rankingRegion)
			rankGroup.GET("/recommend", vfySvr.Verify, rankingRecommend)
			rankGroup.GET("/tag", vfySvr.Verify, rankingTag)
		}
		tagGroup := group.Group("/tag")
		{
			tagGroup.GET("/top", vfySvr.Verify, tagAids)
			tagGroup.GET("/detail", vfySvr.Verify, tagDetail)
		}
		helpGroup := group.Group("/help")
		{
			helpGroup.GET("/list", vfySvr.Verify, helpList)
			helpGroup.GET("/detail", vfySvr.Verify, helpDetail)
			helpGroup.GET("/search", vfySvr.Verify, helpSearch)
		}
		onlineGroup := group.Group("/online")
		{
			onlineGroup.GET("", vfySvr.Verify, onlineInfo)
			onlineGroup.GET("/list", vfySvr.Verify, onlineList)
		}
		viewGroup := group.Group("/view")
		{
			viewGroup.GET("", vfySvr.Verify, authSvr.Guest, view)
			viewGroup.GET("/detail", vfySvr.Verify, authSvr.Guest, detail)
		}
		searchGroup := group.Group("/search")
		{
			searchGroup.GET("/all", vfySvr.Verify, authSvr.Guest, searchAll)
			searchGroup.GET("/type", vfySvr.Verify, authSvr.Guest, searchByType)
			searchGroup.GET("/recommend", vfySvr.Verify, authSvr.Guest, searchRec)
		}
		group.GET("/newlist", vfySvr.Verify, newList)
		group.GET("/zone", vfySvr.Verify, ipZone)
		group.GET("/region/custom", vfySvr.Verify, regionCustom)
		group.GET("/baidu/kv", vfySvr.Verify, kv)
		group.GET("/cmtbox", vfySvr.Verify, cmtbox)
		group.GET("/broadcast/servers", vfySvr.Verify, broadServer)
		group.GET("/bnj2019", vfySvr.Verify, authSvr.Guest, bnj2019)
		group.GET("/bnj2019/aids", vfySvr.Verify, bnj2019Aids)
	}
}

func ping(c *bm.Context) {
	if err := webSvc.Ping(c); err != nil {
		log.Error("web-interface  ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
