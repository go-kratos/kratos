package http

import (
	"go-common/app/service/main/member/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func exp(ctx *bm.Context) {
	arg := new(model.ArgMid2)
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(memberSvc.Exp(ctx, arg.Mid))
}

func level(ctx *bm.Context) {
	arg := new(model.ArgMid2)
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(memberSvc.Level(ctx, arg.Mid))
}

func official(ctx *bm.Context) {
	arg := new(model.ArgMid2)
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(memberSvc.Official(ctx, arg.Mid))
}

func explog(ctx *bm.Context) {
	arg := new(model.ArgMid2)
	if err := ctx.Bind(arg); err != nil {
		return
	}
	arg.RealIP = metadata.String(ctx, metadata.RemoteIP)
	ctx.JSON(memberSvc.ExpLog(ctx, arg.Mid, arg.RealIP))
}

func updateExp(ctx *bm.Context) {
	arg := new(model.ArgAddExp)
	if err := ctx.Bind(arg); err != nil {
		return
	}
	arg.IP = metadata.String(ctx, metadata.RemoteIP)
	ctx.JSON(nil, memberSvc.UpdateExp(ctx, arg))
}

func setExp(ctx *bm.Context) {
	arg := new(model.ArgAddExp)
	if err := ctx.Bind(arg); err != nil {
		return
	}
	arg.IP = metadata.String(ctx, metadata.RemoteIP)
	ctx.JSON(nil, memberSvc.SetExp(ctx, arg))
}

func stat(ctx *bm.Context) {
	arg := new(model.ArgMid2)
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(memberSvc.Stat(ctx, arg.Mid))
}
