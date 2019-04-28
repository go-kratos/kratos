package http

import (
	"net/http"

	"go-common/app/interface/main/app-wall/conf"
	"go-common/app/interface/main/app-wall/service/mobile"
	"go-common/app/interface/main/app-wall/service/offer"
	"go-common/app/interface/main/app-wall/service/operator"
	pingSvr "go-common/app/interface/main/app-wall/service/ping"
	"go-common/app/interface/main/app-wall/service/telecom"
	"go-common/app/interface/main/app-wall/service/unicom"
	"go-common/app/interface/main/app-wall/service/wall"
	"go-common/library/ecode"
	log "go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/proxy"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/http/blademaster/render"
)

var (
	// depend service
	verifySvc *verify.Verify
	authSvc   *auth.Auth
	// self service
	wallSvc     *wall.Service
	offerSvc    *offer.Service
	unicomSvc   *unicom.Service
	mobileSvc   *mobile.Service
	pingSvc     *pingSvr.Service
	telecomSvc  *telecom.Service
	operatorSvc *operator.Service
)

func Init(c *conf.Config) {
	initService(c)
	// init external router
	engineOut := bm.DefaultServer(c.BM.Outer)
	outerRouter(engineOut)
	// init Outer server
	if err := engineOut.Start(); err != nil {
		log.Error("engineOut.Start() error(%v)", err)
		panic(err)
	}
}

// initService init services.
func initService(c *conf.Config) {
	verifySvc = verify.New(nil)
	authSvc = auth.New(&auth.Config{DisableCSRF: true})
	// init self service
	wallSvc = wall.New(c)
	offerSvc = offer.New(c)
	unicomSvc = unicom.New(c)
	mobileSvc = mobile.New(c)
	pingSvc = pingSvr.New(c)
	telecomSvc = telecom.New(c)
	operatorSvc = operator.New(c)
}

func outerRouter(e *bm.Engine) {
	e.Ping(ping)
	// formal api
	proxyHandler := proxy.NewZoneProxy("sh004", "http://sh001-app.bilibili.com")
	w := e.Group("/x/wall")
	{
		w.GET("/get", walls)
		op := w.Group("/operator", authSvc.Guest)
		{
			op.GET("/ip", userOperatorIP)
			op.GET("/m/ip", mOperatorIP)
			op.GET("/reddot", reddot)
		}
		of := w.Group("/offer")
		{
			of.GET("/exist", wallExist)
			of.POST("/click/shike", proxyHandler, wallShikeClick)
			of.GET("/click/dotinapp", wallDotinappClick)
			of.GET("/click/gdt", wallGdtClick)
			of.GET("/click/toutiao", wallToutiaoClick)
			of.POST("/active", proxyHandler, verifySvc.Verify, wallActive)
			of.GET("/active/test", wallTestActive)
			of.POST("/active2", proxyHandler, verifySvc.Verify, wallActive2)
		}
		uc := w.Group("/unicom", proxyHandler)
		{
			// unicomSync
			uc.POST("/orders", ordersSync)
			uc.POST("/advance", advanceSync)
			uc.POST("/flow", flowSync)
			uc.POST("/ip", inIPSync)
			// unicom
			uc.GET("/userflow", verifySvc.Verify, userFlow)
			uc.GET("/user/userflow", bm.CORS(), userFlowState)
			uc.GET("/userstate", verifySvc.Verify, userState)
			uc.GET("/state", verifySvc.Verify, unicomState)
			uc.GET("/m/state", unicomStateM)
			uc.POST("/pack", authSvc.User, bm.CORS(), pack)
			uc.GET("/userip", bm.CORS(), isUnciomIP)
			uc.GET("/user/ip", bm.CORS(), userUnciomIP)
			uc.POST("/order/pay", bm.CORS(), orderPay)
			uc.POST("/order/cancel", bm.CORS(), orderCancel)
			uc.POST("/order/smscode", authSvc.Guest, bm.CORS(), smsCode)
			uc.POST("/order/bind", authSvc.Guest, bm.CORS(), addUnicomBind)
			uc.POST("/order/untie", authSvc.Guest, bm.CORS(), releaseUnicomBind)
			uc.GET("/bind/info", authSvc.Guest, bm.CORS(), userBind)
			uc.GET("/pack/list", authSvc.Guest, bm.CORS(), packList)
			uc.POST("/order/pack/receive", authSvc.Guest, bm.CORS(), packReceive)
			uc.POST("/order/pack/flow", authSvc.Guest, bm.CORS(), flowPack)
			uc.GET("/order/userlog", authSvc.User, bm.CORS(), userBindLog)
			uc.GET("/pack/log", userPacksLog)
			uc.GET("/bind/state", verifySvc.Verify, welfareBindState)
		}
		mb := w.Group("/mobile", proxyHandler)
		{
			mb.POST("/orders.so", ordersMobileSync)
			mb.GET("/activation", verifySvc.Verify, mobileActivation)
			mb.GET("/status", verifySvc.Verify, mobileState)
			mb.GET("/user/status", bm.CORS(), userMobileState)
		}
		tl := w.Group("/telecom", proxyHandler)
		{
			tl.POST("/orders.so", bm.CORS(), telecomOrdersSync)
			tl.POST("/flow.so", bm.CORS(), telecomMsgSync)
			tl.POST("/order/pay", bm.CORS(), telecomPay)
			tl.POST("/order/pay/cancel", bm.CORS(), cancelRepeatOrder)
			tl.GET("/order/consent", verifySvc.Verify, orderConsent)
			tl.GET("/order/list", verifySvc.Verify, orderList)
			tl.GET("/order/user/flow", bm.CORS(), phoneFlow)
			tl.POST("/send/sms", verifySvc.Verify, phoneSendSMS)
			tl.GET("/verification", verifySvc.Verify, phoneVerification)
			tl.GET("/order/state", bm.CORS(), orderState)
		}
	}
}

//returnDataJSON return json no message
func returnDataJSON(c *bm.Context, data map[string]interface{}, err error) {
	code := http.StatusOK
	if data == nil {
		c.JSON(data, err)
		return
	}
	if _, ok := data["message"]; !ok {
		data["message"] = ""
	}
	if err != nil {
		c.Error = err
		bcode := ecode.Cause(err)
		data["code"] = bcode.Code()
	} else {
		if _, ok := data["code"]; !ok {
			data["code"] = ecode.OK
		}
		data["ttl"] = 1
	}
	c.Render(code, render.MapJSON(data))
}
