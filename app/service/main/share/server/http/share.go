package http

import (
	"strconv"

	"go-common/app/service/main/share/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

func add(ctx *bm.Context) {
	p := &model.ShareParams{}
	if err := ctx.Bind(p); err != nil {
		return
	}
	p.IP = metadata.String(ctx, metadata.RemoteIP)
	ctx.JSON(svr.Add(ctx, p))
}

func stat(ctx *bm.Context) {
	var (
		oid int64
		tp  int64
		err error
	)
	params := ctx.Request.Form
	oidStr := params.Get("oid")
	if oid, err = strconv.ParseInt(oidStr, 10, 64); err != nil || oid <= 0 {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	tpStr := params.Get("tp")
	if tp, err = strconv.ParseInt(tpStr, 10, 64); err != nil || tp <= 0 {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	ctx.JSON(svr.Stat(ctx, oid, int(tp)))
}

func stats(ctx *bm.Context) {
	var (
		oids []int64
		tp   int64
		err  error
	)
	params := ctx.Request.Form
	oidsStr := params.Get("oids")
	if oids, err = xstr.SplitInts(oidsStr); err != nil || len(oids) == 0 || len(oids) > 30 {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	tpStr := params.Get("tp")
	if tp, err = strconv.ParseInt(tpStr, 10, 64); err != nil || tp <= 0 {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	ctx.JSON(svr.Stats(ctx, oids, int(tp)))
}
