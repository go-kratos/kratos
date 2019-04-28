package http

import (
	"go-common/app/service/bbq/topic/api"
	"go-common/library/ecode"
	"go-common/library/net/http/blademaster/binding"
	"net/http"

	"go-common/app/service/bbq/topic/internal/service"
	"go-common/library/conf/paladin"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	svc *service.Service
)

// New new a bm server.
func New(s *service.Service) (engine *bm.Engine) {
	var (
		hc struct {
			Server *bm.ServerConfig
		}
	)
	if err := paladin.Get("http.toml").UnmarshalTOML(&hc); err != nil {
		if err != paladin.ErrNotExist {
			panic(err)
		}
	}
	svc = s
	engine = bm.DefaultServer(hc.Server)
	initRouter(engine, verify.New(nil))
	if err := engine.Start(); err != nil {
		panic(err)
	}
	return
}

func initRouter(e *bm.Engine, v *verify.Verify) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/bbq/topic")
	{
		g.GET("/start", howToStart)
		g.POST("/update/state", updateTopicState)
		g.POST("/update/desc", updateTopicDesc)
		g.POST("/stick", stickTopic)
		g.POST("/video/stick", stickTopicVideo)
		g.POST("/video/set/stick", setStickTopicVideo)
		g.GET("/cms/list", cmsTopicList)
		g.GET("/video", topicVideo)
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

func howToStart(c *bm.Context) {
	c.String(0, "golang")
}

func updateTopicState(c *bm.Context) {
	arg := &api.TopicInfo{}
	if err := c.Bind(arg); err != nil {
		c.JSON(nil, ecode.ReqParamErr)
		return
	}
	c.JSON(svc.UpdateTopicState(c, arg))
}

func updateTopicDesc(c *bm.Context) {
	arg := &api.TopicInfo{}
	if err := c.Bind(arg); err != nil {
		c.JSON(nil, ecode.ReqParamErr)
		return
	}
	c.JSON(svc.UpdateTopicDesc(c, arg))
}

func cmsTopicList(c *bm.Context) {
	arg := &api.ListCmsTopicsReq{}
	if err := c.Bind(arg); err != nil {
		c.JSON(nil, ecode.ReqParamErr)
		return
	}
	c.JSON(svc.ListCmsTopics(c, arg))
}

func topicVideo(c *bm.Context) {
	arg := &api.ListCmsTopicsReq{}
	if err := c.Bind(arg); err != nil {
		c.JSON(nil, ecode.ReqParamErr)
		return
	}
	c.JSON(svc.ListCmsTopics(c, arg))
}

func stickTopic(c *bm.Context) {
	arg := &api.StickTopicReq{}
	if err := c.Bind(arg); err != nil {
		c.JSON(nil, ecode.ReqParamErr)
		return
	}
	c.JSON(svc.StickTopic(c, arg))
}

func stickTopicVideo(c *bm.Context) {
	arg := &api.StickTopicVideoReq{}
	if err := c.Bind(arg); err != nil {
		c.JSON(nil, ecode.ReqParamErr)
		return
	}
	c.JSON(svc.StickTopicVideo(c, arg))
}

func setStickTopicVideo(c *bm.Context) {
	arg := &api.SetStickTopicVideoReq{}
	if err := c.BindWith(arg, binding.JSON); err != nil {
		c.JSON(nil, ecode.ReqParamErr)
		return
	}
	c.JSON(svc.SetStickTopicVideo(c, arg))
}
