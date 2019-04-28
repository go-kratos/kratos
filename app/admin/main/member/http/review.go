package http

import (
	"go-common/app/admin/main/member/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func review(ctx *bm.Context) {
	arg := &model.ArgReview{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(svc.Review(ctx, arg))
}

func reviewList(ctx *bm.Context) {
	arg := &model.ArgReviewList{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	arg.IsMonitor = true
	if arg.Pn <= 0 {
		arg.Pn = 1
	}
	if arg.Ps <= 0 {
		arg.Ps = 10
	}
	rws, total, err := svc.Reviews(ctx, arg)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"reviews": rws,
		"page": map[string]int{
			"num":   arg.Pn,
			"size":  arg.Ps,
			"total": total,
		},
	}
	ctx.JSON(res, nil)
}
func reviewFaceList(ctx *bm.Context) {
	arg := &model.ArgReviewList{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	arg.IsMonitor = false
	arg.Property = []int8{model.ReviewPropertyFace}
	if arg.Pn <= 0 {
		arg.Pn = 1
	}
	if arg.Ps <= 0 {
		arg.Ps = 10
	}
	rws, total, err := svc.Reviews(ctx, arg)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"reviews": rws,
		"page": map[string]int{
			"num":   arg.Pn,
			"size":  arg.Ps,
			"total": total,
		},
	}
	ctx.JSON(res, nil)
}

func reviewAudit(ctx *bm.Context) {
	arg := &model.ArgReviewAudit{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	if arg.BlockUser {
		blockArg := model.ArgBatchBlock{}
		if err := ctx.Bind(&blockArg); err != nil {
			return
		}
		// yuzheng: 头像这里的封禁都作为小黑屋封禁
		blockArg.Source = 2
		if !blockArg.Validate() {
			ctx.JSON(nil, ecode.RequestErr)
			return
		}
		arg.ArgBatchBlock = blockArg
	}
	ctx.JSON(nil, svc.ReviewAudit(ctx, arg))
}
