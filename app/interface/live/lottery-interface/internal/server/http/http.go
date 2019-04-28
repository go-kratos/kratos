package http

import (
	"net/http"

	"go-common/app/interface/live/lottery-interface/internal/conf"
	"go-common/app/interface/live/lottery-interface/internal/service"
	v1 "go-common/app/interface/live/lottery-interface/internal/service/v1"
	risk "go-common/app/service/live/live_riskcontrol/api/grpc/v1"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/metadata"

	"go-common/library/net/http/blademaster/render"
)

var (
	vfy *verify.Verify
)

// Init init
func Init(c *conf.Config) {

	vfy = verify.New(c.Verify)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	auth := auth.New(&auth.Config{
		Identify:    conf.Conf.ShortClient,
		DisableCSRF: true,
	})

	g := e.Group("/xlive/lottery-interface")
	{
		g.POST("/v1/storm/Join", abortGET, auth.Guest, riskControl, v1.StormJoin)
		g.GET("/v1/storm/Check", auth.Guest, v1.StormCheck)
	}
}

func abortGET(c *bm.Context) {
	if c.Request.Method == http.MethodGet {
		c.Render(http.StatusOK, render.MapJSON{
			"code": 0,
			"msg":  "",
			"data": []int{},
		})
		c.Abort()
		return
	}
}
func riskControl(c *bm.Context) {
	mid, ok := c.Get("mid")
	if !ok {
		c.Render(http.StatusOK, render.MapJSON{
			"code": 401,
			"msg":  "未登录",
			"data": []int{},
		})
		c.Abort()
		return
	}
	midi := mid.(int64)
	header := make(map[string]string)
	for k := range c.Request.Header {
		if v := c.Request.Header[k]; len(v) > 0 {
			header[k] = v[0]
		}

	}
	header["Platform"] = c.Request.URL.Query().Get("platform")
	resp, err := service.ServiceInstance.IsForbiddenClient.GetForbidden(c, &risk.GetForbiddenReq{
		Uid:    midi,
		Uri:    c.Request.URL.Path,
		Ip:     metadata.String(c, metadata.RemoteIP),
		Method: c.Request.Method,
		Header: header,
	})
	if err != nil {
		log.Info("#call_IsForbiddenClient_fail %s", err.Error())
		c.Render(http.StatusOK, render.MapJSON{
			"code": 400,
			"msg":  "没抢到",
			"data": []int{},
		})
		c.Abort()
		return
	}
	if resp.GetIsForbidden() == risk.GetForbiddenReply_FORBIDDEN {
		log.Info("#IsForbidden# mid= %d", midi)
		c.Render(http.StatusOK, render.MapJSON{
			"code": 400,
			"msg":  "访问被拒绝",
			"data": []int{},
		})
		c.Abort()
		return
	}
}

func ping(ctx *bm.Context) {
	ctx.String(http.StatusOK, "PONG")
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
