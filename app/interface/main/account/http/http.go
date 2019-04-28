package http

import (
	"net/http"

	"go-common/app/interface/main/account/conf"
	"go-common/app/interface/main/account/service/card"
	"go-common/app/interface/main/account/service/coupon"
	"go-common/app/interface/main/account/service/geetest"
	"go-common/app/interface/main/account/service/member"
	"go-common/app/interface/main/account/service/passport"
	"go-common/app/interface/main/account/service/point"
	"go-common/app/interface/main/account/service/realname"
	rls "go-common/app/interface/main/account/service/relation"
	us "go-common/app/interface/main/account/service/usersuit"
	"go-common/app/interface/main/account/service/vip"
	vipverify "go-common/app/service/main/vip/verify"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/supervisor"
	v "go-common/library/net/http/blademaster/middleware/verify"
)

var (
	relationSvc *rls.Service
	memberSvc   *member.Service
	passSvc     *passport.Service
	vipSvc      *vip.Service
	realnameSvc *realname.Service
	usSvc       *us.Service
	couponSvc   *coupon.Service
	pointSvc    *point.Service
	cardSvc     *card.Service
	geetestSvr  *geetest.Service
	// api middleware
	authn        *auth.Auth
	verify       *v.Verify
	anti         *antispam.Antispam
	batchRelAnti *antispam.Antispam
	smsAnti      *antispam.Antispam
	faceAnti     *antispam.Antispam
	vipAnti      *antispam.Antispam
	spv          *supervisor.Supervisor
	// vip third verify
	vipThirdVerify *vipverify.Verify
)

// Init init http sever instance.
func Init(c *conf.Config) {
	// service
	initService(c)
	vipThirdVerify = vipverify.NewThirdVerify(c.VipThirdVerifyConfig)
	// init outer router
	// innerEngine := bm.DefaultServer(c.BM.Inner)
	innerEngine := bm.NewServer(c.BM.Inner)
	innerEngine.Use(bm.Recovery(), bm.Trace(), bm.Logger(), bm.Mobile())
	setupInnerEngine(innerEngine)
	if err := innerEngine.Start(); err != nil {
		log.Error("innerEngine.Start() error(%v)", err)
		panic(err)
	}
}

