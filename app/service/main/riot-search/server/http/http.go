package http

import (
	"net/http"

	"go-common/app/service/main/riot-search/conf"
	"go-common/app/service/main/riot-search/model"
	"go-common/app/service/main/riot-search/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	srv *service.Service
	vfy *verify.Verify
)

// Init init
func Init(c *conf.Config) {
	srv = service.New(c)
	vfy = verify.New(c.Verify)
	engine := bm.DefaultServer(c.BM)
	router(engine)
	if err := engine.Start(); err != nil {
		log.Error("xhttp.Serve error(%v)", err)
		panic(err)
	}
}

func router(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/x/internal/riot-search")
	{
		g.POST("/arc/ids", vfy.Verify, searchIDOnly)
		g.POST("/arc/contents", vfy.Verify, search)
		// debug api
		g.GET("/arc/has", has)
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

// @params RiotSearchReq
// @router post /x/riot-search/arc/aids
// @response IDsResp
func searchIDOnly(c *bm.Context) {
	req := new(model.RiotSearchReq)
	if err := c.Bind(req); err != nil {
		log.Error("request param(%v) error", req)
		return
	}
	c.JSON(srv.SearchIDOnly(c, req), nil)
}

// @params RiotSearchReq
// @router post /x/riot-search/arc/contents
// @response DocumentsResp
func search(c *bm.Context) {
	req := new(model.RiotSearchReq)
	if err := c.Bind(req); err != nil {
		log.Error("request param(%v) error", req)
		return
	}
	c.JSON(srv.Search(c, req), nil)
}

func has(c *bm.Context) {
	req := new(struct {
		ID uint64 `form:"id" validate:"min=0"`
	})
	if err := c.Bind(req); err != nil {
		return
	}
	c.JSON(srv.Has(c, req.ID), nil)
}
