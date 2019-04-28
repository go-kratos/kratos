package http

import (
	"go-common/library/net/metadata"
	"net/http"

	"go-common/app/service/live/userexp/conf"
	"go-common/app/service/live/userexp/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	expSvr *service.Service
)

// Init init account service.
func Init(c *conf.Config, s *service.Service) {
	// init service
	expSvr = s
	// init external router
	engineOut := bm.DefaultServer(c.BM.Inner)
	innerRouter(engineOut)
	// init Inner serve
	if err := engineOut.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
	engineLocal := bm.DefaultServer(c.BM.Local)
	localRouter(engineLocal)
	// init Local serve
	if err := engineLocal.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

// outerRouter init outer router api path.
func innerRouter(e *bm.Engine) {
	// init api
	e.GET("/monitor/ping", ping)
	group := e.Group("/x/internal/liveexp")
	{
		group.GET("/level/get", level)
		group.GET("/level/mulGet", multiGetLevel)
		group.GET("/level/addUexp", addUexp)
		group.GET("/level/addRexp", addRexp)
	}
}

// localRouter init local router.
func localRouter(e *bm.Engine) {
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := expSvr.Ping(c); err != nil {
		log.Error("live-userexp  http service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
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
		//log.Warn("live-userexo infoc get sid failed error(%v)", err)
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

//func checkHealth(c wctx.Context) {
//	_, err := os.Stat(conf.Conf.CheckFile)
//	if os.IsNotExist(err) {
//		http.Error(c.Response(), http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
//	}
//	c.Cancel()
//}
