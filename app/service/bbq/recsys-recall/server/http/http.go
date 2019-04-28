package http

import (
	"fmt"
	"net/http"

	grpc "go-common/app/service/bbq/recsys-recall/api/grpc/v1"
	"go-common/app/service/bbq/recsys-recall/conf"
	"go-common/app/service/bbq/recsys-recall/service"
	"go-common/app/service/bbq/recsys-recall/service/index"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
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
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/bbq/internal/recall")
	{
		g.GET("/start", vfy.Verify, howToStart)
		g.GET("/forward_index", forwardIndex)
		g.GET("/inverted_index", invertedIndex)
		g.GET("/recall", recall)
		g.GET("/videos", videosByIndex)
		g.POST("/new_income", newIncomeVideos)
	}
}

func forwardIndex(c *bm.Context) {
	args := struct {
		Svid uint64 `form:"svid" json:"svid" validate:"required"`
	}{}

	var err error
	if err = c.Bind(&args); err != nil {
		log.Errorv(*c, log.KV("log", err))
		return
	}
	if res := index.Index.Get(args.Svid); res != nil {
		c.String(0, res.String())
		return
	}

	c.String(0, "error: %v", err)
}

func invertedIndex(c *bm.Context) {
	c.JSON(nil, nil)
}

func recall(c *bm.Context) {
	args := struct {
		MID        int64  `json:"mid" form:"mid"`
		BUVID      string `json:"buvid" form:"buvid"`
		TotalLimit int32  `json:"total_limit" form:"total_limit"`
		Tag        string `json:"tag" form:"tag"`
		Name       string `json:"name" form:"name"`
		Scorer     string `json:"scorer" form:"scorer"`
		Filter     string `json:"filter" form:"filter"`
		Ranker     string `json:"ranker" form:"ranker"`
		Priority   int32  `json:"priority" form:"priority"`
		Limit      int32  `json:"limit" form:"limit"`
	}{}

	if err := c.Bind(&args); err != nil {
		return
	}

	req := &grpc.RecallRequest{
		MID:        args.MID,
		BUVID:      args.BUVID,
		TotalLimit: args.TotalLimit,
		Infos: []*grpc.RecallInfo{
			{
				Tag:      args.Tag,
				Name:     args.Name,
				Scorer:   args.Scorer,
				Filter:   args.Filter,
				Ranker:   args.Ranker,
				Priority: args.Priority,
				Limit:    args.Limit,
			},
		},
	}
	c.JSON(srv.Recall(c, req))
}

func newIncomeVideos(c *bm.Context) {
	args := &grpc.NewIncomeVideoRequest{}
	if err := c.BindWith(args, binding.JSON); err != nil {
		return
	}
	fmt.Println(args)
	c.JSON(srv.NewIncomeVideo(c, args))
}

func videosByIndex(c *bm.Context) {
	args := &grpc.VideosByIndexRequest{}
	if err := c.Bind(args); err != nil {
		return
	}
	c.JSON(srv.VideosByIndex(c, args))
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

// example for http request handler
func howToStart(c *bm.Context) {
	c.String(0, "Golang 大法好 !!!")
}
