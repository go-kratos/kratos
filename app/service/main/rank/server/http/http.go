package http

import (
	"net/http"
	"strings"

	"go-common/app/service/main/rank/conf"
	"go-common/app/service/main/rank/model"
	"go-common/app/service/main/rank/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	srv *service.Service
	vfy *verify.Verify
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	srv = s
	vfy = verify.New(c.Verify)
	engine := bm.DefaultServer(c.BM)
	router(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start() error(%v)", err)
		panic(err)
	}
}

func router(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/x/internal/rank")
	{
		g.GET("/do", do)
		g.GET("/mget", mget)
		g.GET("/sort", sort)
		g.GET("/group", vfy.Verify, group)
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

func do(c *bm.Context) {
	arg := new(model.DoReq)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.Do(c, arg), nil)
}

func mget(c *bm.Context) {
	arg := new(model.MgetReq)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.Mget(c, arg))
}

func sort(c *bm.Context) {
	arg := new(struct {
		Business string   `form:"business" validate:"required"`
		Field    string   `form:"field" validate:"required"`
		Order    string   `form:"order" validate:"required"`
		Filters  []string `form:"filters,split"`
		Oids     []int64  `form:"oids,split" validate:"required"`
		Pn       int      `form:"pn"`
		Ps       int      `form:"ps"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	filterMap := make(map[string]string)
	for _, v := range arg.Filters {
		strs := strings.Split(v, "|")
		if len(strs) == 2 {
			filterMap[strs[0]] = strs[1]
		}
	}
	a := new(model.SortReq)
	a.Business = arg.Business
	a.Field = arg.Field
	a.Order = arg.Order
	a.Filters = filterMap
	a.Oids = arg.Oids
	a.Pn = arg.Pn
	a.Ps = arg.Ps
	c.JSON(srv.Sort(c, a))
}

func group(c *bm.Context) {
	arg := new(model.GroupReq)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.Group(c, arg))
}
