package http

import (
	"go-common/app/admin/main/workflow/model/param"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func listCallback(ctx *bm.Context) {
	ctx.JSON(wkfSvc.ListCallback(ctx))
}

func addOrUpCallback(ctx *bm.Context) {
	cbp := &param.AddCallbackParam{}
	if err := ctx.BindWith(cbp, binding.JSON); err != nil {
		return
	}

	if cbp.State > 0 {
		cbp.State = 1
	}

	cbID, err := wkfSvc.AddOrUpCallback(ctx, cbp)
	if err != nil {
		log.Error("wkfSvc.AddUpCallback(%+v) error(%v)", cbp, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	ctx.JSON(map[string]int32{
		"callbackNo": cbID,
	}, nil)
}
