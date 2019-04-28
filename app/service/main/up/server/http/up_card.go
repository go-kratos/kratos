package http

import (
	"strconv"
	"strings"

	"go-common/app/service/main/up/model"
	"go-common/app/service/main/up/service"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
)

func getCardInfo(ctx *blademaster.Context) {

	var r = new(model.GetCardByMidArg)
	if err := ctx.Bind(r); err != nil {
		log.Error("request argument bind fail, err=%v", err)
		err = ecode.RequestErr
		return
	}

	// check params
	mid := r.Mid
	if mid <= 0 {
		log.Error("getCardInfo error mid (%d)", mid)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	card, err := Svc.GetCardInfo(ctx, mid)
	if err != nil {
		service.BmHTTPErrorWithMsg(ctx, ecode.ServerErr, err.Error())
		return
	}
	ctx.JSON(card, err)
}

func listCardBase(ctx *blademaster.Context) {
	mids, err := Svc.ListCardBase(ctx)

	if err != nil {
		service.BmHTTPErrorWithMsg(ctx, ecode.ServerErr, err.Error())
		return
	}

	ctx.JSON(map[string]interface{}{
		"mids": mids,
	}, nil)
}

func listCardDetail(ctx *blademaster.Context) {
	arg := new(model.ListUpCardInfoArg)
	if err := ctx.Bind(arg); err != nil {
		log.Error("request argument bind fail, err=%v", err)
		err = ecode.RequestErr
		return
	}

	offset := (arg.Pn - 1) * arg.Ps
	cards, total, err := Svc.ListCardDetail(ctx, offset, arg.Ps)
	if err != nil {
		service.BmHTTPErrorWithMsg(ctx, ecode.ServerErr, err.Error())
		return
	}

	var data = &model.UpCardInfoPage{
		Cards: cards,
		Page:  &model.Pager{Num: arg.Pn, Size: uint(len(cards)), Total: total},
	}

	ctx.JSON(data, nil)
}

func listCardByMids(ctx *blademaster.Context) {
	arg := new(model.ListCardByMidsArg)
	if err := ctx.Bind(arg); err != nil {
		log.Error("request argument bind fail, err=%v", err)
		err = ecode.RequestErr
		return
	}

	var mids []int64
	midStrs := strings.Split(arg.Mids, ",")
	for _, midStr := range midStrs {
		mid, e := strconv.ParseInt(strings.Trim(midStr, " \n"), 10, 64)
		if e != nil {
			continue
		}
		mids = append(mids, mid)
	}

	cards, err := Svc.GetCardInfoByMids(ctx, mids)
	if err != nil {
		service.BmHTTPErrorWithMsg(ctx, ecode.ServerErr, err.Error())
		return
	}
	ctx.JSON(cards, nil)
}
