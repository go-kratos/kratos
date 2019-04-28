package http

import (
	"encoding/json"
	"net/http"

	"go-common/app/service/bbq/video/conf"
	"go-common/app/service/bbq/video/service"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	srv *service.Service
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	srv = s
	engine := bm.DefaultServer(c.BM.Server)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/bbq/internal")
	{
		g.POST("/sv/trans/back", bvcTransBack)
		g.GET("/sv/trans/commit", bvcTransCommit)
		g.POST("/sv/stat", videoStat)
		g.GET("/create/id", createID)
		g.GET("/sv/create/id", createID)
		// 增加视频播放量数据
		g.POST("/sv/views/add", videoViewsAdd)
		g.POST("/sv/limits/modify", limitsModify)
		g.POST("/sv/play", svPlays)
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

// bindJson 解析application/json请求body
func bindJSON(c *bm.Context, obj interface{}) error {
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(obj)
	if err != nil {
		log.Warn("参数解析失败 %v", err)
		err = ecode.ReqParamErr
		return err
	}
	return nil
}
