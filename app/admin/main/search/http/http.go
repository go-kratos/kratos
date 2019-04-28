package http

import (
	"net/http"

	"go-common/app/admin/main/search/conf"
	"go-common/app/admin/main/search/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	authSrv *permit.Permit
	svr     *service.Service
)

// Init init http
func Init(c *conf.Config, s *service.Service) {
	svr = s
	authSrv = permit.New(c.Auth)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	searchG := e.Group("/x/admin/search")
	{
		// V3新版查询和更新接口
		searchG.GET("/query", querySearch)
		searchG.GET("/query/debug", queryDebug)
		searchG.POST("/upsert", upsert)
		// V2老接口
		searchG.GET("/archive", archiveSearch)
		searchG.GET("/log", logSearch)
		searchG.POST("/log/delete", logDelete)
		searchG.GET("/log/audit", bMlogAudit)
		searchG.GET("/log/audit_group", bMlogAuditGroupBy)
		searchG.GET("/log/user_action", bMlogUserAction)
		// update (deprecated)
		searchG.POST("/archive/update", updateArchive)
		// index insert (deprecated)
		searchG.POST("/copyright/index", copyRight)

		// sven
		mng := searchG.Group("/mng")
		{
			mng.GET("/business/list", businessList)
			mng.GET("/business/all", businessAll)
			mng.GET("/business/info", businessInfo)
			mng.POST("/business/add", addBusiness)
			mng.POST("/business/update", updateBusiness)
			mng.POST("/business/update_app", updateBusinessApp)

			mng.GET("/asset/list", assetList)
			mng.GET("/asset/all", assetAll)
			mng.GET("/asset/info", assetInfo)
			mng.POST("/asset/add", addAsset)
			mng.POST("/asset/update", updateAsset)

			mng.GET("/app/list", appList)
			mng.GET("/app/info", appInfo)
			mng.POST("/app/add", addApp)
			mng.POST("/app/update", updateApp)

			mng.GET("/countlist", countlist)
			mng.GET("/count", count)
			mng.GET("/percent", percent)
		}
		// sven v2
		mng2 := searchG.Group("/mng/v2")
		{
			mng2.GET("/business/all", businessAllV2)
			mng2.GET("/business/info", businessInfoV2)
			mng2.POST("/business/add", businessAdd)
			mng2.POST("/business/update", businessUpdate)

			mng2.GET("/asset/list", assetDBTables)
			mng2.GET("/asset/info", assetInfoV2)
			mng2.GET("/asset/dbconnect", assetDBConnect)
			mng2.POST("/asset/dbadd", assetDBAdd)
			mng2.POST("/asset/tableadd", assetTableAdd)
			mng2.POST("/asset/tableupdate", updateAssetTable)
			mng2.GET("/asset/showtables", assetShowTables)
			mng2.GET("/asset/tablefields", assetTableFields)

			mng2.GET("/cluster/owners", clusterOwners)
		}
	}
}

// ping check health
func ping(ctx *bm.Context) {
	if err := svr.Ping(ctx); err != nil {
		ctx.Error = err
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
