package http

import (
	"net/http"

	"go-common/app/interface/main/player/conf"
	"go-common/app/interface/main/player/service"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	authSvr   *auth.Auth
	vfySvr    *verify.Verify
	playSvr   *service.Service
	playInfoc *infoc.Infoc
)

// Init init http.
func Init(c *conf.Config) error {
	authSvr = auth.New(c.Auth)
	vfySvr = verify.New(c.Verify)
	playSvr = service.New(c)
	engine := bm.DefaultServer(c.BM.Outer)
	outerRouter(engine)
	internalRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		return err
	}
	// init infoc
	if c.Infoc2 != nil {
		playInfoc = infoc.New(c.Infoc2)
	}
	return nil
}

func outerRouter(e *bm.Engine) {
	e.GET("/monitor/ping", ping)
	e.GET("/x/player.so", bm.CORS(), authSvr.Guest, player)
	group := e.Group("/x/player", bm.CORS())
	{
		group.GET("/policy", authSvr.Guest, policy)
		group.GET("/carousel.so", carousel)
		group.GET("/view", view)
		group.GET("/matsuri", matPage)
		group.GET("/pagelist", pageList)
		group.GET("/videoshot", videoShot)
		group.GET("/playurl/token", authSvr.User, playURLToken)
		group.GET("/playurl", authSvr.Guest, playurl)
	}
}

func internalRouter(e *bm.Engine) {
	group := e.Group("/x/internal/player")
	{
		group.GET("/playurl", vfySvr.Verify, authSvr.Guest, playurl)
	}
}

func ping(c *bm.Context) {
	if err := playSvr.Ping(c); err != nil {
		log.Error("player service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
