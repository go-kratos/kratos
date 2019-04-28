package http

import (
	"strconv"

	"go-common/app/admin/main/workflow/model/param"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func addEvent(ctx *bm.Context) {
	ep := &param.EventParam{}
	if err := ctx.BindWith(ep, binding.JSON); err != nil {
		return
	}

	if !ep.ValidComponent() {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	eid, err := wkfSvc.AddEvent(ctx, ep)
	if err != nil {
		log.Error("wkfSvc.AddEvent(%v) error(%v)", ep, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	ctx.JSON(map[string]int64{
		"eventNo": eid,
	}, nil)
}

func batchAddEvent(ctx *bm.Context) {
	bep := &param.BatchEventParam{}
	if err := ctx.BindWith(bep, binding.JSON); err != nil {
		return
	}

	if !bep.ValidComponent() {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	eids, err := wkfSvc.BatchAddEvent(ctx, bep)
	if err != nil {
		log.Error("wkfSvc.BatchAddEvent(%v) error(%v)", bep, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	ctx.JSON(map[string][]int64{
		"eventNo": eids,
	}, nil)
}

func eventList(ctx *bm.Context) {
	params := ctx.Request.Form
	cidStr := params.Get("cid")
	// check params
	cid, err := strconv.ParseInt(cidStr, 10, 32)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", cidStr, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	ctx.JSON(wkfSvc.ListEvent(ctx, cid))
}
