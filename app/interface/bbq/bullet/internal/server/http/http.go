package http

import (
	"net/http"
	"strings"

	"go-common/app/interface/bbq/bullet/api"
	"go-common/app/interface/bbq/bullet/internal/conf"
	"go-common/app/interface/bbq/bullet/internal/model"
	"go-common/app/interface/bbq/bullet/internal/service"
	xauth "go-common/app/interface/bbq/common/auth"
	chttp "go-common/app/interface/bbq/common/http"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	antispam "go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/verify"

	"github.com/golang/protobuf/ptypes/empty"
)

var (
	// TODO:verify
	vfy            *verify.Verify
	svc            *service.Service
	authSrv        *xauth.BannedAuth
	logger         *chttp.UILog
	bulletAntiSpam *antispam.Antispam
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	svc = s
	initAntiSpam(c)
	vfy = verify.New(c.Verify)
	authSrv = xauth.NewBannedAuth(c.Auth, c.OnlineMySQL)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
	logger = chttp.New(c.Infoc)
}

func initAntiSpam(c *conf.Config) {
	var antiConfig *antispam.Config
	var exists bool
	if antiConfig, exists = c.AntiSpam["bullet"]; !exists {
		panic("lose bullet anti_spam config")
	}
	bulletAntiSpam = antispam.New(antiConfig)
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/bbq/bullet", wrapBBQ)
	{
		g.GET("/content/list", authSrv.Guest, contentList)
		g.GET("/content/get", authSrv.Guest, contentGet)
		g.POST("/content/post", authSrv.User, phoneCheck, bulletAntiSpam.ServeHTTP, contentPost)
	}
}

func ping(ctx *bm.Context) {
	if err := svc.Ping(ctx); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

//wrapRes 为返回头添加BBQ自定义字段
func wrapBBQ(ctx *bm.Context) {
	chttp.WrapHeader(ctx)
}

func uiLog(ctx *bm.Context, action int, ext interface{}) {
	logger.Infoc(ctx, action, ext)
}

func contentGet(c *bm.Context) {
	if conf.Conf.BulletConfig.CloseRead {
		c.JSON([]*api.Bullet{}, nil)
		return
	}

	arg := &api.ListBulletReq{}
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.Oid == 0 {
		res := new([]interface{})
		c.JSON(res, ecode.DanmuGetErr)
		return
	}

	if midValue, exists := c.Get("mid"); exists {
		arg.Mid = midValue.(int64)
	}

	c.JSON(svc.ContentGet(c, arg))
}

func contentList(c *bm.Context) {
	if conf.Conf.BulletConfig.CloseRead {
		c.JSON(api.ListBulletReply{HasMore: false}, nil)
		return
	}

	arg := &api.ListBulletReq{}
	if err := c.Bind(arg); err != nil {
		return
	}

	if arg.Oid == 0 {
		res := new([]interface{})
		c.JSON(res, ecode.DanmuGetErr)
		return
	}

	if midValue, exists := c.Get("mid"); exists {
		arg.Mid = midValue.(int64)
	}

	c.JSON(svc.ContentList(c, arg))
}

func contentPost(c *bm.Context) {
	if conf.Conf.BulletConfig.CloseWrite {
		c.JSON(struct{}{}, nil)
		return
	}

	arg := &api.Bullet{}
	if err := c.Bind(arg); err != nil {
		return
	}

	// 这里客户端是有限制的，所以这里就不单独给出toast了
	if arg.Oid == 0 || arg.Content == "" || strings.Count(arg.Content, "") > model.BulletMaxLen {
		log.Warnw(c, "log", "post error", "req", arg, "content_len", strings.Count(arg.Content, ""))
		c.JSON(nil, ecode.DanmuPostErr)
		return
	}

	if midValue, exists := c.Get("mid"); exists {
		arg.Mid = midValue.(int64)
	}
	dmid, err := svc.ContentPost(c, arg)
	c.JSON(new(empty.Empty), err)

	uiLog(c, model.ActionDanmaku, struct {
		DMID int64 `json:"dmid"`
	}{
		DMID: dmid,
	})
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
	err := svc.PhoneCheck(ctx, mid)
	if err != nil {
		ctx.JSON(nil, err)
		ctx.Abort()
		return
	}
}
