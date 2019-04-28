package http

import (
	"net/http"

	"go-common/app/interface/main/activity/conf"
	"go-common/app/interface/main/activity/service/bnj"
	"go-common/app/interface/main/activity/service/bws"
	"go-common/app/interface/main/activity/service/kfc"
	"go-common/app/interface/main/activity/service/like"
	"go-common/app/interface/main/activity/service/sports"
	"go-common/app/interface/main/activity/service/timemachine"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	likeSvc   *like.Service
	sportsSvc *sports.Service
	matchSvc  *like.Service
	bwsSvc    *bws.Service
	tmSvc     *timemachine.Service
	bnjSvc    *bnj.Service
	kfcSvc    *kfc.Service
	authSvc   *auth.Auth
	vfySvc    *verify.Verify
)

// Init int http service
func Init(c *conf.Config) {
	initService(c)
	engine := bm.DefaultServer(c.HTTPServer)
	outerRouter(engine)
	internalRouter(engine)
	// init Outer serve
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// initService init service
func initService(c *conf.Config) {
	authSvc = auth.New(c.Auth)
	vfySvc = verify.New(c.Verify)
	likeSvc = like.New(c)
	sportsSvc = sports.New(c)
	matchSvc = like.New(c)
	bwsSvc = bws.New(c)
	tmSvc = timemachine.New(c)
	bnjSvc = bnj.New(c)
	kfcSvc = kfc.New(c)
}

//CloseService close all service
func CloseService() {
	likeSvc.Close()
	bnjSvc.Close()
	kfcSvc.Close()
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	group := e.Group("/x/activity", bm.CORS())
	{
		group.GET("/subject", subject)
		group.POST("/vote", authSvc.User, vote)
		group.GET("/ltime", ltime)
		group.GET("/object/group", groupData)
		group.GET("/tag/object", tagList)
		group.GET("/tag/object/stats", tagStats)
		group.GET("/region/object", regionList)
		group.GET("/subject/total/stat", subjectStat)
		group.GET("/view/rank", viewRank)
		group.POST("/likeact", authSvc.User, likeAct)
		group.GET("/likeact/list", authSvc.Guest, likeActList)
		group.POST("/missiongroup/like", authSvc.User, missionLike)
		group.POST("/missiongroup/likeact", authSvc.User, missionLikeAct)
		group.GET("/missiongroup/info", authSvc.User, missionInfo)
		group.GET("/missiongroup/tops", missionTops)
		group.GET("/missiongroup/user", missionUser)
		group.GET("/missiongroup/rank", authSvc.User, missionRank)
		group.GET("/missiongroup/friends", authSvc.User, missionFriends)
		group.GET("/missiongroup/award", authSvc.User, missionAward)
		group.POST("/missiongroup/achievement", authSvc.User, missionAchieve)
		group.POST("/up/act", authSvc.User, storyKingAct)
		group.GET("/up/left", authSvc.User, storyKingLeft)
		group.GET("/up/list", authSvc.Guest, upList)
		spGroup := group.Group("/sports")
		{
			spGroup.GET("/qq", qq)
			spGroup.GET("/news", news)
		}
		matchGroup := group.Group("/match")
		{
			matchGroup.GET("", matchs)
			matchGroup.GET("/unstart", authSvc.Guest, unStart)
			matchGroup.POST("/cache/clear", clearCache)
			guGroup := matchGroup.Group("/guess")
			{
				guGroup.GET("", authSvc.User, guess)
				guGroup.GET("/list", authSvc.User, listGuess)
				guGroup.POST("/add", authSvc.User, addGuess)
			}
			foGroup := matchGroup.Group("/follow")
			{
				foGroup.GET("", authSvc.User, follow)
				foGroup.POST("/add", authSvc.User, addFollow)
			}
		}
		tmGroup := group.Group("/timemachine")
		{
			tmGroup.GET("/2018", authSvc.User, timemachine2018)
			tmGroup.GET("/2018/raw", authSvc.User, timemachine2018Raw)
			tmGroup.GET("/2018/cache", authSvc.User, timemachine2018Cache)
		}
		bwsGroup := group.Group("/bws")
		{
			bwsGroup.GET("/user", authSvc.Guest, user)
			bwsGroup.GET("/points", points)
			bwsGroup.GET("/point", point)
			bwsGroup.GET("/achievements", achievements)
			bwsGroup.GET("/achievement", achievement)
			bwsGroup.POST("/point/unlock", authSvc.User, unlock)
			bwsGroup.POST("/binding", authSvc.User, binding)
			bwsGroup.POST("/award", authSvc.User, award)
			bwsGroup.GET("/lottery", authSvc.User, lottery)
			bwsGroup.GET("/lottery/check", authSvc.User, lotteryCheck)
			bwsGroup.GET("/redis/check", authSvc.User, redisInfo)
			bwsGroup.GET("/key/info", authSvc.User, keyInfo)
			bwsGroup.GET("/admin/check", authSvc.User, adminInfo)
		}
		bnjGroup := group.Group("/bnj2019")
		{
			bnjGroup.GET("/preview", authSvc.Guest, previewInfo)
			// TODO remove guest check
			bnjGroup.GET("/timeline", authSvc.Guest, timeline)
			bnjGroup.POST("/fail", fail)
			bnjGroup.POST("/reset", authSvc.User, reset)
			bnjGroup.POST("/reward", authSvc.User, reward)
		}
		kfcGroup := group.Group("/kfc")
		{
			kfcGroup.GET("/info", authSvc.User, kfcInfo)
			kfcGroup.GET("/use", kfcUse)
		}
	}
}

func internalRouter(e *bm.Engine) {
	group := e.Group("/x/internal/activity")
	{
		group.GET("/subject", vfySvc.Verify, subject)
		group.POST("/vote", vfySvc.Verify, vote)
		group.GET("/ltime", vfySvc.Verify, ltime)
		group.GET("/reddot", vfySvc.Verify, redDot)
		group.GET("/reddot/clear", vfySvc.Verify, authSvc.Guest, clearRedDot)
		group.GET("/object/stat/set", vfySvc.Verify, setSubjectStat)
		group.GET("/view/rank/set", vfySvc.Verify, setViewRank)
		group.GET("/like/content/set", vfySvc.Verify, setLikeContent)
		group.GET("/likeact/add", vfySvc.Verify, addLikeAct)
		group.GET("/likeact/cache", vfySvc.Verify, likeActCache)
		group.GET("/oids/info", vfySvc.Verify, likeOidsInfo)
		spGroup := group.Group("/sports")
		{
			spGroup.GET("/qq", vfySvc.Verify, qq)
			spGroup.GET("/news", vfySvc.Verify, news)
		}
		mactchGroup := group.Group("/match")
		{
			mactchGroup.GET("", matchs)
			mactchGroup.GET("/unstart", vfySvc.Verify, unStart)
			mactchGroup.POST("/cache/clear", clearCache)
			guGroup := mactchGroup.Group("/guess")
			{
				guGroup.GET("", vfySvc.Verify, guess)
				guGroup.GET("/list", vfySvc.Verify, listGuess)
				guGroup.POST("/add", vfySvc.Verify, addGuess)
			}
			foGroup := mactchGroup.Group("/follow")
			{
				foGroup.GET("", vfySvc.Verify, follow)
				foGroup.POST("/add", vfySvc.Verify, addFollow)
			}
		}
		initGroup := group.Group("/init")
		{
			initGroup.GET("/subject", vfySvc.Verify, subjectInit)
			initGroup.GET("/like", vfySvc.Verify, likeInit)
			initGroup.GET("/likeact", vfySvc.Verify, likeActCountInit)
			initGroup.GET("/subject/list", vfySvc.Verify, subjectLikeListInit)
		}
		//tmGroup := group.Group("/timemachine")
		//{
		//	tmGroup.GET("/start", startTmProc)
		//	tmGroup.GET("/stop", stopTmProc)
		//}
		kfcIGroup := group.Group("/kfc")
		{
			kfcIGroup.POST("/deliver", vfySvc.Verify, deliverKfc)
		}
		group.GET("/bnj2019/time/del", delTime)
	}
}

func ping(c *bm.Context) {
	if err := likeSvc.Ping(c); err != nil {
		log.Error("activity interface ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(nil, nil)
}
