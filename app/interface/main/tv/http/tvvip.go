package http

import (
	"fmt"
	"net/http"
	"strings"

	"go-common/app/interface/main/tv/conf"
	tvmdl "go-common/app/interface/main/tv/model/tvvip"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/net/http/blademaster/render"
	"go-common/library/net/metadata"
)

const (
	userAgentWechat = "MicroMessenger"
	userAgentAliPay = "AlipayClient"
	userAgentIphone = "iPhone"
	userAgentIpad   = "iPad"
	agentAndroid    = "Android"

	platformAndroid = 1
	platformIos     = 3
	platformOther   = 4

	errPageURL = "https://www.bilibili.com/blackboard/activity-kWDq8R7f6R.html?code=%d"

	errOrderInvalid        = 93018
	errQrDeviceUnsupported = 93019
	errOrderUnknownErr     = 93020
	errPanelErr            = 93021
	errBuyRateExceededErr  = 93022

	ystErrBadRequest  = "93030"
	ystErrNotFound    = "93031"
	ystErrInternalErr = "93032"
)

func ystErrResp(result string, msg string) map[string]interface{} {
	return map[string]interface{}{
		"result": result,
		"msg":    msg,
	}
}

func ystRender(ctx *bm.Context, data map[string]interface{}) {
	ctx.Render(http.StatusOK, render.MapJSON(data))
}

func errPage(err error) string {
	if ecode.EqualError(ecode.TVIPQrDevideUnsupported, err) {
		return fmt.Sprintf(errPageURL, errQrDeviceUnsupported)
	}
	if ecode.EqualError(ecode.TVIPTokenErr, err) {
		return fmt.Sprintf(errPageURL, errOrderInvalid)
	}
	if ecode.EqualError(ecode.TVIPTokenExpire, err) {
		return fmt.Sprintf(errPageURL, errOrderInvalid)
	}
	if ecode.EqualError(ecode.TVIPDupOrderNo, err) {
		return fmt.Sprintf(errPageURL, errOrderInvalid)
	}

	if ecode.EqualError(ecode.TVIPPanelNotSuitalbe, err) {
		return fmt.Sprintf(errPageURL, errPanelErr)
	}
	if ecode.EqualError(ecode.TVIPPanelNotFound, err) {
		return fmt.Sprintf(errPageURL, errPanelErr)
	}
	if ecode.EqualError(ecode.TVIPBuyNumExceeded, err) {
		return fmt.Sprintf(errPageURL, errPanelErr)
	}
	if ecode.EqualError(ecode.TVIPBuyRateExceeded, err) {
		return fmt.Sprintf(errPageURL, errBuyRateExceededErr)
	}
	if ecode.EqualError(ecode.TVIPMVipRateExceeded, err) {
		return fmt.Sprintf(errPageURL, errBuyRateExceededErr)
	}
	log.Error("errPage(%+v) err(UnExpectedErr)", err)
	return fmt.Sprintf(errPageURL, errOrderUnknownErr)
}

func payTypeFromUA(ctx *bm.Context) (payType string, err error) {
	ua := ctx.Request.UserAgent()
	if strings.Contains(ua, userAgentWechat) {
		return "", ecode.TVIPQrDevideUnsupported
		//return "wechat", nil
	} else if strings.Contains(ua, userAgentAliPay) {
		return "alipay", nil
	} else {
		return "", ecode.TVIPQrDevideUnsupported
	}
}

func platformFromUA(ctx *bm.Context) (platform int8, err error) {
	ua := ctx.Request.UserAgent()
	if strings.Contains(ua, userAgentIpad) {
		return platformIos, nil
	} else if strings.Contains(ua, userAgentIphone) {
		return platformIos, nil
	} else if strings.Contains(ua, agentAndroid) {
		return platformAndroid, nil
	}
	return platformOther, nil
}

func isIpValid(ip string) bool {
	log.Info("ip: %s ipWhiteList: %+v", ip, conf.Conf.IP.White.TvVip)
	if len(conf.Conf.IP.White.TvVip) == 0 {
		return true
	}
	for _, wip := range conf.Conf.IP.White.TvVip {
		if wip == ip {
			return true
		}
	}
	return false
}

// VipInfo implementation
func vipInfo(ctx *bm.Context) {
	mid, ok := ctx.Get("mid")
	if !ok {
		ctx.JSON(nil, ecode.NoLogin)
		return
	}
	ctx.JSON(tvVipSvc.VipInfo(ctx, mid.(int64)))
}

// ystVipInfo implementation
func ystVipInfo(ctx *bm.Context) {
	req := new(tvmdl.YstUserInfoReq)
	if err := ctx.BindWith(req, binding.Query); err != nil {
		ystRender(ctx, ystErrResp(ystErrBadRequest, err.Error()))
		return
	}
	res, err := tvVipSvc.YstVipInfo(ctx, req.Mid, req.Sign)
	if err != nil && ecode.EqualError(ecode.NothingFound, err) {
		ystRender(ctx, ystErrResp(ystErrNotFound, err.Error()))
		return
	}
	if err != nil && ecode.EqualError(ecode.TVIPSignErr, err) {
		ystRender(ctx, ystErrResp(ystErrBadRequest, "SignErr"))
		return
	}
	if err != nil {
		ystRender(ctx, ystErrResp(ystErrInternalErr, err.Error()))
		return
	}
	data := map[string]interface{}{
		"mid":            res.Mid,
		"status":         res.Status,
		"overdue_time":   res.OverdueTime,
		"pay_channel_id": res.PayChannelId,
		"vip_type":       res.VipType,
		"pay_type":       res.PayType,
		"result":         res.Result,
		"msg":            res.Msg,
	}
	ctx.Render(http.StatusOK, render.MapJSON(data))
}

