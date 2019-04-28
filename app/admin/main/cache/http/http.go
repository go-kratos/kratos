package http

import (
	"net/http"

	"go-common/app/admin/main/cache/conf"
	"go-common/app/admin/main/cache/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	srv *service.Service
)

// Init init
func Init(c *conf.Config) {
	srv = service.New(c)
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
	g := e.Group("/x/admin/cache")
	{
		g.GET("/clusters", clusters)
		g.GET("/cluster", cluster)
		g.GET("/cluster/detail", clusterDtl)
		g.POST("/cluster/add", addCluster)
		g.POST("/cluster/del", delCluster)
		g.POST("/cluster/node/modify", modifyCluster)
		g.GET("/cluster/toml", toml)
		g.POST("/cluster/from/yml", addFromYml)
	}
	ol := e.Group("/x/admin/cache/overlord")
	{
		ol.GET("/clusters", overlordClusters)
		ol.POST("/del/cluster", overlordDelCluster)
		ol.POST("/del/node", overlordDelNode)
		ol.GET("/ops/names", overlordOpsClusterNames)
		ol.GET("/ops/nodes", overlordOpsNodes)
		ol.POST("/import/ops/cluster", overlordImportCluster)
		ol.POST("/new/ops/node", overlordClusterNewNode)
		ol.POST("/replace/ops/node", overlordClusterReplaceNode)
		ol.GET("/app/clusters", overlordAppClusters)
		ol.GET("/app/can/bind/clusters", overlordAppNeedClusters)
		ol.POST("/app/cluster/bind", overlordAppClusterBind)
		ol.POST("/app/cluster/del", overlordAppClusterDel)
		ol.GET("/app/appids", overlordAppAppIDs)
		ol.GET("/app/toml", overlordToml)
	}
}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("cache-admin ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
