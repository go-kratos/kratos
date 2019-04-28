package http

import (
	"net/http"

	"go-common/app/interface/openplatform/article/conf"
	"go-common/app/interface/openplatform/article/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/cache"
	"go-common/library/net/http/blademaster/middleware/cache/store"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	artSrv    *service.Service
	authSvr   *auth.Auth
	vfySvr    *verify.Verify
	antispamM *antispam.Antispam
	// cache components
	cacheSvr *cache.Cache
	deg      *cache.Degrader
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	authSvr = auth.New(c.Auth)
	vfySvr = verify.New(c.Verify)
	artSrv = s
	antispamM = antispam.New(c.Antispam)
	cacheSvr = cache.New(store.NewMemcache(c.DegradeConfig.Memcache))
	deg = cache.NewDegrader(c.DegradeConfig.Expire)
	// init outer router
	engine := bm.DefaultServer(c.BM)
	outerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// outerRouter init outer router
func outerRouter(r *bm.Engine) {
	r.Ping(ping)
	r.Register(register)
	cr := r.Group("/x/article")
	{
		cr.GET("/recommends", authSvr.Guest, recommends)
		cr.GET("/recommends/plus", authSvr.Guest, recommendsPlus)
		cr.GET("/home", authSvr.Guest, cacheSvr.Cache(deg.Args("pn", "ps", "device", "mobi_app", "build"), nil), home)
		cr.GET("/view", authSvr.Guest, view)
		cr.GET("/metas", metas)
		cr.GET("/card", card)
		cr.GET("/cards", cards)
		cr.GET("/notice", notice)
		cr.GET("/addview", authSvr.Guest, addView)
		cr.POST("/addshare", authSvr.Guest, addShare)
		cr.GET("/viewinfo", authSvr.Guest, viewInfo)
		cr.GET("/actinfo", actInfo)
		cr.POST("/like", authSvr.User, like)
		cr.GET("/applyinfo", authSvr.Guest, applyInfo)
		cr.GET("/is_author", authSvr.User, isAuthor)
		cr.POST("/author/add", authSvr.User, addAuthor)
		cr.POST("/apply", authSvr.User, apply)
		cr.POST("/complaints", authSvr.User, addComplaint)
		cr.GET("/list", list)
		cr.GET("/categories", categories)
		cr.GET("/anniversary", authSvr.User, anniversaryInfo)
		cr.GET("/sentinel/config", authSvr.Guest, sentinel)
		ccr := cr.Group("/favorites", authSvr.User)
		{
			ccr.POST("/add", addFavorite)
			ccr.POST("/del", delFavorite)
			ccr.GET("/list", favorites)
			ccr.GET("/list/all", allFavorites)
		}
		cr.GET("/archives", archives)
		cr.GET("/early", earlyArticles)
		cr.GET("/more", authSvr.Guest, moreArts)
		ccr = cr.Group("/rank")
		{
			ccr.GET("/categories", rankCategories)
			ccr.GET("/list", authSvr.Guest, ranks)
		}
		ccr = cr.Group("/user", authSvr.User)
		{
			ccr.GET("/notice", userNotice)
			ccr.POST("/notice/update", updateUserNotice)
		}
		// read list
		cr.GET("/list/articles", authSvr.Guest, listArticles)
		cr.GET("/list/web/articles", authSvr.Guest, webListArticles)
		cr.GET("/listinfo", listInfo)
		cr.GET("/up/lists", upLists)
		cr.GET("/hotspots", authSvr.Guest, hotspotArts)
		cr.GET("/authors", authSvr.User, authors)
		ccr = cr.Group("/creative", authSvr.User)
		{
			cr1 := ccr.Group("/list")
			{
				cr1.GET("/all", lists)
				cr1.POST("/add", addList)
				cr1.POST("/del", delList)
				cr1.POST("/update", updateListArticles)
				cr1.GET("/articles/all", listAllArticles)
				cr1.GET("/articles/can_add", canAddArts)
				cr1.POST("/articles/update", updateArticleList)
			}
			// creative article
			ccr.POST("/article/submit", webSubArticle)
			ccr.POST("/article/update", webUpdateArticle)
			ccr.POST("/draft/addupdate", webSubmitDraft)
			ccr.GET("/draft/view", webDraft)
			ccr.GET("/draft/list", webDraftList)
			ccr.GET("/draft/count", draftCount)
			ccr.GET("/article/view", webArticle)
			ccr.GET("/article/list", webArticleList)
			ccr.GET("/app/pre", creatorArticlePre)
			ccr.POST("/upload/image", antispamM.ServeHTTP, uploadImage)
			ccr.POST("/draft/delete", deleteDraft)
			ccr.POST("/article/delete", delArticle)
			ccr.POST("/article/withdraw", withdrawArticle)
			ccr.POST("/article/capture", antispamM.ServeHTTP, articleCapture)
			ccr.POST("/segment", segment)
		}
		// article read ping for timing
		cr.GET("/read/ping", authSvr.Guest, readPing)
	}

	cr = r.Group("/x/internal/article", vfySvr.Verify)
	{
		cr.GET("/meta", meta)
		cr.GET("/metas", metas)
		cr.GET("/list", list)
		cr.GET("/view", view)
		cr.GET("/recommends/all", allRecommends)
		cr.POST("/refresh_list", refreshLists)
		cr.POST("/rebuild_allrc", rebuildAllListReadCount)
		cr.POST("/lock", addCheatFilter)
		cr.POST("/unlock", delCheatFilter)
	}
}

func ping(c *bm.Context) {
	if err := artSrv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

func buvid(c *bm.Context) string {
	buvid := c.Request.Header.Get(_headerBuvid)
	if buvid == "" {
		cookie, _ := c.Request.Cookie(_buvid)
		if cookie != nil {
			buvid = cookie.Value
		}
	}
	return buvid
}
