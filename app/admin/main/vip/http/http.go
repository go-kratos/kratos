package http

import (
	"go-common/app/admin/main/vip/conf"
	"go-common/app/admin/main/vip/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	// depend service
	vipSvc  *service.Service
	authSvc *permit.Permit
	cf      *conf.Config
)

//Init init http
func Init(c *conf.Config) {
	cf = c
	initService(c)
	// init external router
	engine := bm.DefaultServer(c.BM)
	initRouter(engine)
	// init Outer serve
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	vipSvc = service.New(c)
	authSvc = permit.New(c.Auth)
}

// initRouter init outer router api path.
func initRouter(e *bm.Engine) {
	e.Ping(moPing)
	group := e.Group("/x/admin/vip", authSvc.Verify())
	{
		monthGroup := group.Group("/month", authSvc.Permit("VIP_MONTH"))
		{
			monthGroup.GET("/list", monthList)
			monthGroup.POST("/edit", monthEdit)
			monthGroup.GET("/price/list", priceList)
			monthGroup.POST("/price/add", priceAdd)
			monthGroup.POST("/price/edit", priceEdit)
		}
		poolGroup := group.Group("/pool", authSvc.Permit("VIP_POOL"))
		{
			poolGroup.GET("/list", queryPool)
			poolGroup.GET("/info", getPool)
			poolGroup.POST("/save", savePool)

			batchGroup := poolGroup.Group("/batch")
			{
				batchGroup.GET("/list", queryBatch)
				batchGroup.GET("/info", getBatch)
				batchGroup.POST("/add", addBatch)
				batchGroup.POST("/edit", saveBatch)
				batchGroup.POST("/consume", grantResouce)
			}

		}
		batchCodeGroup := group.Group("/batchCode", authSvc.Permit("VIP_BATCH_CODE"))
		{
			batchCodeGroup.GET("/list", batchCodes)
			batchCodeGroup.POST("/save", saveBatchCode)
			batchCodeGroup.POST("/frozen", frozenBatchCode)
			batchCodeGroup.GET("/export", exportCodes)
			codeGroup := batchCodeGroup.Group("/code")
			{
				codeGroup.GET("/list", codes)
				codeGroup.POST("/frozen", frozenCode)
			}
		}
		pushGroup := group.Group("/push", authSvc.Permit("VIP_PUSH"))
		{
			pushGroup.GET("/list", pushs)
			pushGroup.POST("/save", savePush)
			pushGroup.GET("/info", push)
			pushGroup.GET("/del", delPush)
			pushGroup.GET("/disable", disablePush)

		}

		vipGroup := group.Group("/user", authSvc.Permit("VIP_USER"))
		{
			vipGroup.POST("/drawback", drawback)
			vipGroup.GET("/log/list", historyList)
			vipGroup.GET("/info", vipInfo)
		}
		bizGroup := group.Group("/biz", authSvc.Permit("VIP_BIZ"))
		{
			bizGroup.GET("/list", businessList)
			bizGroup.POST("/add", addBusiness)
			bizGroup.GET("/info", business)
			bizGroup.POST("/edit", updateBusiness)
		}
		verGroup := group.Group("/version", authSvc.Permit("VIP_VERSION"))
		{
			verGroup.GET("/list", versions)
			verGroup.POST("/edit", updateVersion)
		}
		tipsGroup := group.Group("/tips", authSvc.Permit("VIP_TIPS"))
		{
			tipsGroup.GET("/list", tips)
			tipsGroup.GET("/info", tipbyid)
			tipsGroup.POST("/add", tipadd)
			tipsGroup.POST("/edit", tipupdate)
			tipsGroup.GET("/delete", tipdelete)
			tipsGroup.POST("/expire", tipexpire)
		}
		panelGroup := group.Group("/panel")
		{
			panelGroup.GET("/conf/types", authSvc.Permit("VIP_PRICE_PANEL"), vipPanelTypes)
			panelGroup.GET("/conf/list", authSvc.Permit("VIP_PRICE_PANEL"), vipPriceConfigs)
			panelGroup.GET("/conf/info", authSvc.Permit("VIP_PRICE_PANEL"), vipPriceConfigID)
			panelGroup.POST("/conf/add", authSvc.Permit("VIP_PRICE_PANEL_CONTROL"), addVipPriceConfig)
			panelGroup.POST("/conf/up", authSvc.Permit("VIP_PRICE_PANEL_CONTROL"), upVipPriceConfig)
			panelGroup.POST("/conf/del", authSvc.Permit("VIP_PRICE_PANEL_CONTROL"), delVipPriceConfig)
			panelGroup.GET("/conf/dprice/list", authSvc.Permit("VIP_PRICE_PANEL"), vipDPriceConfigs)
			panelGroup.GET("/conf/dprice/info", authSvc.Permit("VIP_PRICE_PANEL"), vipDPriceConfigID)
			panelGroup.POST("/conf/dprice/add", authSvc.Permit("VIP_PRICE_PANEL_CONTROL"), addVipDPriceConfig)
			panelGroup.POST("/conf/dprice/up", authSvc.Permit("VIP_PRICE_PANEL_CONTROL"), upVipDPriceConfig)
			panelGroup.POST("/conf/dprice/del", authSvc.Permit("VIP_PRICE_PANEL_CONTROL"), delVipDPriceConfig)
		}
		pgGroup := group.Group("/privilege", authSvc.Permit("VIP_PRIVILEGE"))
		{
			pgGroup.GET("/list", privileges)
			pgGroup.POST("/update/state", updatePrivilegeState)
			pgGroup.POST("/delete", deletePrivilege)
			pgGroup.POST("/update/order", updateOrder)
			pgGroup.POST("/add", addPrivilege)
			pgGroup.POST("/modify", updatePrivilege)
		}
		jointlyGroup := group.Group("/jointly", authSvc.Permit("VIP_JOINTLY"))
		{
			jointlyGroup.GET("/list", jointlys)
			jointlyGroup.POST("/add", addJointly)
			jointlyGroup.POST("/modify", modifyJointly)
			jointlyGroup.POST("/delete", deleteJointly)
		}
		refundGroup := group.Group("/order", authSvc.Permit("VIP_ORDER"))
		{
			refundGroup.GET("/list", orderList)
			refundGroup.POST("/refund", authSvc.Permit("VIP_REFUND"), refund)
		}
		dialogGroup := group.Group("/dialog", authSvc.Permit("VIP_DIALOG"))
		{
			dialogGroup.GET("/list", dialogList)
			dialogGroup.GET("/info", dialogInfo)
			dialogGroup.POST("/save", dialogSave)
			dialogGroup.POST("/enable", dialogEnable)
			dialogGroup.POST("/del", dialogDel)
		}
		platformGroup := group.Group("/platform", authSvc.Permit("VIP_PLATFORM"))
		{
			platformGroup.GET("/list", platformList)
			platformGroup.GET("/info", platformInfo)
			platformGroup.POST("/save", platformSave)
			platformGroup.POST("/del", platformDel)
		}
		welfareGroup := group.Group("/welfare", authSvc.Permit("VIP_WELFARE"))
		{
			welfareGroup.POST("/type/save", welfareTypeSave)
			welfareGroup.POST("/type/state", welfareTypeState)
			welfareGroup.GET("/type/list", welfareTypeList)
			welfareGroup.POST("/save", welfareSave)
			welfareGroup.POST("/state", welfareState)
			welfareGroup.GET("/list", welfareList)
			welfareGroup.POST("/batch/upload", welfareBatchUpload)
			welfareGroup.GET("/batch/list", welfareBatchList)
			welfareGroup.POST("/batch/state", welfareBatchState)
		}
	}
}

func moPing(c *bm.Context) {
	vipSvc.Ping(c)
}
