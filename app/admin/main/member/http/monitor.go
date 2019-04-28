package http

import (
	"go-common/app/admin/main/member/model"
	bm "go-common/library/net/http/blademaster"
)

func monitors(ctx *bm.Context) {
	arg := &model.ArgMonitor{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	if arg.Pn <= 0 {
		arg.Pn = 1
	}
	if arg.Ps <= 0 {
		arg.Ps = 10
	}
	mns, total, err := svc.Monitors(ctx, arg)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"monitors": mns,
		"page": map[string]int{
			"num":   arg.Pn,
			"size":  arg.Ps,
			"total": total,
		},
	}
	ctx.JSON(res, nil)
}

func addMonitor(ctx *bm.Context) {
	arg := &model.ArgAddMonitor{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	arg.OperatorID = operatorID(ctx)
	ctx.JSON(nil, svc.AddMonitor(ctx, arg))
}

func delMonitor(ctx *bm.Context) {
	arg := &model.ArgDelMonitor{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	arg.OperatorID = operatorID(ctx)
	ctx.JSON(nil, svc.DelMonitor(ctx, arg))
}

func operatorID(ctx *bm.Context) int64 {
	uidI, ok := ctx.Get("uid")
	if !ok {
		return 0
	}
	uid, _ := uidI.(int64)
	return uid
}
