package http

import (
	"strconv"
	"strings"

	"go-common/app/service/main/member/api"
	"go-common/app/service/main/member/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func realnameStatus(ctx *bm.Context) {
	var (
		err    error
		mid    int64
		status model.RealnameStatus
		params = ctx.Request.Form
		midStr = params.Get("mid")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if status, err = memberSvc.RealnameStatus(ctx, mid); err != nil {
		log.Error("%+v", err)
		ctx.JSON(nil, err)
		return
	}
	var resData struct {
		Status model.RealnameStatus `json:"status"`
	}
	resData.Status = status
	ctx.JSON(resData, nil)
}

func realnameInfo(ctx *bm.Context) {
	var (
		err    error
		mid    int64
		params = ctx.Request.Form
		midStr = params.Get("mid")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	ctx.JSON(memberSvc.RealnameBrief(ctx, mid))
}

func realnameTelCapture(ctx *bm.Context) {
	var (
		err    error
		mid    int64
		params = ctx.Request.Form
		midStr = params.Get("mid")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	_, err = memberSvc.RealnameTelCapture(ctx, mid)
	ctx.JSON(nil, err)
}

func realnameCheckTelCapture(ctx *bm.Context) {
	var (
		err   error
		param = &model.ParamRealnameTelCaptureCheck{}
	)
	if err = ctx.Bind(param); err != nil {
		return
	}
	ctx.JSON(nil, memberSvc.RealnameTelCaptureCheck(ctx, param.MID, param.Capture))
}

func realnameApplyStatus(ctx *bm.Context) {
	var (
		err    error
		mid    int64
		params = ctx.Request.Form
		midStr = params.Get("mid")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	ctx.JSON(memberSvc.RealnameApplyStatus(ctx, mid))
}

func realnameApply(ctx *bm.Context) {
	var (
		err    error
		params = ctx.Request.Form
		// res    = c.Result()

		midStr        = params.Get("mid")
		mid           int64
		realname      = params.Get("real_name")
		cardTypeStr   = params.Get("card_type")
		cardType      int64
		cardNum       = params.Get("card_num")
		countryStr    = params.Get("country")
		country       int64
		captureStr    = params.Get("capture")
		capture       int64
		handIMGToken  = params.Get("img1_token")
		frontIMGToken = params.Get("img2_token")
		backIMGToken  = params.Get("img3_token")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if cardType, err = strconv.ParseInt(cardTypeStr, 10, 64); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if country, err = strconv.ParseInt(countryStr, 10, 64); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if capture, err = strconv.ParseInt(captureStr, 10, 64); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	ctx.JSON(nil, memberSvc.RealnameApply(ctx, mid, int(capture), realname, int8(cardType), cardNum, int16(country), handIMGToken, frontIMGToken, backIMGToken))
}

func realnameAdult(ctx *bm.Context) {
	var (
		err    error
		params = ctx.Request.Form

		midStr = params.Get("mid")
		mid    int64
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	var resp struct {
		Type model.RealnameAdultType `json:"type"`
	}
	resp.Type, err = memberSvc.RealnameAdult(ctx, mid)
	ctx.JSON(resp, err)
}

func realnameCheck(ctx *bm.Context) {
	var (
		err   error
		param = &model.ParamRealnameCheck{}
		resp  bool
	)
	if err = ctx.Bind(param); err != nil {
		return
	}
	resp, err = memberSvc.RealnameCheck(ctx, param.MID, param.CardType, strings.TrimSpace(param.CardCode))
	ctx.JSON(resp, err)
}

func realnameStrippedInfo(ctx *bm.Context) {
	arg := &model.ArgMid2{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(memberSvc.RealnameStrippedInfo(ctx, arg.Mid))
}

func realnameMidByCard(ctx *bm.Context) {
	arg := &api.MidByRealnameCardsReq{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	reply, err := memberSvc.MidByRealnameCard(ctx, arg)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(reply.CodeToMid, nil)
}