func setupInnerEngine(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)

	// member
	mr := e.Group("/x/member", bm.CSRF())
	mr.GET("/notice", authn.User, notice)
	mr.POST("/notice/close", authn.User, closeNotice)

	sec := mr.Group("/security", authn.User)
	sec.GET("/status", status)
	sec.POST("/feedback", feedback)
	sec.POST("/close", closeNotify)

	app := mr.Group("/app", authn.UserMobile)
	app.GET("/nickfree", nickFree)
	app.POST("/uname/update", spv.ServeHTTP, updateUname)
	app.POST("/sign/update", spv.ServeHTTP, updateSign)
	app.POST("/sex/update", updateSex)
	app.POST("/birthday/update", updateBirthday)
	app.POST("/face/update", spv.ServeHTTP, faceAnti.ServeHTTP, updateFace)
	app.POST("/pendant/equip", pendantEquip)
	app.GET("/point/flag", pointFlagMobile)

	web := mr.Group("/web", authn.User)
	web.POST("/update", spv.ServeHTTP, update)
	web.POST("/face/update", spv.ServeHTTP, faceAnti.ServeHTTP, updateFace)
	web.POST("/sign/update", spv.ServeHTTP, updateSign)
	web.POST("/uname/update", spv.ServeHTTP, updateUname)
	web.GET("/account", account)
	web.GET("/reply/list", replyHistoryList)
	web.GET("/coin/log", logCoin)
	web.GET("/moral/log", logMoral)
	web.GET("/exp/log", logExp)
	web.GET("/exp/reward", reward)
	web.GET("/login/log", logLogin)
	web.GET("/coin", coin)
	web.POST("/birthday/update", updateBirthday)
	web.POST("/pendant/equip", pendantEquip)
	web.GET("/point/flag", pointFlag)

	sudo := mr.Group("/sudo", sudo)
	sudo.POST("/notify-purge-cache", notityPurgeCache)

	// captcha 验证码
	cap := mr.Group("/captcha", authn.UserWeb)
	cap.GET("/geetest", getChallenge)           //获取极验图形验证
	cap.POST("/geetest/check", geetestValidate) //校验极验验证码

	// vip
	vip := e.Group("/x/vip", bm.CSRF())
	vip.GET("/code/verify", authn.User, codeVerify)
	vip.POST("/code/open", authn.User, codeOpen)
	vip.GET("/tips", authn.User, tips)
	vip.GET("/price/panel", authn.User, vipAnti.ServeHTTP, vipPanel)
	vip.GET("/coupon/list", authn.User, vipAnti.ServeHTTP, couponList)
	vip.GET("/coupon/usable", authn.User, vipAnti.ServeHTTP, couponUsable)
	vip.POST("/coupon/unlock", authn.User, vipAnti.ServeHTTP, couponUnlock)
	vip.GET("/privilege/bysid", vipAnti.ServeHTTP, privilegeBySid)
	vip.GET("/privilege/bytype", vipAnti.ServeHTTP, privilegeByType)
	vip.GET("/manager", authn.User, vipAnti.ServeHTTP, vipManagerInfo)
	vip.POST("/unfreeze", authn.UserWeb, vipAnti.ServeHTTP, unfrozen)
	vip.GET("/frozenTime", authn.UserWeb, vipAnti.ServeHTTP, frozenTime)
	vip.GET("/code/openeds", vipAnti.ServeHTTP, codeOpeneds)
	vip.GET("/public/panel", authn.Guest, vipAnti.ServeHTTP, publicPriceList)
	vip.POST("/batch/use", vipAnti.ServeHTTP, useBatch)
	vip.GET("/order/status", authn.UserWeb, orderStatus)
	vip.GET("/price/panel/v8", authn.Guest, vipAnti.ServeHTTP, vipPanelV8)
	vip.GET("/prize/cards", authn.UserWeb, vipAnti.ServeHTTP, prizeCards)
	vip.POST("/prize/draw", authn.UserWeb, vipAnti.ServeHTTP, prizeDraw)
	vip.GET("/resource/banner", authn.UserMobile, vipAnti.ServeHTTP, resourceBanner) // 大会员落地页
	vip.GET("/resource/buy", authn.UserMobile, vipAnti.ServeHTTP, resourceBuy)       // 大会员购买页

	vip.GET("/coupon/usable/v2", authn.Guest, vipAnti.ServeHTTP, couponBySuitIDV2)
	vip.GET("/price/panel/v9", authn.Guest, vipAnti.ServeHTTP, vipPanelV9)

	// associate
	vip.GET("/associate/info", authn.User, bindInfoByMid)
	vip.GET("/associate/panel", authn.User, actlimit, associatePanel)
	vip.POST("/associate/create/order", authn.User, actlimit, createAssociateOrder)
	vip.GET("/associate/ele/oauth", authn.User, actlimit, eleOAuthURL)
	vip.GET("/associate/ele/redpackets", redpackets)
	vip.GET("/associate/ele/specailfoods", specailfoods)

	//vip welfare
	vip.GET("/welfare/list", authn.Guest, welfareList)
	vip.GET("/welfare/type", authn.Guest, welfareTypeList)
	vip.GET("/welfare/info", authn.Guest, welfareInfo)
	vip.POST("/welfare", authn.User, vipAnti.ServeHTTP, receiveWelfare)
	vip.GET("/welfare/my", authn.User, myWelfare)

	// vip third verify.

	// ele oauth callback.
	e.GET("/x/oauth2/v1/callback", authn.User, actlimit, openAuthCallBack)

	// bilibili third.
	oauth2 := e.Group("/x/oauth2/v1", iplimit, vipThirdVerify.Verify, bm.CSRF())
	oauth2.POST("/access_token", openIDByOAuth2Code)

	// bilibili third vip.
	vipThird := oauth2.Group("/vip", iplimit, openlimit)
	vipThird.POST("/bind", openBindByOutOpenID)
	vipThird.GET("/account", userInfoByOpenID)
	vipThird.POST("/delivery", bilibiliVipGrant)
	vipThird.POST("/prize/grant", bilibiliPrizeGrant)

	// vip app api
	vip.GET("/v1/order/status", authn.UserMobile, orderStatus)
	vip.GET("/v1/frozenTime", authn.UserMobile, vipAnti.ServeHTTP, frozenTime)
	vip.POST("/v1/unfrozen", authn.UserMobile, vipAnti.ServeHTTP, unfrozen)
	vip.GET("/v2/tips", authn.User, tipsv2)
	vip.GET("/v2/price/panel", authn.Guest, vipAnti.ServeHTTP, vipPanelV2)

	invite := mr.Group("/invite", authn.UserWeb)
	invite.GET("/stat", inviteStat)
	invite.POST("/buy", buy)
	invite.POST("/apply", apply)

	// medal 勋章
	medal := mr.Group("/medal")
	// medal.GET("/user/info", authn.User, medalHomeInfo)
	medal.GET("/user/info", authn.User, medalUserInfo)
	medal.POST("/install", authn.User, medalInstall)
	medal.GET("/popup", authn.User, medalPopup)
	medal.GET("/my/info", authn.User, medalMyInfo)
	medal.GET("/all/info", authn.User, medalAllInfo)

	official := mr.Group("/official", bm.CSRF(), authn.User)
	official.POST("/submit", submitOffical)
	official.GET("/doc", officialDoc)
	official.GET("/conditions", officialConditions)
	official.POST("/upload/image", uploadImage)
	official.POST("/mobile/verify", smsAnti.ServeHTTP, mobileVerify)
	official.GET("/permission", officialPermission)
	official.GET("/auto/fill/doc", officialAutoFillDoc)
	official.GET("/monthly/times", monthlyOfficialSubmittedTimes)

	identifyG := mr.Group("/identify", authn.UserWeb)
	identifyG.GET("/info", identifyInfo)

	// realname
	realnameG := mr.Group("/realname", authn.User)
	realnameG.GET("/channel", realnameChannel)
	realnameG.GET("/status", realnameStatus)
	realnameG.GET("/apply/status", realnameApplyStatus)
	realnameG.POST("/apply", realnameApply)
	realnameG.POST("/upload", realnameUpload)
	realnameG.GET("/preview", realnamePreview)
	realnameG.GET("/countrylist", realnameCountryList)
	realnameG.GET("/card/types", realnameCardTypes)
	realnameG.GET("/v2/card/types", realnameCardTypesV2)
	realnameG.POST("/tel/capture", realnameTelCapture)
	realnameG.GET("/tel/info", realnameTelInfo)
	realnameG.GET("/captcha", realnameCaptcha)
	realnameG.GET("/captcha/refresh", realnameCaptchaRefresh)
	realnameG.GET("/captcha/confirm", realnameCaptchaConfirm)
	realnameG.POST("/alipay/apply", realnameAlipayApply)
	realnameG.GET("/alipay/confirm", realnameAlipayConfirm)

	// passport
	passportR := mr.Group("/passport", authn.User)
	passportR.GET("/testUserName", testUserName)

	// member v2
	memberV2 := e.Group("x/member/v2", bm.CSRF())
	memberV2.GET("/notice", authn.UserMobile, noticeV2)
	memberV2.POST("/notice/close", authn.UserMobile, closeNoticeV2)

	// relation
	relationG := e.Group("/x/relation", bm.CSRF())
	relationG.GET("", authn.User, relation)
	relationG.GET("/relations", authn.User, relations)
	relationG.GET("/blacks", authn.User, blacks)
	relationG.GET("/whispers", authn.User, whispers)
	relationG.GET("/friends", authn.User, friends)
	relationG.POST("/modify", authn.User, anti.ServeHTTP, modify)
	relationG.POST("/batch/modify", authn.User, batchRelAnti.ServeHTTP, batchModify)
	relationG.GET("/followings", authn.Guest, followings)
	relationG.GET("/same/followings", authn.User, sameFollowings)
	relationG.GET("/followers", authn.Guest, followers)
	relationG.GET("/stat", authn.Guest, stat)
	relationG.GET("/stats", authn.Guest, stats)
	relationG.GET("/recommend", authn.User, recommend)
	relationG.GET("/recommend/followlist_empty", authn.User, bm.Mobile(), recommendFollowlistEmpty)
	relationG.GET("/recommend/answer_ok", authn.User, bm.Mobile(), recommendAnswerOK)
	relationG.GET("/recommend/tag_suggest", authn.User, bm.Mobile(), recommendTagSuggest)
	relationG.GET("/recommend/tag_suggest/detail", authn.User, bm.Mobile(), recommendTagSuggestDetail)
	// relation tag
	relationG.GET("/tag", authn.User, tag)
	relationG.GET("/tags", authn.User, tags)
	relationG.POST("/tag/special/add", authn.User, addSpecial)
	relationG.POST("/tag/special/del", authn.User, delSpecial)
	relationG.GET("/tag/special", authn.User, special)
	relationG.GET("/tag/user", authn.User, tagUser)
	relationG.POST("/tag/create", authn.User, tagCreate)
	relationG.POST("/tag/update", authn.User, tagUpdate)
	relationG.POST("/tag/del", authn.User, tagDel)
	relationG.POST("/tags/addUsers", authn.User, tagsAddUsers)
	relationG.POST("/tags/copyUsers", authn.User, tagsCopyUsers)
	relationG.POST("/tags/moveUsers", authn.User, tagsMoveUsers)
	// for mobile.
	relationG.GET("/tag/m/tags", authn.User, mobileTags)
	// 提示用户关注该up主
	relationG.POST("/prompt", authn.User, prompt)
	relationG.POST("/prompt/close", authn.User, closePrompt)
	// 粉丝提醒功能
	relationG.GET("/followers/unread", authn.User, unread)
	relationG.POST("/followers/unread/reset", authn.User, unreadReset)
	relationG.GET("/followers/unread/count", authn.User, unreadCount)
	relationG.POST("/followers/unread/count/reset", authn.User, unreadCountReset)
	relationG.GET("/followers/notify", authn.User, followerNotifySetting)
	relationG.POST("/followers/notify/enable", authn.User, enableFollowerNotify)
	relationG.POST("/followers/notify/disable", authn.User, disableFollowerNotify)
	// achieve
	relationG.POST("/achieve/award/get", authn.User, achieveGet)
	relationG.GET("/achieve/award", authn.Guest, achieve)

	// pendant group
	pendant := e.Group("/x/pendant", bm.CSRF())
	// pendant with web
	pendant.GET("/current", authn.UserWeb, pendantCurrent)
	pendant.GET("/all", authn.UserWeb, pendantAll)
	pendant.GET("/my", authn.UserWeb, pendantMy)
	pendant.GET("/myHistory", authn.UserWeb, pendantMyHistory)
	pendant.GET("/bigEntry", authn.UserWeb, pendantEntry)
	pendant.GET("/vipRecommend", authn.UserWeb, pendantVIP)
	pendant.POST("/checkOrder", authn.UserWeb, pendantCheckOrder)
	pendant.POST("/vipGet", authn.UserWeb, pendantVIPGet)
	pendant.POST("/order", authn.UserWeb, pendantOrder)
	// pendant with app
	pendant.GET("/pointEntry", authn.UserMobile, pendantEntry)
	// pendent with vri
	pendant.GET("/single", verify.Verify, pendantSingle)

	coupon := e.Group("/x/coupon", bm.CSRF(), authn.User)
	coupon.GET("/allowance/list", allowanceList)
	coupon.GET("/list", couponPage)
	coupon.GET("/code/verify", captchaToken)
	coupon.POST("/code/exchange", useCouponCode)

	point := e.Group("/x/point", bm.CSRF(), authn.User)
	point.GET("/info", pointInfo)
	point.GET("/history", pointPage)

	card := e.Group("/x/card", bm.CSRF())
	card.GET("/bymid", authn.User, userCard)
	card.GET("/info", cardInfo)
	card.GET("/hots", cardHots)
	card.GET("/groups", authn.Guest, cardGroups)
	card.GET("/bytype", cardsByGid)
	card.POST("/equip", authn.User, equip)
	card.POST("/demount", authn.User, demount)
}

