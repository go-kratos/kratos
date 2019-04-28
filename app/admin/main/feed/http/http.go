package http

import (
	"go-common/app/admin/main/feed/conf"
	bfssvr "go-common/app/admin/main/feed/service/bfs"
	"go-common/app/admin/main/feed/service/channel"
	"go-common/app/admin/main/feed/service/common"
	"go-common/app/admin/main/feed/service/egg"
	pgcsvr "go-common/app/admin/main/feed/service/pgc"
	"go-common/app/admin/main/feed/service/popular"
	"go-common/app/admin/main/feed/service/search"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	authSvc    *permit.Permit
	eggSvc     *egg.Service
	bfsSvc     *bfssvr.Service
	searchSvc  *search.Service
	pgcSvr     *pgcsvr.Service
	chanelSvc  *channel.Service
	popularSvc *popular.Service
	cardSvc    *channel.Service
	commonSvc  *common.Service
)

// initService init service
func initService(c *conf.Config) {
	authSvc = permit.New(c.Auth)
	eggSvc = egg.New(c)
	bfsSvc = bfssvr.New(c)
	searchSvc = search.New(c)
	pgcSvr = pgcsvr.New(c)
	chanelSvc = channel.New(c)
	cardSvc = channel.New(c)
	popularSvc = popular.New(c)
	commonSvc = common.New(c)
}

// Init init http sever instance.
func Init(c *conf.Config) {
	initService(c)
	engine := bm.DefaultServer(c.HTTPServer)
	innerRouter(engine)
	// init internal server
	if err := engine.Start(); err != nil {
		log.Error("httpx.Serve error(%v)", err)
		panic(err)
	}
}

// innerRouter
func innerRouter(e *bm.Engine) {
	// ping monitor
	e.GET("/monitor/ping", ping)
	//modules color eggs
	feed := e.Group("/x/admin/feed")
	{
		feed.POST("/upload", clientUpload)
		//对外 搜索
		feed.GET("/eggSearch", searchEgg)
		//对外 web
		feed.GET("/eggSearchWeb", SearchEggWeb)
		common := feed.Group("/common")
		{
			common.GET("/card/titlePreview", cardPreview2)
			common.GET("/log/action", actionLog)
			common.GET("/pgc/season", getPgcSeason)
			common.GET("/pgc/seasons", getPgcSeasons)
			common.GET("/pgc/ep", getPgcEp)
			common.GET("/card/type", cardType)
		}
		egg := feed.Group("/egg")
		{
			egg.POST("/add", addEgg)
			egg.GET("/index", indexEgg)
			egg.POST("/update", updateEgg)
			egg.POST("/publish", pubEgg)
			egg.POST("/delete", delEgg)
		}
		//对外
		open := feed.Group("/open")
		{
			//search
			open.POST("/search/addHotword", openAddHotword)   //搜索 添加热词
			open.POST("/search/addDarkword", openAddDarkword) //搜索 添加黑马词
			open.GET("/search/blackList", openBlacklist)      //搜索 黑名单
			open.GET("/search/hotwords", openHotList)         //搜索 热词
			open.GET("/search/darkword", openDarkword)        //搜索 获取黑马词
			open.GET("/search/webSearch", openSearchWeb)      //web 搜索
			open.POST("/ai/addPopStars", aiAddPopularStars)   //AI 添加新星卡片
		}
		search := feed.Group("/search", authSvc.Permit("SEARCH_HOTWORD"))
		{
			search.GET("/blackList", blackList)
			search.POST("/addBlack", addBlack)
			search.POST("/delBlack", delBlack)
			search.GET("/hot", HotList)
			search.POST("/addInter", addInter)
			search.POST("/updateInter", updateInter)
			search.POST("/deleteHot", deleteHot)
			search.POST("/updateSearch", updateSearch)
			search.POST("/publishHot", publishHotWord)
			search.POST("/publishDark", publishDarkWord)
			search.GET("/dark", darkList)
			search.POST("/delDark", deleteDark)
		}
		searchWeb := feed.Group("/search/web")
		{
			searchWeb.GET("/card/list", searchWebCardList)
			searchWeb.POST("/card/add", addSearchWebCard)
			searchWeb.POST("/card/update", upSearchWebCard)
			searchWeb.POST("/card/delete", delSearchWebCard)
			searchWeb.GET("/list", searchWebList)
			searchWeb.POST("/add", addSearchWeb)
			searchWeb.POST("/update", upSearchWeb)
			searchWeb.POST("/delete", delSearchWeb)
			searchWeb.POST("/opt", optSearchWeb)
		}
		cardsetup := feed.Group("/cardsetup")
		{
			cardsetup.POST("/add", addCardSetup)
			cardsetup.GET("/list", cardSetupList)
			cardsetup.POST("/delete", delCardSetup)
			cardsetup.POST("/update", updateCardSetup)
		}
		channel := feed.Group("/channel")
		{
			tab := channel.Group("/tab")
			{
				tab.GET("/list", tabList)
				tab.POST("/add", addTab)
				tab.POST("/update", updateTab)
				tab.POST("/delete", deleteTab)
				tab.POST("/offline", offlineTab)
			}
		}
		popular := feed.Group("/popular")
		{
			eventTopic := popular.Group("/event_topic")
			{
				eventTopic.GET("/list", eventTopicList)
				eventTopic.POST("/add", addEventTopic)
				eventTopic.POST("/update", upEventTopic)
				eventTopic.GET("/delete", delEventTopic)
			}
			stars := popular.Group("/stars")
			{
				stars.GET("/list", popularStarsList)
				stars.POST("/add", addPopularStars)
				stars.POST("/update", updatePopularStars)
				stars.POST("/delete", deletePopularStars)
				stars.POST("/reject", rejectPopularStars)
			}
		}
	}
}

// ping check server ok.
func ping(c *bm.Context) {

}
