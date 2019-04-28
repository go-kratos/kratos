package http

import (
	"net/http"

	"go-common/app/service/video/stream-mng/conf"
	"go-common/app/service/video/stream-mng/middleware"
	"go-common/app/service/video/stream-mng/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	srv     *service.Service
	vfy     *verify.Verify
	authSvr *auth.Auth
)

// Init init
func Init(c *conf.Config, svc *service.Service) {
	srv = svc
	vfy = verify.New(c.Verify)
	authSvr = auth.New(c.Auth)
	engine := bm.DefaultServer(c.BM)
	engine.Use(middleware.Logger())
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/video/stream-mng")
	{
		g.GET("/", alive)

		g.POST("/stream/backup", createBackupStream)
		g.POST("/stream/offical", createOfficalStream)
		g.POST("/stream/validate", streamValidate)
		g.GET("/stream/old/getbyroomid", getOldStreamInfoByRoomID)
		g.GET("/stream/old/getbyname", getOldStreamInfoByStreamName)

		g.GET("/notifymaskbyroomid", saveMaskByRoomID)         //控制一个主流是否需要转蒙版直播流
		g.GET("/notifymaskbystreamname", saveMaskByStreamName) //控制一个主流是否可提供蒙版给PLAYURL
		g.POST("/addhotstream", addHotStream)                  //增加热流到redis
		g.GET("/getstream", getStream)
		g.GET("/getmultistreams", getMultiStreams)

		g.GET("/stream/getRoomIdByStreamName", getRoomIDByStreamName)
		g.GET("/stream/getStreamNameByRoomId", getStreamNameByRoomID)

		g.POST("/stream/changeSrc", changeSrc)
		g.GET("/stream/cut", cutStream)
		g.GET("/stream/cutmobilestream", authSvr.User, cutStreamByMobile)
		g.GET("/stream/getLastTime", getStreamLastTime)
		g.GET("/stream/getAdapterStream", getAdapterStreamByStreamName)
		g.GET("/stream/getSrcByRoomID", getSrcByRoomID)
		g.GET("/stream/getLineListByRoomID", getLineListByRoomID)
		g.GET("/shot/getSinglePic", getSingleScreenShot)
		g.GET("/shot/getMultiPic", getMultiScreenShot)
		g.GET("/shot/getOriginPic", getOriginScreenShotPic)
		g.GET("/shot/getperiodpic", getTimePeriodScreenShot)
		g.POST("/stream/clearstreamstatus", clearStreamStatus)
		g.GET("/stream/getRoomRtmp", getRoomRtmp)                   // 拜年祭推流码
		g.GET("/stream/getUpStreamRtmp", getUpStreamRtmp)           // 后台调用，无需鉴权
		g.GET("/stream/getmobilertmp", authSvr.User, getMobileRtmp) // 移动端调用
		g.GET("/stream/getwebrtmp", authSvr.User, getWebRtmp)       // 被web端和pc_link调用
		g.GET("/stream/live", checkLiveStreamList)

		// 删除room_id缓存的接口，防止缓存问题出现的bug
		g.GET("/stream/clearcache", clearRoomCacheByRID)

		// 查询更改记录
		g.GET("/change/getchangeLog", getChangeLogByRoomID)

		// 查询统计上行调度信息
		g.GET("/summary/upstream", getSummaryUpStreamRtmp)
		g.GET("/summary/isp", getSummaryUpStreamISP)
		g.GET("/summary/country", getSummaryUpStreamCountry)
		g.GET("/summary/platform", getSummaryUpStreamPlatform)
		g.GET("/summary/city", getSummaryUpStreamCity)
	}

	g2 := e.Group("/live_stream/v1/StreamThird")
	{
		g2.POST("/stream_validate", streamValidate)
		g2.POST("/open_notify", openNotify)
		g2.POST("/close_notify", closeNotify)
	}

	g7 := e.Group("/live_stream/v1/StreamList")
	{
		g7.GET("/get_stream_by_roomId", authSvr.User, getWebRtmp)
	}

	g8 := e.Group("/live_stream/v1/UpStreamExt")
	{
		g8.GET("/get_by_room", authSvr.User, getMobileRtmp)
		g8.GET("/pause_by_room", authSvr.User, cutStreamByMobile)
	}
}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

func alive(c *bm.Context) {
	c.String(0, "Golang 大法好 !!!")
}
