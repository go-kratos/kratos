package http

import (
	"net/http"

	"go-common/app/admin/main/bfs/conf"
	"go-common/app/admin/main/bfs/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	srv *service.Service
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	srv = s
	engine := bm.DefaultServer(c.BM)
	router(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func router(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	group := e.Group("/x/admin/bfs")
	{
		group.GET("/clusters", clusters)
		group.GET("/total", bfsTotal)
		group.GET("/rack", rackMeta)
		group.GET("/volume", volumeMeta)
		group.GET("/group", groupMeta)
		group.POST("/group/compact", compact)
		group.POST("/group/status", setGroupStatus)
		group.POST("/group/add_volume", addVolume)
		group.POST("/group/add_free_volume", addFreeVolume)
	}
}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("bfs ping error(%v)", err)
		c.Error = err
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
