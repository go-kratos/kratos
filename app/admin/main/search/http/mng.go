package http

import (
	"go-common/app/admin/main/search/model"
	bm "go-common/library/net/http/blademaster"
)

func businessList(ctx *bm.Context) {
	p := &model.ParamMngBusiness{}
	if err := ctx.Bind(p); err != nil {
		return
	}
	list, total, err := svr.BusinessList(ctx, p.Name, p.Pn, p.Ps)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	data["list"] = list
	data["page"] = &model.Page{
		Pn:    p.Pn,
		Ps:    p.Ps,
		Total: total,
	}
	ctx.JSON(data, nil)
}

func businessAll(ctx *bm.Context) {
	ctx.JSON(svr.BusinessAll(ctx))
}

func businessInfo(ctx *bm.Context) {
	p := &model.ParamMngBusiness{}
	if err := ctx.Bind(p); err != nil {
		return
	}
	ctx.JSON(svr.BusinessInfo(ctx, p.ID))
}

func addBusiness(ctx *bm.Context) {
	p := &model.ParamMngBusiness{}
	if err := ctx.Bind(p); err != nil {
		return
	}
	b := &model.MngBusiness{Name: p.Name, Desc: p.Desc, AppsJSON: "[]"}
	ctx.JSON(svr.AddBusiness(ctx, b))
}

func updateBusiness(ctx *bm.Context) {
	p := &model.ParamMngBusiness{}
	if err := ctx.Bind(p); err != nil {
		return
	}
	b := &model.MngBusiness{ID: p.ID, Name: p.Name, Desc: p.Desc, AppsJSON: p.Apps}
	ctx.JSON(nil, svr.UpdateBusiness(ctx, b))
}

func updateBusinessApp(ctx *bm.Context) {
	p := &model.ParamMngBusinessApp{}
	if err := ctx.Bind(p); err != nil {
		return
	}
	ctx.JSON(nil, svr.UpdateBusinessApp(ctx, p.Business, p.App, p.IncrWay, p.IsJob, p.IncrOpen))
}

func assetList(ctx *bm.Context) {
	p := &model.ParamMngAsset{}
	if err := ctx.Bind(p); err != nil {
		return
	}
	list, total, err := svr.AssetList(ctx, p.Type, p.Name, p.Pn, p.Ps)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	data["list"] = list
	data["page"] = &model.Page{
		Pn:    p.Pn,
		Ps:    p.Ps,
		Total: total,
	}
	ctx.JSON(data, nil)
}

func assetAll(ctx *bm.Context) {
	ctx.JSON(svr.AssetAll(ctx))
}

func assetInfo(ctx *bm.Context) {
	p := &model.ParamMngAsset{}
	if err := ctx.Bind(p); err != nil {
		return
	}
	ctx.JSON(svr.AssetInfo(ctx, p.ID))
}

func addAsset(ctx *bm.Context) {
	p := &model.ParamMngAsset{}
	if err := ctx.Bind(p); err != nil {
		return
	}
	a := &model.MngAsset{Name: p.Name, Type: p.Type, Config: p.Config, Desc: p.Desc}
	ctx.JSON(svr.AddAsset(ctx, a))
}

func updateAsset(ctx *bm.Context) {
	p := &model.ParamMngAsset{}
	if err := ctx.Bind(p); err != nil {
		return
	}
	a := &model.MngAsset{ID: p.ID, Name: p.Name, Type: p.Type, Config: p.Config, Desc: p.Desc}
	ctx.JSON(nil, svr.UpdateAsset(ctx, a))
}

func appList(ctx *bm.Context) {
	p := &model.MngApp{}
	if err := ctx.Bind(p); err != nil {
		return
	}
	ctx.JSON(svr.AppList(ctx, p.Business))
}

func appInfo(ctx *bm.Context) {
	p := &model.MngApp{}
	if err := ctx.Bind(p); err != nil {
		return
	}
	ctx.JSON(svr.AppInfo(ctx, p.ID))
}

func addApp(ctx *bm.Context) {
	p := &model.MngApp{}
	if err := ctx.Bind(p); err != nil {
		return
	}
	ctx.JSON(svr.AddApp(ctx, p))
}

func updateApp(ctx *bm.Context) {
	p := &model.MngApp{}
	if err := ctx.Bind(p); err != nil {
		return
	}
	ctx.JSON(nil, svr.UpdateApp(ctx, p))
}

func countlist(ctx *bm.Context) {
	ctx.JSON(svr.MngCountList(ctx))
}

func count(ctx *bm.Context) {
	p := &model.MngCount{}
	if err := ctx.Bind(p); err != nil {
		return
	}
	ctx.JSON(svr.MngCount(ctx, p))
}

func percent(ctx *bm.Context) {
	p := &model.MngCount{}
	if err := ctx.Bind(p); err != nil {
		return
	}
	ctx.JSON(svr.MngPercent(ctx, p))
}
