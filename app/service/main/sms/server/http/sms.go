package http

import (
	pb "go-common/app/service/main/sms/api"
	bm "go-common/library/net/http/blademaster"
)

func send(ctx *bm.Context) {
	req := new(pb.SendReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(smsSvc.Send(ctx, req))
}

func sendBatch(ctx *bm.Context) {
	req := new(pb.SendBatchReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(smsSvc.SendBatch(ctx, req))
}
