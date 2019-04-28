package http

import (
	"go-common/app/service/openplatform/anti-fraud/conf"
	"go-common/app/service/openplatform/anti-fraud/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"net/http"
)

var (
	vfySvc *verify.Verify
	svc    *service.Service
)

// Init init
func Init(c *conf.Config) {
	initService(c)
	// init router
	engine := bm.DefaultServer(nil)
	innerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	vfySvc = verify.New(nil)
	svc = service.New(c)
}

// innerRouter init inner router api path.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)

	//group := e.Group("/openplatform/internal/anti/fraud",idfSvc.Verify)
	group := e.Group("/openplatform/internal/antifraud")
	{

		group.GET("/qusb/info", qusBankInfo)    //题库单条信息
		group.GET("/qusb/list", qusBankList)    //题库列表
		group.POST("/qusb/add", qusBankAdd)     //题库添加
		group.POST("/qusb/del", qusBankDel)     //题库删除
		group.POST("/qusb/update", qusBankEdit) //题库修改
		group.POST("/qusb/check", qusBankCheck) // 答题

		group.GET("/qs/info", qusInfo)           //题目信息
		group.GET("/qs/list", qusList)           //题目列表
		group.POST("/qs/add", qusAdd)            //题目添加
		group.POST("/qs/update", qusUpdate)      //题目更新
		group.POST("/qs/del", qusDel)            //题目删除
		group.GET("/qs/get", getQuestion)        //题目获取
		group.POST("/qs/answer", answerQuestion) // 答题

		group.POST("/bind", questionBankBind)     // 绑定题库
		group.POST("/unbind", questionBankUnbind) // 解绑题库
		group.POST("/bind/bank", getBankBind)     // 查询已绑定的题库
		group.GET("/bind/items", getBindItems)    // 查询绑定到题库的 items

		group.GET("/risk/check", riskCheck)       //风险检验
		group.POST("/risk/check/v2", riskCheckV2) //风险检验v2
		group.GET("/graph/prepare", graphPrepare) //拉起图形验证
		group.POST("/graph/check", graphCheck)    //图形验证
	}

	group2 := e.Group("/openplatform/admin/antifraud/shield")
	{
		group2.GET("/risk/ip/list", ipList)       //ip列表
		group2.GET("/risk/ip/detail", ipDetail)   //ip详情列表
		group2.GET("/risk/uid/list", uidList)     //uid列表
		group2.GET("/risk/uid/detail", uidDetail) //uid详情列表
		group2.POST("/risk/ip/black", ipBlack)    //设置ip黑名单
		group2.POST("/risk/uid/black", uidBlack)  //设置uid黑名单

	}
}

// outerRouter init outer router.
//func outerRouter(e *bm.Engine) {
//}

// ping check server ok.
func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("open-abtest http service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
