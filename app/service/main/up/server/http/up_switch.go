package http

import (
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func switchSet(ctx *bm.Context) {
	param := ctx.Request.Form
	midStr := param.Get("mid")
	fromStr := param.Get("from")
	stateStr := param.Get("state")
	// check params
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if mid <= 0 {
		log.Error("switchSet error mid (%d)", mid)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	from, err := strconv.ParseUint(fromStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fromStr, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	state, err := strconv.Atoi(stateStr)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", stateStr, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	row, err := Svc.SetSwitch(ctx, mid, state, uint8(from))
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(row, nil)
}

func upSwitch(ctx *bm.Context) {
	param := ctx.Request.Form
	midStr := param.Get("mid")
	fromStr := param.Get("from")
	// check params
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if mid <= 0 {
		log.Error("upSwitch error mid(%d)", mid)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	from, err := strconv.Atoi(fromStr)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fromStr, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	state, err := Svc.UpSwitchs(ctx, mid, uint8(from))
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(map[string]interface{}{
		"state": state,
	}, nil)
}
