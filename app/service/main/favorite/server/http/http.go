package http

import (
	"net/http"
	"strconv"

	"go-common/app/service/main/favorite/conf"
	"go-common/app/service/main/favorite/service"
	"go-common/library/log"
	"go-common/library/log/anticheat"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/supervisor"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	favSvc      *service.Service
	verifySvc   *verify.Verify
	antispamM   *antispam.Antispam
	supervisorM *supervisor.Supervisor
	collector   *anticheat.AntiCheat
)

// Init init router
func Init(c *conf.Config, svc *service.Service) {
	verifySvc = verify.New(c.Verify)
	antispamM = antispam.New(c.Antispam)
	supervisorM = supervisor.New(c.Supervisor)
	favSvc = svc
	if c.Infoc2 != nil {
		collector = anticheat.New(c.Infoc2)
	}
	// init outer router
	engineOut := bm.DefaultServer(c.BM)
	internalRouter(engineOut)
	// init serve
	if err := engineOut.Start(); err != nil {
		log.Error("engineOut.Start() error(%v)", err)
		panic(err)
	}
}

// internalRouter init internal router api path
func internalRouter(e *bm.Engine) {
	// init api
	e.Ping(ping)
	e.Register(register)
	favV3 := e.Group("/x/internal/v3/fav")
	{
		favV3.GET("", verifySvc.Verify, setMid, Favorites)
		favV3.GET("/test", setMid, Favorites)
		favV3.GET("/tlists", verifySvc.Verify, setMid, tlists)
		favV3.GET("/recents", verifySvc.VerifyUser, recentFavs)
		favV3.GET("/batch", verifySvc.Verify, batchFavs)
		favV3.POST("/add", verifySvc.VerifyUser, addFav)
		favV3.POST("/del", verifySvc.VerifyUser, delFav)
		favV3.POST("/madd", verifySvc.VerifyUser, multiAddFavs)
		favV3.POST("/mdel", verifySvc.VerifyUser, multiDelFavs)
		favV3.POST("/move", verifySvc.VerifyUser, moveFavs)
		favV3.POST("/copy", verifySvc.VerifyUser, copyFavs)
		favV3.GET("/favored", verifySvc.VerifyUser, isFavored)
		favV3.GET("/favoreds", verifySvc.VerifyUser, isFavoreds)
		favV3.GET("/users", verifySvc.Verify, userList)
		favV3.GET("/count", verifySvc.Verify, oidCount)
		favV3.GET("/counts", verifySvc.Verify, oidsCount)
		favV3.GET("/default", verifySvc.VerifyUser, inDefaultFolder)
	}
	folderV3 := e.Group("/x/internal/v3/fav/folder")
	{
		folderV3.GET("", verifySvc.Verify, setMid, userFolders)
		folderV3.GET("/multi", verifySvc.Verify, folders)
		folderV3.GET("/info", verifySvc.Verify, folderInfo)
		folderV3.GET("/count", verifySvc.Verify, cntUserFolders)
		folderV3.POST("/add", verifySvc.VerifyUser, addFolder)
		folderV3.POST("/update", verifySvc.VerifyUser, updateFolder)
		folderV3.POST("/del", verifySvc.VerifyUser, delFolder)
		folderV3.POST("/rename", verifySvc.VerifyUser, renameFolder)
		folderV3.POST("/public", verifySvc.VerifyUser, upAttrFolder)
		folderV3.POST("/sort", verifySvc.VerifyUser, sortFolders)
		folderV3.GET("/cleaned", verifySvc.VerifyUser, isCleaned)
		folderV3.POST("/clean", verifySvc.VerifyUser, cleanInvalidFavs)
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

// ping check server ok.
func ping(c *bm.Context) {
	if err := favSvc.Ping(c); err != nil {
		log.Error("favorite http service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// register check server ok.
func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
