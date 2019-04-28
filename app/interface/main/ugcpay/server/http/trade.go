package http

import (
	api "go-common/app/interface/main/ugcpay/api/http"
	"go-common/app/interface/main/ugcpay/model"
	bm "go-common/library/net/http/blademaster"
)

func tradeCreate(ctx *bm.Context) {
	var (
		err      error
		arg      = &api.ArgTradeCreate{}
		resp     = &api.RespTradeCreate{}
		mid, _   = ctx.Get("mid")
		platform = ctx.Request.FormValue("platform")
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if resp.OrderID, resp.PayData, err = srv.TradeCreate(ctx, mid.(int64), platform, arg.OID, arg.OType, arg.Currency); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(resp, err)
}

func tradeQuery(ctx *bm.Context) {
	var (
		err   error
		arg   = &api.ArgTradeOrder{}
		resp  *api.RespTradeOrder
		order *model.TradeOrder
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if order, err = srv.TradeQuery(ctx, arg.OrderID); err != nil {
		ctx.JSON(nil, err)
		return
	}
	resp = &api.RespTradeOrder{
		OrderID:  order.OrderID,
		MID:      order.MID,
		Biz:      order.Biz,
		Platform: order.Platform,
		OID:      order.OID,
		OType:    order.OType,
		Fee:      order.Fee,
		Currency: order.Currency,
		PayID:    order.PayID,
		State:    order.State,
		Reason:   order.Reason,
	}
	ctx.JSON(resp, err)
}

func tradeConfirm(ctx *bm.Context) {
	var (
		err   error
		arg   = &api.ArgTradeOrder{}
		resp  *api.RespTradeOrder
		order *model.TradeOrder
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if order, err = srv.TradeConfirm(ctx, arg.OrderID); err != nil {
		ctx.JSON(nil, err)
		return
	}
	resp = &api.RespTradeOrder{
		OrderID:  order.OrderID,
		MID:      order.MID,
		Biz:      order.Biz,
		Platform: order.Platform,
		OID:      order.OID,
		OType:    order.OType,
		Fee:      order.Fee,
		Currency: order.Currency,
		PayID:    order.PayID,
		State:    order.State,
		Reason:   order.Reason,
	}
	ctx.JSON(resp, err)
}

func tradeCancel(ctx *bm.Context) {
	var (
		err error
		arg = &api.ArgTradeOrder{}
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(nil, srv.TradeCancel(ctx, arg.OrderID))
}
