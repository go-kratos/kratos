package http

import (
	"encoding/json"
	"net/http"

	"go-common/app/service/main/account/conf"
	"go-common/app/service/main/account/model"
	"go-common/app/service/main/account/service"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	v "go-common/library/net/http/blademaster/middleware/verify"
)

var (
	accSvc *service.Service
	verify *v.Verify
)

// Init init account service.
func Init(c *conf.Config, s *service.Service) {
	accSvc = s
	verify = v.New(c.Verify)
	// engine
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	// init inner server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start() error(%v)", err)
		panic(err)
	}
}

func filterByAppkey(keys []string) func(*bm.Context) {
	allowed := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		allowed[k] = struct{}{}
	}
	return func(ctx *bm.Context) {
		req := ctx.Request
		params := req.Form
		appkey := params.Get("appkey")
		if _, ok := allowed[appkey]; !ok {
			log.Error("appkey: %s try to access %s failed", appkey, req.URL)
			ctx.JSON(nil, ecode.AccessDenied)
			ctx.Abort()
			return
		}
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	group := e.Group("/x/internal/v3/account", verify.Verify)
	{
		group.GET("/info", info)
		group.GET("/info/by/name", infoByName)
		group.GET("/infos", infos)
		group.GET("/card", card)
		group.GET("/cards", cards)
		group.GET("/vip", vip)
		group.GET("/vips", vips)
		group.GET("/profile", profile)
		group.GET("/profile/stat", profileWithStat)
		group.GET("/privacy", filterByAppkey(conf.Conf.AppkeyFilter.Privacy), privacy)
		group.GET("/cache/del", cacheDel)
		group.POST("/cache/clear", cacheClear)
	}
	v2Group := e.Group("/x/internal/account/v2", verify.Verify)
	{
		v2Group.GET("/myinfo", v2MyInfo)
		v2Group.GET("/userinfo", v2MyInfo)
	}
	v1Group := e.Group("/x/internal/account", verify.Verify)
	{
		v1Group.GET("/info", v1Info)
		v1Group.GET("/infos", v1Infos)
		v1Group.GET("/card", v1Card)
		v1Group.GET("/vip", v1Vip)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := accSvc.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// register support discovery.
func register(c *bm.Context) {
	c.JSON(map[string]struct{}{}, nil)
}

// cache del
func cacheDel(c *bm.Context) {
	p := new(model.ParamModify)
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(nil, accSvc.DelCache(c, p.Mid, p.ModifiedAttr))
}

// cache clear
func cacheClear(c *bm.Context) {
	p := new(model.ParamMsg)
	if err := c.Bind(p); err != nil {
		return
	}
	var m struct {
		New struct {
			Mid int64 `json:"mid"`
		} `json:"new,omitempty"`
		Mid    int64  `json:"mid"`
		Action string `json:"action"`
	}
	if err := json.Unmarshal([]byte(p.Msg), &m); err != nil {
		c.JSON(nil, err)
		return
	}
	mid := m.Mid
	if mid == 0 {
		mid = m.New.Mid
	}
	if mid == 0 {
		log.Warn("cache clear no mid msg(%s)", p.Msg)
		return
	}
	log.Info("Try to delete cache with mid: %d and param: %+v", mid, p)
	c.JSON(nil, accSvc.DelCache(c, mid, m.Action))
}
