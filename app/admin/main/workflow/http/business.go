package http

import (
	"net/http"
	"net/url"
	"strings"

	"go-common/app/admin/main/workflow/model/param"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/net/http/blademaster/render"
)

func busMetaList(ctx *bm.Context) {
	req := ctx.Request
	itemType := req.Form.Get("item_type")
	ctx.JSON(wkfSvc.ListMeta(ctx, itemType))
}

func listBusAttr(ctx *bm.Context) {
	ctx.JSON(wkfSvc.ListBusAttr(ctx))
}

func listBusAttrV3(ctx *bm.Context) {
	ctx.JSON(wkfSvc.ListBusAttrV3(ctx))
}

func addOrUpdateBusAttr(ctx *bm.Context) {
	abap := &param.AddBusAttrParam{}
	if err := ctx.BindWith(abap, binding.FormPost); err != nil {
		return
	}
	ctx.JSON(nil, wkfSvc.AddOrUpdateBusAttr(ctx, abap))
}

func setSwitch(ctx *bm.Context) {
	bs := new(param.BusAttrButtonSwitch)
	if err := ctx.BindWith(bs, binding.FormPost); err != nil {
		return
	}
	ctx.JSON(nil, wkfSvc.SetSwitch(ctx, bs))
}

func setShortCut(ctx *bm.Context) {
	sc := new(param.BusAttrButtonShortCut)
	if err := ctx.BindWith(sc, binding.FormPost); err != nil {
		return
	}
	if len(sc.ShortCut) != 1 { // only support char
		ctx.Render(http.StatusOK, render.JSON{
			Code:    ecode.RequestErr.Code(),
			Message: "short cut only length 1",
			Data:    nil,
		})
		ctx.Abort()
		return
	}
	sc.ShortCut = strings.ToUpper(sc.ShortCut)
	ctx.JSON(nil, wkfSvc.SetShortCut(ctx, sc))
}

func setExtAPI(ctx *bm.Context) {
	ea := new(param.BusAttrExtAPI)
	if err := ctx.BindWith(ea, binding.FormPost); err != nil {
		return
	}
	if ea.ExternalAPI != "" {
		if _, err := url.Parse(ea.ExternalAPI); err != nil {
			ctx.Render(http.StatusOK, render.JSON{
				Code:    ecode.RequestErr.Code(),
				Message: err.Error(),
				Data:    nil,
			})
			ctx.Abort()
			return
		}
	}
	ctx.JSON(nil, wkfSvc.SetExtAPI(ctx, ea))
}

func mngTag(ctx *bm.Context) {
	ctx.JSON(wkfSvc.ManagerTag(ctx))
}

func userBlockInfo(ctx *bm.Context) {
	bi := new(param.BlockInfo)
	if err := ctx.Bind(bi); err != nil {
		return
	}
	ctx.JSON(wkfSvc.UserBlockInfo(ctx, bi))
}

func srcList(ctx *bm.Context) {
	src := new(param.Source)
	if err := ctx.Bind(src); err != nil {
		return
	}
	ctx.JSON(wkfSvc.SourceList(ctx, src))
}
