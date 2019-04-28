package http

import (
	"net/http"
	"strconv"

	"go-common/app/interface/main/tag/conf"
	"go-common/app/interface/main/tag/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	svr       *service.Service
	authSvc   *auth.Auth
	verifySvc *verify.Verify
)

// Init .
func Init(c *conf.Config, s *service.Service) {
	svr = s
	authSvc = auth.New(c.Auth)
	verifySvc = verify.New(c.Verify)
	engine := bm.DefaultServer(c.BM)
	router(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func router(e *bm.Engine) {
	e.Ping(ping)
	// 公网服务接口
	outer := e.Group("/x/tag")
	{
		outer.GET("/info", authSvc.Guest, info)
		outer.GET("/minfo", authSvc.Guest, mInfo)
		outer.GET("/detail", authSvc.Guest, detail)
		outer.GET("/hots", authSvc.Guest, hotTags)
		outer.GET("/change/similar", changeSim)
		//视频详情页相关API.
		outer.POST("/archive/add", authSvc.User, addArcTagForOuter)
		outer.POST("/archive/del", authSvc.User, delArcTagForOuter)
		outer.POST("/archive/like", authSvc.User, likeArcTag)
		outer.POST("/archive/hate", authSvc.User, hateArcTag)
		outer.POST("/archive/like2", authSvc.User, likeArcTag)
		outer.POST("/archive/hate2", authSvc.User, hateArcTag)
		outer.POST("/archive/report", authSvc.User, reportArcTag) // 详情页举报
		outer.POST("/archive/add/report", authSvc.User, logReport)
		outer.POST("/archive/del/report", authSvc.User, logReport)
		outer.GET("/archive/tags", authSvc.Guest, arcTags)
		outer.GET("/archive/log", authSvc.Guest, arcTagLog)
		//订阅相关API
		outer.POST("/subscribe/add", authSvc.User, addSub)
		outer.POST("/subscribe/cancel", authSvc.User, cancelSub)
		outer.GET("/subscribe/tags", authSvc.User, subTags)
		outer.GET("/subscribe/archives", authSvc.User, subArcs) // for app feed
		// rank相关API.
		outer.GET("/ranking/archives", newArcs) // newest archives
	}
	// 内网服务接口
	inner := e.Group("/x/internal/tag")
	{
		inner.GET("/info", verifySvc.Verify, setMid, info)
		inner.GET("/minfo", verifySvc.Verify, setMid, mInfo)
		inner.GET("/infos", verifySvc.Verify, setMid, infos)
		inner.POST("/batch/info", verifySvc.Verify, setMid, tagBatchInfo)
		inner.GET("/detail", verifySvc.Verify, setMid, detail)
		inner.GET("/detail/ranking", verifySvc.Verify, detailRankArc)
		inner.GET("/hots", verifySvc.Verify, setMid, hotTags)
		inner.GET("/similar", verifySvc.Verify, similarTags)
		inner.GET("/change/similar", changeSim)
		inner.POST("/activity/add", verifySvc.Verify, addActivityTag)
		inner.GET("/recommand", recommandTag) // 投稿推荐tag
		inner.GET("/synonym", synonymTag)     // 大数据同义词tag
		// dynamic-servicey依赖接口
		inner.GET("/hotmap", verifySvc.Verify, hotMap)
		inner.GET("/prids", verifySvc.Verify, prids)
		inner.GET("/check", verifySvc.VerifyUser, checkTagName)
		// 视频详情页相关API.
		inner.GET("/archive/tags", verifySvc.Verify, setMid, arcTags)
		inner.GET("/archive/multi/tags", verifySvc.Verify, setMid, multiArcTags)
		inner.POST("/archive/upbind", verifySvc.Verify, setMid, upBind)
		inner.POST("/archive/adminbind", verifySvc.Verify, setMid, adminBind)
		//订阅相关API
		inner.POST("/subscribe/add", verifySvc.VerifyUser, addSub)
		inner.POST("/subscribe/cancel", verifySvc.VerifyUser, cancelSub)
		inner.GET("/subscribe/tags", verifySvc.Verify, setMid, subTags)
		inner.GET("/subscribe/archives", verifySvc.Verify, setMid, subArcs)
		inner.GET("/subscribe/custom/sort/tags", verifySvc.Verify, setMid, customSortTags)
		inner.POST("/subscribe/custom/sort/update", verifySvc.VerifyUser, upCustomSortTags)

		inner.GET("/ranking/archives", verifySvc.Verify, newArcs)
	}
	// 内网服务V2接口 通用平台接入业务
	innerV2 := e.Group("/x/internal/v2/tag")
	{
		innerV2.GET("/tags", verifySvc.Verify, setMid, platformListTag)  // 资源绑定tag列表
		innerV2.POST("/bind/up", verifySvc.Verify, platformUpBind)       // up bind
		innerV2.POST("/bind/admin", verifySvc.Verify, platformAdminBind) // admin bind
	}
	innerChannel := e.Group("/x/internal/channel")
	{
		innerChannel.GET("/category", verifySvc.Verify, channelCategory)
		innerChannel.GET("/rule", verifySvc.Verify, channelRule)
		innerChannel.GET("/list", verifySvc.Verify, channeList)
		innerChannel.GET("/recommand", verifySvc.Verify, channelRecommand)
		innerChannel.GET("/discover", verifySvc.Verify, channelDiscover)
		innerChannel.GET("/square", verifySvc.Verify, channelSquare)
		innerChannel.GET("/detail", verifySvc.Verify, setMid, channelDetail)
		innerChannel.GET("/subscribe/customs", verifySvc.Verify, setMid, channelSubCustoms)
		innerChannel.GET("/subscribe/customs/update", verifySvc.VerifyUser, upChannelSubCustoms)
		innerChannel.GET("/resources", verifySvc.Verify, channelResource)
		innerChannel.GET("/resource/checkback", verifySvc.Verify, resourceCheckBack)
		innerChannel.GET("/resource/infos", verifySvc.Verify, resourceInfos)
	}
}

func ping(c *bm.Context) {
	if err := svr.Ping(c); err != nil {
		log.Error("tag service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func setMid(c *bm.Context) {
	var (
		err error
		mid int64
	)
	req := c.Request
	midStr := req.Form.Get("mid")
	if midStr != "" {
		mid, err = strconv.ParseInt(midStr, 10, 64)
		if err != nil {
			c.JSON(nil, err)
			c.Abort()
			return
		}
	}
	c.Set("mid", mid)
}
