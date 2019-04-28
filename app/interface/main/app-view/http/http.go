package http

import (
	"go-common/app/interface/main/app-view/conf"
	"go-common/app/interface/main/app-view/service/report"
	"go-common/app/interface/main/app-view/service/view"
	"go-common/library/log"
	"go-common/library/log/anticheat"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/proxy"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/queue/databus"
)

var (
	viewSvr   *view.Service
	reportSvr *report.Service
	authSvr   *auth.Auth
	verifySvc *verify.Verify
	// databus
	userActPub *databus.Databus
	dislikePub *databus.Databus
	collector  *anticheat.AntiCheat
)

type userAct struct {
	Client   string `json:"client"`
	Buvid    string `json:"buvid"`
	Mid      int64  `json:"mid"`
	Time     int64  `json:"time"`
	From     string `json:"from"`
	Build    string `json:"build"`
	ItemID   string `json:"item_id"`
	ItemType string `json:"item_type"`
	Action   string `json:"action"`
	ActionID string `json:"action_id"`
	Extra    string `json:"extra"`
}

type cmDislike struct {
	ID         int64  `json:"id"`
	Buvid      string `json:"buvid"`
	Goto       string `json:"goto"`
	Mid        int64  `json:"mid"`
	ReasonID   int64  `json:"reason_id"`
	CMReasonID int64  `json:"cm_reason_id"`
	UpperID    int64  `json:"upper_id"`
	Rid        int64  `json:"rid"`
	TagID      int64  `json:"tag_id"`
	ADCB       string `json:"ad_cb"`
}

// Init init http
func Init(c *conf.Config) {
	initService(c)
	collector = anticheat.New(c.InfocCoin)
	// init external router
	engineOut := bm.DefaultServer(c.BM.Outer)
	outerRouter(engineOut)
	// init outer server
	if err := engineOut.Start(); err != nil {
		log.Error("engineOut.Start() error(%v)", err)
		panic(err)
	}
}

func initService(c *conf.Config) {
	verifySvc = verify.New(nil)
	authSvr = auth.New(nil)
	viewSvr = view.New(c)
	reportSvr = report.New(c)
	// databus
	userActPub = databus.New(c.UseractPub)
	dislikePub = databus.New(c.DislikePub)
}

func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	// view
	proxyHandler := proxy.NewZoneProxy("sh004", "http://sh001-app.bilibili.com")
	view := e.Group("/x/v2/view")
	view.GET("", verifySvc.Verify, authSvr.GuestMobile, viewIndex)
	view.GET("/page", verifySvc.Verify, authSvr.GuestMobile, viewPage)
	view.GET("/video/shot", verifySvc.Verify, videoShot)
	view.POST("/share/add", verifySvc.Verify, authSvr.GuestMobile, addShare)
	view.POST("/coin/add", proxyHandler, authSvr.UserMobile, addCoin)
	view.POST("/fav/add", authSvr.UserMobile, addFav)
	view.POST("/ad/dislike", authSvr.GuestMobile, adDislike)
	view.GET("/report", verifySvc.Verify, copyWriter)
	view.POST("/report/add", authSvr.UserMobile, addReport)
	view.POST("/report/upload", verifySvc.Verify, upload)
	view.POST("/like", proxyHandler, authSvr.UserMobile, like)
	view.POST("/dislike", authSvr.UserMobile, dislike)
	view.POST("/vip/playurl", authSvr.UserMobile, vipPlayURL)
	view.GET("/follow", authSvr.GuestMobile, follow)
	view.GET("/upper/recmd", authSvr.GuestMobile, upperRecmd)
	view.POST("/like/triple", authSvr.UserMobile, likeTriple)
	// bnj2019
	view.GET("/bnj2019", verifySvc.Verify, authSvr.GuestMobile, bnj2019)
	view.GET("/bnj2019/list", verifySvc.Verify, authSvr.GuestMobile, bnjList)
	view.GET("/bnj2019/item", verifySvc.Verify, authSvr.GuestMobile, bnjItem)
}
