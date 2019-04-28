package http

import (
	"net/http"

	"go-common/app/interface/main/dm2/conf"
	"go-common/app/interface/main/dm2/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	dmSvc       *service.Service
	antispamSvc *antispam.Antispam
	authSvc     *auth.Auth
	verifySvc   *verify.Verify
)

// Init http init.
func Init(c *conf.Config, s *service.Service) {
	dmSvc = s
	authSvc = auth.New(c.Auth)
	verifySvc = verify.New(c.Verify)
	antispamSvc = antispam.New(c.Antispam)
	engine := bm.NewServer(c.HTTPServer)
	engine.Use(bm.Recovery(), bm.Trace(), bm.Logger())
	innerRouter(engine)
	outerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start() error(%v)", err)
		panic(err)
	}
}

// innerRouter init router api path.
func innerRouter(e *bm.Engine) {
	group := e.Group("/x/internal/v2/dm", bm.CSRF())
	{
		group.GET("/search", verifySvc.VerifyUser, dmUpSearch)
		group.GET("/recent", verifySvc.VerifyUser, dmUpRecent)
		group.GET("/distribution", verifySvc.Verify, dmDistribution)
		group.POST("/edit/state", verifySvc.VerifyUser, editState)
		group.POST("/edit/pool", verifySvc.VerifyUser, editPool)
		group.POST("/mask/update", verifySvc.Verify, updateMask)
		group.POST("/subtitle/lan/add", verifySvc.Verify, subtitleLanAdd)
		group.POST("/subtitle/upos/callback", waveFormCallBack)

	}
}

func outerRouter(e *bm.Engine) {
	e.GET("/monitor/ping", ping)
	e.GET("/x/v1/dm/list.so", dmXML)
	group := e.Group("/x/v2/dm", bm.CSRF())
	{
		group.GET("/view", authSvc.Guest, view)

		group.GET("", dm)
		group.GET("/search", authSvc.User, dmUpSearch)
		group.GET("/ajax", ajaxDM)
		group.GET("/list.so", authSvc.Guest, dmSeg)
		group.GET("/list", authSvc.Guest, dmSegV2)
		group.GET("/judge/list", authSvc.Guest, judgeDM)
		group.GET("/thumbup/stats", authSvc.Guest, thumbupStats)
		group.GET("/history", authSvc.User, antispamSvc.ServeHTTP, dmHistory)
		group.GET("/history/list", authSvc.User, antispamSvc.ServeHTTP, dmHistoryV2)
		group.GET("/history/index", authSvc.User, antispamSvc.ServeHTTP, dmHistoryIndex)
		group.POST("/thumbup/add", authSvc.User, thumbupDM)
		group.POST("/post", authSvc.User, dmPost)
		group.POST("/edit/state", authSvc.User, editState)
		group.POST("/edit/pool", authSvc.User, editPool)
		group.POST("/filter/up/add", authSvc.User, antispamSvc.ServeHTTP, addUpFilterID)
		group.GET("/recent", authSvc.User, dmUpRecent)
		group.GET("/upper/config", authSvc.User, dmUpConfig)
		group.POST("/advance/config", authSvc.User, upAdvancePermit)
		group.GET("/filter/up", authSvc.User, upFilters)
		group.POST("/filter/up/edit", authSvc.User, editUpFilters)
		group.GET("/advert", authSvc.Guest, dmAdvert)

		subtitle := group.Group("/subtitle")
		{
			subtitle.GET("/lans", authSvc.User, subtitleLans)
			subtitle.POST("/del", authSvc.User, subtitleDel)
			subtitle.POST("/lock", authSvc.User, subtitleLock)
			subtitle.POST("/sign", authSvc.User, subtitleSign)
			subtitle.GET("/show", authSvc.User, subtitleShow)
			subtitle.GET("/archive/name", authSvc.User, subtitleArchiveName)
			subtitle.POST("/draft/save", authSvc.User, draftSave)
			subtitle.POST("/assit/audit", authSvc.User, assitAudit)
			subtitle.GET("/permission", authSvc.User, subtitlePermission)
			subtitle.GET("/waveform", authSvc.User, waveForm)
			subtitle.POST("/filter", authSvc.User, subtitleFilter)
			subtitle.GET("/report/tag", authSvc.User, subtitleReportTag)
			subtitle.POST("/report/add", authSvc.User, subtitleReportAdd)

			subtitle.GET("/search/assist", authSvc.User, searchAssist)
			subtitle.GET("/search/author/list", authSvc.User, authorList)
		}
	}

}

// ping check server ok.
func ping(c *bm.Context) {
	if err := dmSvc.Ping(c); err != nil {
		log.Error("dm2 service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
