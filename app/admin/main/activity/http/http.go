package http

import (
	"go-common/app/admin/main/activity/conf"
	"go-common/app/admin/main/activity/service"
	"go-common/app/admin/main/activity/service/kfc"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	actSrv  *service.Service
	authSrv *permit.Permit
	kfcSrv  *kfc.Service
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	actSrv = s
	kfcSrv = kfc.New(c)
	authSrv = permit.New(c.Auth)
	engine := bm.DefaultServer(c.HTTPServer)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("httpx.Serve error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/x/admin/activity")
	{
		g.GET("/arcs", archives)
		gapp := g.Group("/matchs", authSrv.Permit("ACT_MATCHS_MGT_TEST"))
		{
			gapp.POST("/add", addMatch)
			gapp.POST("/save", saveMatch)
			gapp.GET("/info", matchInfo)
			gapp.GET("/list", matchList)
		}
		gappO := g.Group("/matchs/object", authSrv.Permit("ACT_MATCHS_MGT_TEST"))
		{
			gappO.POST("/add", addMatchObject)
			gappO.POST("/save", saveMatchObject)
			gappO.GET("/info", matchObjectInfo)
			gappO.GET("/list", matchObjectList)
		}
		gappSuject := g.Group("/subject")
		{
			gappSuject.GET("/list", listInfosAll)
			gappSuject.GET("/videos", videoList)
			gappSuject.POST("/add", addActSubject)
			gappSuject.POST("/up", updateInfoAll)
			gappSuject.GET("/protocol", subPro)
			gappSuject.GET("/conf", timeConf)
			gappSuject.GET("/articles", article)
		}
		gappLikes := g.Group("/likes")
		{
			gappLikes.GET("/list", likesList)
			gappLikes.GET("/lids", likes)
			gappLikes.POST("/add", addLike)
			gappLikes.POST("/up", upLike)
			gappLikes.POST("/up/reply", upListContent)
			gappLikes.POST("/up/wid", upWid)
			gappLikes.POST("/add/pic", addPic)
			gappLikes.POST("/batch/wid", batchLikes)
		}
		gappKfc := g.Group("kfc")
		{
			gappKfc.GET("/list", kfcList)
		}
		gappBws := g.Group("/bws")
		{
			gappBws.POST("/add", addBws)
			gappBws.POST("/save", saveBws)
			gappBws.GET("/info", bwsInfo)
			gappBws.GET("/list", bwsList)
			gappAchievement := gappBws.Group("/achievement")
			{
				gappAchievement.POST("/add", addBwsAchievement)
				gappAchievement.POST("/save", saveBwsAchievement)
				gappAchievement.GET("/info", bwsAchievement)
				gappAchievement.GET("/list", bwsAchievements)
			}
			gappField := gappBws.Group("/field")
			{
				gappField.POST("/add", addBwsField)
				gappField.POST("/save", saveBwsField)
				gappField.GET("/info", bwsField)
				gappField.GET("/list", bwsFields)
			}
			gappPoint := gappBws.Group("/point")
			{
				gappPoint.POST("/add", addBwsPoint)
				gappPoint.POST("/save", saveBwsPoint)
				gappPoint.GET("/info", bwsPoint)
				gappPoint.GET("/list", bwsPoints)
			}
			gappUser := gappBws.Group("/user")
			{
				gappUser.POST("/add", addBwsUser)
				gappUser.POST("/save", saveBwsUser)
				gappUser.GET("/info", bwsUser)
				gappUser.GET("/list", bwsUsers)
				gappUserAchievement := gappUser.Group("/achievement")
				{
					gappUserAchievement.POST("/add", addBwsUserAchievement)
					gappUserAchievement.POST("/save", saveBwsUserAchievement)
					gappUserAchievement.GET("/info", bwsUserAchievement)
					gappUserAchievement.GET("/list", bwsUserAchievements)
				}
				gappUserPoint := gappUser.Group("/point")
				{
					gappUserPoint.POST("/add", addBwsUserPoint)
					gappUserPoint.POST("/save", saveBwsUserPoint)
					gappUserPoint.GET("/info", bwsUserPoint)
					gappUserPoint.GET("/list", bwsUserPoints)
				}
			}
		}

	}
}

func ping(c *bm.Context) {
	if err := actSrv.Ping(c); err != nil {
		c.Error = err
		c.AbortWithStatus(503)
	}
}
