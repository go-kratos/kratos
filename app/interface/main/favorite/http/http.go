package http

import (
	"net/http"
	"strconv"

	"go-common/app/interface/main/favorite/conf"
	"go-common/app/interface/main/favorite/service"
	"go-common/library/log"
	"go-common/library/log/anticheat"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/supervisor"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	favSvc      *service.Service
	authSvc     *auth.Auth
	verifySvc   *verify.Verify
	antispamM   *antispam.Antispam
	supervisorM *supervisor.Supervisor
	collector   *anticheat.AntiCheat
)

// Init init router
func Init(c *conf.Config, svc *service.Service) {
	verifySvc = verify.New(c.Verify)
	authSvc = auth.New(c.Auth)
	antispamM = antispam.New(c.Antispam)
	supervisorM = supervisor.New(c.Supervisor)
	favSvc = svc
	if c.Infoc2 != nil {
		collector = anticheat.New(c.Infoc2)
	}
	// init outer router
	engineOut := bm.DefaultServer(c.BM)
	outerRouter(engineOut)
	internalRouter(engineOut)
	// init Out serve
	if err := engineOut.Start(); err != nil {
		log.Error("engineOut.Start() error(%v)", err)
		panic(err)
	}
}

// outerRouter init outer router api path
func outerRouter(e *bm.Engine) {
	// init api
	e.GET("/monitor/ping", ping)
	folderG := e.Group("/x/v2/fav/folder")
	{
		folderG.GET("", authSvc.Guest, videoFolders)
		folderG.POST("/add", authSvc.User, antispamM.ServeHTTP, supervisorM.ServeHTTP, addVideoFolder)
		folderG.POST("/del", authSvc.User, antispamM.ServeHTTP, delVideoFolder)
		folderG.POST("/rename", authSvc.User, antispamM.ServeHTTP, supervisorM.ServeHTTP, renameVideoFolder)
		folderG.POST("/public", authSvc.User, antispamM.ServeHTTP, upStateVideoFolder)
		folderG.POST("/sort", authSvc.User, antispamM.ServeHTTP, sortVideoFolders)
	}
	videoG := e.Group("/x/v2/fav/video")
	{
		videoG.GET("", authSvc.Guest, favVideo)
		videoG.GET("/tlist", authSvc.Guest, tidList)
		videoG.GET("/newest", authSvc.User, favVideoNewest)
		videoG.POST("/add", authSvc.User, antispamM.ServeHTTP, addFavVideo)
		videoG.POST("/del", authSvc.User, antispamM.ServeHTTP, delFavVideo)
		videoG.POST("/mdel", authSvc.User, antispamM.ServeHTTP, delFavVideos)
		videoG.POST("/move", authSvc.User, antispamM.ServeHTTP, moveFavVideos)
		videoG.POST("/copy", authSvc.User, antispamM.ServeHTTP, copyFavVideos)
		videoG.GET("/favoureds", authSvc.User, isFavoureds)
		videoG.GET("/favoured", authSvc.User, isFavoured)
		videoG.GET("/default", authSvc.User, inDefaultFav)
		videoG.GET("/cleaned", authSvc.User, isCleaned)
		videoG.POST("/clean", authSvc.User, cleanInvalidArcs)
	}
	topicG := e.Group("/x/v2/fav/topic")
	{
		topicG.POST("/add", authSvc.User, antispamM.ServeHTTP, addFavTopic)
		topicG.POST("/del", authSvc.User, antispamM.ServeHTTP, delFavTopic)
		topicG.GET("/favoured", authSvc.User, isTopicFavoured)
		topicG.GET("", authSvc.User, favTopics)
	}
}

// internalRouter init internal router api path
func internalRouter(e *bm.Engine) {
	// init api
	folderG := e.Group("/x/internal/v2/fav/folder")
	{
		folderG.GET("", verifySvc.Verify, setMid, videoFolders)
		folderG.POST("/add", verifySvc.VerifyUser, addVideoFolder)
		folderG.POST("/del", verifySvc.VerifyUser, delVideoFolder)
		folderG.POST("/rename", verifySvc.VerifyUser, renameVideoFolder)
		folderG.POST("/public", verifySvc.VerifyUser, upStateVideoFolder)
		folderG.POST("/sort", verifySvc.VerifyUser, sortVideoFolders)
	}
	videoG := e.Group("/x/internal/v2/fav/video")
	{
		videoG.GET("", verifySvc.Verify, setMid, favVideo)
		videoG.GET("/tlist", verifySvc.Verify, setMid, tidList)
		videoG.GET("/newest", verifySvc.VerifyUser, favVideoNewest)
		videoG.POST("/add", verifySvc.VerifyUser, addFavVideo)
		videoG.POST("/del", verifySvc.VerifyUser, delFavVideo)
		videoG.POST("/mdel", verifySvc.VerifyUser, delFavVideos)
		videoG.POST("/move", verifySvc.VerifyUser, moveFavVideos)
		videoG.POST("/copy", verifySvc.VerifyUser, copyFavVideos)
		videoG.GET("/favoureds", verifySvc.VerifyUser, isFavoureds)
		videoG.GET("/favoured", verifySvc.VerifyUser, isFavoured)
		videoG.GET("/default", verifySvc.VerifyUser, inDefaultFav)
		videoG.GET("/cleaned", verifySvc.VerifyUser, isCleaned)
		videoG.POST("/clean", verifySvc.VerifyUser, cleanInvalidArcs)
	}
	topicG := e.Group("/x/internal/v2/fav/topic")
	{
		topicG.POST("/add", verifySvc.VerifyUser, addFavTopic)
		topicG.POST("/del", verifySvc.VerifyUser, delFavTopic)
		topicG.GET("/favoured", verifySvc.VerifyUser, isTopicFavoured)
		topicG.GET("", verifySvc.VerifyUser, favTopics)
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
		if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
			c.JSON(nil, err)
			c.Abort()
			return
		}
	}
	c.Set("mid", mid)
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := favSvc.Ping(c); err != nil {
		log.Error("favorite http service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
