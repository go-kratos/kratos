package http

import (
	bm "go-common/library/net/http/blademaster"
)

func listChallActivity(ctx *bm.Context) {
	v := &struct {
		Business int8  `form:"business" validate:"required,gt=0"`
		Cid      int64 `form:"cid" validate:"required,gt=0"`
	}{}
	if err := ctx.Bind(v); err != nil {
		return
	}

	ctx.JSON(wkfSvc.ActivityList(ctx, v.Business, v.Cid))
}
