package http

import (
	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/param"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func addReply(ctx *bm.Context) {
	var (
		eid int64
		err error
	)

	e := &param.EventParam{}
	if err = ctx.BindWith(e, binding.FormPost); err != nil {
		return
	}
	e.AdminID, e.AdminName = adminInfo(ctx)

	if eid, err = wkfSvc.AddEvent(ctx, e); err != nil {
		ctx.JSON(nil, err)
		return
	}

	// 管理员回复同步修改 business_state
	if e.Event == 1 {
		if err = wkfSvc.UpChallBusState(ctx, e.Cid, e.AdminID, e.AdminName, model.FeedbackReplyNotRead); err != nil {
			ctx.JSON(nil, err)
			return
		}
	}

	ctx.JSON(map[string]int64{
		"eventNo": eid,
	}, nil)
}

func batchAddReply(ctx *bm.Context) {
	var (
		eids []int64
		err  error
	)
	bep := &param.BatchEventParam{}
	if err = ctx.BindWith(bep, binding.FormPost); err != nil {
		return
	}
	bep.AdminID, bep.AdminName = adminInfo(ctx)

	eids, err = wkfSvc.BatchAddEvent(ctx, bep)
	if err != nil {
		log.Error("wkfSvc.BatchAddEvent(%v) error(%v)", bep, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	// 管理员回复同步修改 business_state
	if bep.Event == 1 {
		if err = wkfSvc.BatchUpChallBusState(ctx, bep.Cids, bep.AdminID, bep.AdminName, model.FeedbackReplyNotRead); err != nil {
			ctx.JSON(nil, err)
			return
		}
	}

	ctx.JSON(map[string][]int64{
		"eventNo": eids,
	}, nil)
}
