package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func moral(ctx *bm.Context) {
	var (
		err error
		mid int64
		// moral  *model.Moral
		params = ctx.Request.Form
		midStr = params.Get("mid")
		// res    = c.Result()
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		// res["code"] = ecode.RequestErr
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	// if moral, err = memberSvc.Moral(c, mid); err != nil {
	// 	log.Error("memberSvc.Moral(%d) error(%v)", mid, err)
	// 	res["code"] = err
	// 	return
	// }
	// res["data"] = moral
	ctx.JSON(memberSvc.Moral(ctx, mid))
}

func moralLog(ctx *bm.Context) {
	var (
		err    error
		mid    int64
		params = ctx.Request.Form
		midStr = params.Get("mid")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	ctx.JSON(memberSvc.MoralLog(ctx, mid))
}
