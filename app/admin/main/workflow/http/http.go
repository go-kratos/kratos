package http

import (
	"go-common/app/admin/main/workflow/service"
	"go-common/library/conf/paladin"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	authSvc *permit.Permit
	wkfSvc  *service.Service
)

// Init http server
func Init(s *service.Service) {
	var (
		hc struct {
			BM     *bm.ServerConfig
			Permit *permit.Config
		}
	)
	if err := paladin.Get("http.toml").UnmarshalTOML(&hc); err != nil {
		panic(err)
	}
	// init service
	iniService(hc.Permit, s)
	// init internal router
	engine := bm.DefaultServer(hc.BM)
	//global timeout
	setupInnerEngine(engine)
	// init internal server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func iniService(auth *permit.Config, s *service.Service) {
	authSvc = permit.New(auth)
	wkfSvc = s
}

// innerRouter
func setupInnerEngine(e *bm.Engine) {
	// monitor ping
	e.Ping(ping)

	workflow := e.Group("/x/admin/workflow")

	// platform
	platform := workflow.Group("/platform", authSvc.Permit(""))
	{
		platform.GET("/count", platformChallCount)
		platform.GET("/list/pending", platformChallListPending)
		platform.GET("/list/handling", platformHandlingChalllist)
		platform.GET("/list/done", platformDoneChallList)
		platform.GET("/list/created", platformCreatedChallList)
		platform.GET("/release", platformRelease)
		platform.GET("/checkin", platformCheckIn)
	}

	// challenge
	challenge := workflow.Group("/challenge")
	{
		challenge.GET("/list", challList)
		challenge.GET("/list2", authSvc.Permit(""), challListCommon)
		challenge.GET("/detail", challDetail)
		challenge.GET("/activity/list", listChallActivity)
		challenge.GET("/business/list", listChallBusiness) // Deprecated
		challenge.POST("/update", authSvc.Permit(""), upChall)
		challenge.POST("/reset", authSvc.Permit(""), rstChallResult)
		challenge.POST("/state/set", authSvc.Permit(""), setChallResult)             // new api logic
		challenge.POST("/state/batch/set", authSvc.Permit(""), batchSetChallResult)  //new api logic
		challenge.POST("/batch/result/set", authSvc.Permit(""), batchSetChallResult) // todo: deprecated
		challenge.POST("/extra/update", authSvc.Permit(""), upChallExtra)
		challenge.POST("/extra/batch/update", authSvc.Permit(""), batchUpChallExtra)
		challenge.POST("/business/state/set", authSvc.Permit(""), upChallBusState)
		challenge.POST("/business/state/batch/set", authSvc.Permit(""), batchUpChallBusState)

		// manager-v4 used
		challenge.POST("/business/busState/update", upBusChallsBusState) //Deprecated
		challenge.POST("/reset/business/state/batch/set", upBusChallsBusState)

		// challenge event
		event := challenge.Group("/event")
		event.POST("/add", addEvent)
		event.POST("/batch/add", batchAddEvent)
		event.GET("/list", eventList)

		// call reply/add sync add event and set business_state
		reply := challenge.Group("/reply")
		reply.POST("/add", authSvc.Permit(""), addReply)
		reply.POST("/batch/add", authSvc.Permit(""), batchAddReply)
	}

	// business
	business := workflow.Group("/business")
	{
		business.GET("/meta/list", busMetaList)

		// callback
		callback := business.Group("/callback")
		{
			callback.GET("/list", listCallback)
			callback.POST("/add", authSvc.Permit(""), addOrUpCallback)
			callback.POST("/external/api/set", setExtAPI)
		}

		// attr
		busAttr := business.Group("/attr")
		{
			busAttr.GET("/list", listBusAttr)
			busAttr.POST("/add", addOrUpdateBusAttr)
		}
	}

	// todo: workflow_v3
	v3 := workflow.Group("/v3")
	challenge3 := v3.Group("/challenge", authSvc.Permit(""))
	{
		challenge3.GET("/list", challListV3)
		challenge3.POST("/update", upChallV3)
		challenge3.POST("/reset", rstChallResultV3)
		challenge3.POST("/state/set", setChallStateV3)
		challenge3.POST("/business/state/set", upChallBusStateV3)
		challenge3.POST("/extra/update", upChallExtraV3)
	}
	group3 := v3.Group("/group", authSvc.Permit(""))
	{
		group3.GET("/list", groupListV3)
		group3.GET("/pending/count", countPendingGroup)
		group3.POST("/role/set", setGroupRole)     //角色流转
		group3.POST("/state/set", setGroupStateV3) //工单状态变更/处罚
		group3.POST("/extra/update", upGroupExtra)
		group3.POST("/public/referee/set", setPublicReferee) //移交众裁
	}

	bus3 := v3.Group("/business")
	{
		attr3 := bus3.Group("/attr")
		{
			attr3.GET("/list", listBusAttrV3)
			attr3.POST("/button/switch", setSwitch)
			attr3.POST("/button/shortcut/set", setShortCut)
		}
		mng := bus3.Group("/manager")
		{
			mng.GET("/tag", mngTag)
		}
		source := bus3.Group("/source")
		{
			source.GET("/list", srcList)
		}
	}

	v3.GET("/user/block/info", userBlockInfo) //单个用户封禁次数/封禁状态
}

// ping check server ok.
func ping(ctx *bm.Context) {
	if err := wkfSvc.Ping(ctx); err != nil {
		ctx.Error = err
		ctx.AbortWithStatus(503)
	}
}
