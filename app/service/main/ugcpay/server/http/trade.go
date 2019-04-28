package http

import (
	"net/http"

	api "go-common/app/service/main/ugcpay/api/http"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func tradePayRefund(ctx *bm.Context) {
	var (
		err error
		arg = &api.ArgTradeRefund{}
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(nil, srv.TradeRefund(ctx, arg.OrderID))
}

func tradePayRefunds(ctx *bm.Context) {
	var (
		err error
		arg = &api.ArgTradeRefunds{}
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if len(arg.OrderIDs) > 20 {
		err = ecode.RequestErr
		return
	}
	ctx.JSON(nil, srv.TradeRefunds(ctx, arg.OrderIDs))
}

func tradePayCallback(ctx *bm.Context) {
	var (
		err    error
		arg    = &api.ArgTradeCallback{}
		retMSG string
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if retMSG, err = srv.TradePayCallback(ctx, arg.MSGID, arg.MSGContent); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.String(http.StatusOK, retMSG)
}

func tradePayRefundCallback(ctx *bm.Context) {
	var (
		err    error
		arg    = &api.ArgTradeCallback{}
		retMSG string
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if retMSG, err = srv.TradeRefundCallback(ctx, arg.MSGID, arg.MSGContent); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.String(http.StatusOK, retMSG)
}

func tradePayRechargeCallback(ctx *bm.Context) {
	var (
		err    error
		arg    = &api.ArgTradeCallback{}
		retMSG string
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if retMSG, err = srv.TradeRefundCallback(ctx, arg.MSGID, arg.MSGContent); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.String(http.StatusOK, retMSG)
}
