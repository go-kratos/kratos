package http

import (
	http2 "net/http"

	"go-common/app/admin/main/videoup-task/conf"
	"go-common/app/admin/main/videoup-task/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	srv     *service.Service
	permSrv *permit.Permit
	vfySvc  *verify.Verify
)

//Init init http
func Init(conf *conf.Config, s *service.Service) {
	srv = s
	permSrv = permit.New(conf.Auth)
	vfySvc = verify.New(nil)
	engine := bm.DefaultServer(conf.BM)
	innerRoute(engine)

	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func innerRoute(engine *bm.Engine) {
	engine.Ping(ping)
	g := engine.Group("/x/admin/vt")
	{
		v := g.Group("/video", permSrv.Permit("TASK_QA_VIDEO"))
		{
			v.GET("/list", list)
			v.GET("/detail", detail)
			v.POST("/submit", submit)
		}

		t := g.Group("/task", permSrv.Permit(""))
		{
			w := t.Group("/", permSrv.Permit("TASKWEIGHT"))
			{
				w.GET("/weightconfig/maxweight", maxweight)
				w.POST("/weightconfig/add", addwtconf)
				w.POST("/weightconfig/del", delwtconf)
				w.GET("/weightconfig/list", listwtconf)
				w.GET("/weightlog/list", listwtlog)
				w.GET("/wcv/show", show)
				w.POST("/wcv/set", set)
			}

			r := t.Group("/review")
			{
				r.GET("/config/list", listreviews)
				r.POST("/config/add", addreview)
				r.POST("/config/edit", editreview)
				r.POST("/config/delete", delreview)
			}
			c := t.Group("consumer")
			{
				c.GET("/on", checkgroup(), on)
				c.GET("/off", checkgroup(), off) //自己退出
				c.POST("/forceoff", forceoff)    //强制踢出
			}

			t.GET("/online", permSrv.Permit("ONLINE"), online)
			t.GET("/inoutlist", inoutlist)
			t.POST("/delay", checkowner(), delay)
			t.POST("/free", taskfree)
		}
	}

	g = engine.Group("/vt", vfySvc.Verify)
	{
		v := g.Group("/video")
		{
			v.POST("/add", add)
			v.POST("/uputime", upTaskUTime)
		}

		g.GET("/report/memberstats", memberStats)
		r := g.Group("review")
		{
			r.POST("/check", checkReview)
		}

		t := g.Group("task")
		{
			t.GET("/tooks", taskTooks)
			t.GET("/next", next)
			t.GET("/list", listTask)
			t.GET("/info", info)

		}
	}
}

func ping(ctx *bm.Context) {
	if srv.Ping(ctx) != nil {
		ctx.AbortWithStatus(http2.StatusServiceUnavailable)
		ctx.Done()
	}
}

func getUIDName(ctx *bm.Context) (uid int64, username string) {
	if uidi, _ := ctx.Get("uid"); uidi != nil {
		uid = uidi.(int64)
	}
	if name, _ := ctx.Get("username"); name != nil {
		username = name.(string)
	}
	return
}