// PanelInfo .
func panelInfo(ctx *bm.Context) {
	mid, ok := ctx.Get("mid")
	if !ok {
		ctx.JSON(nil, ecode.NoLogin)
		return
	}
	res, err := tvVipSvc.PanelInfo(ctx, mid.(int64))
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(res.PriceConfigs, err)
}

func guestPanelInfo(ctx *bm.Context) {
	res, err := tvVipSvc.GuestPanelInfo(ctx)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(res.PriceConfigs, err)
}

func createQr(ctx *bm.Context) {
	req := new(tvmdl.CreateQrReq)
	if err := ctx.Bind(req); err != nil {
		ctx.JSON(nil, err)
		return
	}
	if req.AppChannel == "" {
		log.Warn("createQr(%+v), msg(EmptyAppChannel)", req)
	}
	buvid := ctx.Request.Header.Get("buvid")
	if buvid == "" {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	req.Guid = buvid
	res, err := tvVipSvc.CreateQr(ctx, req)
	ctx.JSON(res, err)
}

func createGuestQr(ctx *bm.Context) {
	req := new(tvmdl.CreateGuestQrReq)
	if err := ctx.Bind(req); err != nil {
		ctx.JSON(nil, err)
		return
	}
	if req.AppChannel == "" {
		log.Warn("createQr(%+v), msg(EmptyAppChannel)", req)
	}
	buvid := ctx.Request.Header.Get("buvid")
	if buvid == "" {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	req.Guid = buvid
	res, err := tvVipSvc.CreateGuestQr(ctx, req)
	ctx.JSON(res, err)
}

func createOrder(ctx *bm.Context) {
	var (
		err error
	)
	req := new(tvmdl.CreateOrderReq)
	if err = ctx.Bind(req); err != nil {
		ctx.JSON(nil, err)
		return
	}
	if req.Platform, err = platformFromUA(ctx); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if req.PaymentType, err = payTypeFromUA(ctx); err != nil {
		ctx.Redirect(302, errPage(err))
		return
	}
	ip := metadata.String(ctx, metadata.RemoteIP)
	if ip == "" {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	res, err := tvVipSvc.CreateOrder(ctx, ip, req)
	if err != nil {
		ctx.Redirect(302, errPage(err))
		return
	}
	ctx.Redirect(302, res.CodeUrl)
}

func createGuestOrder(ctx *bm.Context) {
	var (
		err error
	)
	mid, ok := ctx.Get("mid")
	if !ok {
		ctx.JSON(nil, ecode.NoLogin)
		return
	}
	req := new(tvmdl.CreateGuestOrderReq)
	if err = ctx.Bind(req); err != nil {
		ctx.JSON(nil, err)
		return
	}
	if req.Platform, err = platformFromUA(ctx); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if req.PaymentType, err = payTypeFromUA(ctx); err != nil {
		ctx.Redirect(302, errPage(err))
		return
	}
	ip := metadata.String(ctx, metadata.RemoteIP)
	if ip == "" {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	res, err := tvVipSvc.CreateGuestOrder(ctx, mid.(int64), ip, req)
	if err != nil {
		ctx.Redirect(302, errPage(err))
		return
	}
	ctx.Redirect(302, res.CodeUrl)
}

func tokenStatus(ctx *bm.Context) {
	query := ctx.Request.URL.Query()
	tokens := query["token"]
	res, err := tvVipSvc.TokenInfo(ctx, tokens)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(res.Tokens, err)
}

func payCallback(ctx *bm.Context) {
	ip := metadata.String(ctx, metadata.RemoteIP)
	if !isIpValid(ip) {
		log.Error("payCallback(%s) err(InvalidIP)", ip)
		ystRender(ctx, ystErrResp(ystErrBadRequest, "InvalidIP"))
		return
	}
	req := new(tvmdl.YstPayCallbackReq)
	if err := ctx.BindWith(req, binding.JSON); err != nil {
		ystRender(ctx, ystErrResp(ystErrBadRequest, err.Error()))
		return
	}
	res := tvVipSvc.PayCallback(ctx, req)
	data := map[string]interface{}{
		"traceno": res.TraceNo,
		"result":  res.Result,
		"msg":     res.Msg,
	}
	ystRender(ctx, data)
}

func wxContractCallback(ctx *bm.Context) {
	req := new(tvmdl.WxContractCallbackReq)
	if err := ctx.BindWith(req, binding.JSON); err != nil {
		ystRender(ctx, ystErrResp(ystErrBadRequest, err.Error()))
		return
	}
	res := tvVipSvc.WxContractCallback(ctx, req)
	data := map[string]interface{}{
		"contract_id": res.ContractId,
		"result":      res.Result,
		"msg":         res.Msg,
	}
	ystRender(ctx, data)
}
