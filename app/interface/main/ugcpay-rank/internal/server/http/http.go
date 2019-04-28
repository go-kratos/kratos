package http

import (
	"go-common/app/interface/main/ugcpay-rank/internal/conf"
	"go-common/app/interface/main/ugcpay-rank/internal/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
)

var (
	svc   *service.Service
	authM *auth.Auth
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	svc = s
	authM = auth.New(conf.Conf.Auth)
	engine := bm.DefaultServer(conf.Conf.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/x/ugcpay-rank")
	{
		g1 := g.Group("/v1/elec", authM.GuestMobile)
		{
			g1.GET("/month/up", elecMonthUP)
			g1.GET("/month", elecMonth)
			g1.GET("/all/av", elecAllAV)
		}
		g1 = g.Group("/elec", authM.GuestWeb)
		{
			g1.GET("/month/up", elecMonthUP)
			g1.GET("/month", elecMonth)
			g1.GET("/all/av", elecAllAV)
		}
	}
}

func ping(ctx *bm.Context) {
}
