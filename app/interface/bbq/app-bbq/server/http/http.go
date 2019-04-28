package http

import (
	"fmt"
	"go-common/library/ecode"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"net/http"

	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/conf"
	"go-common/app/interface/bbq/app-bbq/service"
	xauth "go-common/app/interface/bbq/common/auth"
	chttp "go-common/app/interface/bbq/common/http"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/trace"
)

var (
	srv              *service.Service
	vfy              *verify.Verify
	authSrv          *xauth.BannedAuth
	cfg              *conf.Config
	logger           *chttp.UILog
	likeAntiSpam     *antispam.Antispam
	relationAntiSpam *antispam.Antispam
	replyAntiSpam    *antispam.Antispam
	uploadAntiSpam   *antispam.Antispam
	reportAntiSpam   *antispam.Antispam
)

// Init init
func Init(c *conf.Config) {
	cfg = c
	initAntiSpam(c)
	logger = chttp.New(c.Infoc)
	srv = service.New(c)
	vfy = verify.New(c.Verify)
	authSrv = xauth.NewBannedAuth(c.Auth, c.MySQL)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func initAntiSpam(c *conf.Config) {
	var antiConfig *antispam.Config
	var exists bool
	if antiConfig, exists = c.AntiSpam["like"]; !exists {
		panic("lose like anti_spam config")
	}
	relationAntiSpam = antispam.New(antiConfig)
	if antiConfig, exists = c.AntiSpam["relation"]; !exists {
		panic("lose relation anti_spam config")
	}
	likeAntiSpam = antispam.New(antiConfig)
	if antiConfig, exists = c.AntiSpam["reply"]; !exists {
		panic("lose reply anti_spam config")
	}
	replyAntiSpam = antispam.New(antiConfig)
	if antiConfig, exists = c.AntiSpam["upload"]; !exists {
		panic("lose upload anti_spam config")
	}
	uploadAntiSpam = antispam.New(antiConfig)

	if antiConfig, exists = c.AntiSpam["report"]; !exists {
		panic("lose report anti_spam config")
	}
	reportAntiSpam = antispam.New(antiConfig)
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/bbq/app-bbq", wrapBBQ)
	{
		//用户登录
		g.GET("/user/login", authSrv.User, login)
		g.POST("/user/logout", authSrv.Guest, bm.Mobile(), pushLogout)

		//用户相关
		g.GET("/user/base", authSrv.User, userBase)
		// 所有字段都需要携带修改
		g.POST("/user/base/edit", authSrv.User, userBaseEdit)

		g.POST("/user/like/add", authSrv.User, likeAntiSpam.ServeHTTP, addUserLike)
		g.POST("/user/like/cancel", authSrv.User, likeAntiSpam.ServeHTTP, cancelUserLike)
		g.GET("/user/like/list", userLikeList)
		g.POST("/user/unlike", authSrv.User, likeAntiSpam.ServeHTTP, userUnLike)
		g.GET("/user/follow/list", authSrv.Guest, userFollowList)
		g.GET("/user/fan/list", authSrv.Guest, userFanList)
		g.GET("/user/black/list", authSrv.User, userBlackList)
		g.POST("/user/relation/modify", authSrv.User, relationAntiSpam.ServeHTTP, userRelationModify)

		g.GET("/search/hot/word", hotWord)

		// feed关注列表页
		g.GET("/feed/list", authSrv.User, feedList)
		// feed关注页红点
		g.GET("/feed/update_num", authSrv.User, feedUpdateNum)
		// space发布列表页
		g.GET("/space/sv/list", authSrv.Guest, spaceSvList)
		// space 用户详情况主／客
		g.GET("/space/user/profile", authSrv.Guest, spaceUserProfile)
		// 详情页up主发布列表
		g.GET("/detail/sv/list", authSrv.Guest, detailSvList)

		//视频相关
		g.GET("/sv/list", authSrv.Guest, bm.Mobile(), svList)
		g.GET("/sv/playlist", authSrv.Guest, bm.Mobile(), svPlayList) // playurl过期时候请求
		g.GET("/sv/detail", authSrv.Guest, svDetail)
		g.GET("/sv/stat", authSrv.Guest, bm.Mobile(), svStatistics)
		g.GET("/sv/relate", authSrv.Guest, svRelList)
		g.POST("/sv/del", authSrv.User, svDel)
		//搜索相关
		g.GET("/search/sv", authSrv.Guest, videoSearch)
		g.GET("/search/user", authSrv.Guest, userSearch)
		g.GET("/search/sug", authSrv.Guest, sug)
		g.GET("/search/topic", authSrv.Guest, topicSearch)
		//发现页
		g.GET("/discovery", authSrv.Guest, discoveryList)
		//话题详情页
		g.GET("/topic/detail", authSrv.Guest, topicDetail)

		// 用户location
		g.GET("/location/all", authSrv.User, locationAll)
		g.GET("/location", authSrv.User, location)

		//图片上传
		g.POST("/img/upload", authSrv.User, uploadAntiSpam.ServeHTTP, upload)

		// 客户端分享链接
		g.GET("/share", authSrv.Guest, shareURL)
		g.GET("/share/callback", authSrv.Guest, shareCallback)

		// 邀请函接口（内测版，公测删除）内测取消
		// g.GET("/invitation/download", invitationDownload)

		// App全局设置接口
		g.GET("/setting", appSetting)
		g.GET("/package", appPackage)
	}
	//评论组
	r := e.Group("/bbq/app-bbq/reply", wrapBBQ, commentInit)
	{
		//评论相关
		r.GET("/cursor", commentCloseRead, authSrv.Guest, commentCursor)
		r.POST("/add", commentCloseWrite, authSrv.User, phoneCheck, replyAntiSpam.ServeHTTP, commentAdd)
		r.POST("/action", commentCloseWrite, authSrv.User, likeAntiSpam.ServeHTTP, commentLike)
		r.GET("/", commentCloseRead, authSrv.Guest, commentList)
		r.GET("/reply/cursor", commentCloseRead, authSrv.Guest, commentSubCursor)
	}

	// 举报接口
	report := e.Group("/bbq/app-bbq/report", wrapBBQ)
	{
		report.GET("/config", authSrv.Guest, bm.Mobile(), reportConfig)
		report.POST("/report", authSrv.Guest, bm.Mobile(), reportAntiSpam.ServeHTTP, reportReport)
	}

	// 播放数据收集
	d := e.Group("/bbq/app-bbq/data", wrapBBQ)
	{
		d.GET("/collect", authSrv.Guest, bm.Mobile(), videoPlay)
	}

	// 通知中心，需要登录
	p := e.Group("/bbq/app-bbq/notice/center", authSrv.User, wrapBBQ)
	{
		p.GET("/num", noticeNum)
		p.GET("/overview", noticeOverview)
		p.GET("/list", noticeList)
	}

	// 推送相关
	push := e.Group("/bbq/app-bbq/push", wrapBBQ, authSrv.Guest, bm.Mobile())
	{
		push.POST("/register", pushRegister)
		push.GET("/callback", pushCallback)
	}

	//视频上传相关
	upload := e.Group("/bbq/app-bbq/upload/sv", authSrv.Guest)
	{
		upload.POST("/preupload", perUpload)
		upload.POST("/callback", callBack)
		upload.GET("/check", authSrv.User, uploadCheck)
		upload.POST("/homeimg", authSrv.User, homeimg)
	}
}

func commentCloseWrite(ctx *bm.Context) {
	if conf.Conf.Comment.CloseWrite {
		ctx.JSON(struct{}{}, ecode.OK)
		ctx.Abort()
	}
}
func commentCloseRead(ctx *bm.Context) {
	if conf.Conf.Comment.CloseRead {
		ctx.JSON(struct{}{}, ecode.OK)
		ctx.Abort()
	}
}

//wrapRes 为返回头添加BBQ自定义字段
func wrapBBQ(ctx *bm.Context) {
	chttp.WrapHeader(ctx)

	// Base params
	req := ctx.Request
	base := new(v1.Base)
	ctx.Bind(base)
	base.BUVID = req.Header.Get("Buvid")
	ctx.Set("BBQBase", base)

	// QueryID
	qid := base.QueryID
	if base.QueryID == "" {
		tracer, _ := trace.FromContext(ctx.Context)
		qid = fmt.Sprintf("%s", tracer)
	}
	ctx.Set("QueryID", qid)
}

// phoneCheck 进行手机校验
func phoneCheck(ctx *bm.Context) {
	midValue, exists := ctx.Get("mid")
	if !exists {
		ctx.JSON(nil, ecode.NoLogin)
		ctx.Abort()
		return
	}
	mid := midValue.(int64)
	err := srv.PhoneCheck(ctx, mid)
	if err != nil {
		ctx.JSON(nil, err)
		ctx.Abort()
		return
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

func uiLog(ctx *bm.Context, action int, ext interface{}) {
	logger.Infoc(ctx, action, ext)
}
