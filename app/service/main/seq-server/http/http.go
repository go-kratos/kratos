package http

import (
	"net/http"

	"go-common/app/service/main/seq-server/conf"
	"go-common/app/service/main/seq-server/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	seqSvc *service.Service
	vrfSvc *verify.Verify
)

// Init init http router.
func Init(c *conf.Config, s *service.Service) {
	seqSvc = s
	vrfSvc = verify.New(nil)
	en := bm.DefaultServer(c.BM)
	innerRouter(en)
	// init internal server
	if err := en.Start(); err != nil {
		log.Error("xhttp.Serve error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	seq := e.Group("/x/internal/seq", vrfSvc.Verify)
	{
		seq.POST("/id", seqID)
		seq.POST("/id32", seqID32)
		seq.POST("/maxseq", maxSeq)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := seqSvc.Ping(c); err != nil {
		log.Error("seq-server ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// register check server ok.
func register(c *bm.Context) {
	c.JSON(nil, nil)
}
