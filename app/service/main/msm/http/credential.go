package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const (
	_contentType = "Content-Type"
	_urlJSON     = "application/json"
)

// credential verify.
func credential(ctx *bm.Context) {
	var (
		appIDStr  string
		signature string
	)
	req := ctx.Request
	params := req.Form
	header := req.Header
	if header.Get(_contentType) == _urlJSON {
		appIDStr = header.Get("App-Tree-ID")
		signature = header.Get("Signature")
		header.Del("Signature")
	} else {
		appIDStr = params.Get("app_tree_id")
		signature = params.Get("signature")
		params.Del("signature")
	}
	appID, _ := strconv.ParseInt(appIDStr, 10, 64)
	if appID == 0 || signature == "" {
		ctx.JSON(nil, ecode.RequestErr)
		ctx.Abort()
		return
	}
	if ok := svr.CheckSign(appID, signature); !ok {
		ctx.JSON(nil, ecode.SignCheckErr)
		ctx.Abort()
	}
}