// ping check server ok.
func ping(c *bm.Context) {
	var err error
	if err = memberSvc.Ping(c); err != nil {
		c.JSON(nil, err)
		log.Error("service ping error(%v)", err)
		c.Writer.WriteHeader(http.StatusServiceUnavailable)
	}
}

func register(ctx *bm.Context) {
	ctx.JSON(nil, nil)
}

func initService(c *conf.Config) {
	relationSvc = rls.New(c)
	memberSvc = member.New(c)
	realnameSvc = realname.New(c, conf.RsaPub(), conf.RsaPriv(), conf.AlipayPub(), conf.AlipayBiliPriv())
	passSvc = passport.New(c)
	cardSvc = card.New(c)
	vipSvc = vip.New(c)
	usSvc = us.New(c)
	geetestSvr = geetest.New(c)
	// api middleware
	authn = auth.New(c.AuthN)
	verify = v.New(c.Verify)
	anti = antispam.New(c.Antispam)
	batchRelAnti = antispam.New(c.BatchRelAntispam)
	smsAnti = antispam.New(c.SMSAntispam)
	faceAnti = antispam.New(c.FaceAntispam)
	vipAnti = antispam.New(c.VIPAntispam)
	spv = supervisor.New(c.Supervisor)
	couponSvc = coupon.New(c)
	pointSvc = point.New(c)
}
