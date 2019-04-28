package http

import (
	"net/http"

	"go-common/app/admin/main/esports/conf"
	"go-common/app/admin/main/esports/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	esSvc *service.Service
	//idfSvc  *identify.Identify
	permitSvc *permit.Permit
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	esSvc = s
	permitSvc = permit.New(c.Permit)
	engine := bm.DefaultServer(c.BM)
	authRouter(engine)
	// init internal server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func authRouter(e *bm.Engine) {
	e.Ping(ping)
	group := e.Group("/x/admin/esports", permitSvc.Permit("ESPORTS_ADMIN"))
	{
		matchGroup := group.Group("/matchs")
		{
			matchGroup.GET("/info", matchInfo)
			matchGroup.GET("/list", matchList)
			matchGroup.POST("/add", addMatch)
			matchGroup.POST("/save", editMatch)
			matchGroup.POST("/forbid", forbidMatch)
		}
		seasonGroup := group.Group("/seasons")
		{
			seasonGroup.GET("/info", seasonInfo)
			seasonGroup.GET("/list", seasonList)
			seasonGroup.POST("/add", addSeason)
			seasonGroup.POST("/save", editSeason)
			seasonGroup.POST("/forbid", forbidSeason)
		}
		contestGroup := group.Group("/contest")
		{
			contestGroup.GET("/info", contestInfo)
			contestGroup.GET("/list", contestList)
			contestGroup.POST("/add", addContest)
			contestGroup.POST("/save", editContest)
			contestGroup.POST("/forbid", forbidContest)
		}
		gameGroup := group.Group("/games")
		{
			gameGroup.GET("/info", gameInfo)
			gameGroup.GET("/list", gameList)
			gameGroup.POST("/add", addGame)
			gameGroup.POST("/save", editGame)
			gameGroup.POST("/forbid", forbidGame)
			gameGroup.GET("/types", types)
		}
		teamGroup := group.Group("/teams")
		{
			teamGroup.GET("/info", teamInfo)
			teamGroup.GET("/list", teamList)
			teamGroup.POST("/add", addTeam)
			teamGroup.POST("/save", editTeam)
			teamGroup.POST("/forbid", forbidTeam)
		}
		tagGroup := group.Group("/tags")
		{
			tagGroup.GET("/info", tagInfo)
			tagGroup.GET("/list", tagList)
			tagGroup.POST("/add", addTag)
			tagGroup.POST("/save", editTag)
			tagGroup.POST("/forbid", forbidTag)
		}
		arcGroup := group.Group("/arcs")
		{
			arcGroup.GET("/list", arcList)
			arcGroup.POST("/edit", editArc)
			arcGroup.POST("/batch/add", batchAddArc)
			arcGroup.POST("/batch/edit", batchEditArc)
			arcGroup.POST("/batch/del", batchDelArc)
			arcGroup.POST("/import/csv", arcImportCSV)
		}
		actGroup := group.Group("/active")
		{
			actGroup.GET("", listAct)
			actGroup.POST("/add", addAct)
			actGroup.POST("/edit", editAct)
			actGroup.POST("/forbid", forbidAct)
			dGroup := actGroup.Group("/detail")
			{
				dGroup.GET("/list", listDetail)
				dGroup.POST("/add", addDetail)
				dGroup.POST("/edit", editDetail)
				dGroup.POST("/forbid", forbidDetail)
				dGroup.POST("/online", onLine)
			}
			tGroup := actGroup.Group("/tree")
			{
				tGroup.GET("/list", listTree)
				tGroup.POST("/add", addTree)
				tGroup.POST("/edit", editTree)
				tGroup.POST("/del", delTree)
			}
		}
	}
}

func ping(c *bm.Context) {
	if err := esSvc.Ping(c); err != nil {
		log.Error("esport-admin ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
