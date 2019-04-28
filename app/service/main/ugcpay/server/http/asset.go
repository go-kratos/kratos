package http

import (
	api "go-common/app/service/main/ugcpay/api/http"
	"go-common/app/service/main/ugcpay/model"
	bm "go-common/library/net/http/blademaster"
)

func assetQuery(ctx *bm.Context) {
	var (
		err   error
		arg   = &api.ArgAssetQuery{}
		resp  = &api.RespAssetQuery{}
		asset *model.Asset
		pp    map[string]int64
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if asset, pp, err = srv.AssetQuery(ctx, arg.OID, arg.OType, arg.Currency); err != nil {
		ctx.JSON(nil, err)
		return
	}
	resp.Parse(asset, pp)
	ctx.JSON(resp, err)
}

func assetRegister(ctx *bm.Context) {
	var (
		err error
		arg = &api.ArgAssetRegister{}
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(nil, srv.AssetRegister(ctx, arg.MID, arg.OID, arg.OType, arg.Currency, arg.Price))
}

func assetRelation(ctx *bm.Context) {
	var (
		err   error
		arg   = &api.ArgAssetRelation{}
		resp  = &api.RespAssetRelation{}
		state string
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if state, err = srv.AssetRelation(ctx, arg.MID, arg.OID, arg.OType); err != nil {
		ctx.JSON(nil, err)
		return
	}
	resp.State = state
	ctx.JSON(resp, err)
}

func assetRelationDetail(ctx *bm.Context) {
	var (
		err   error
		arg   = &api.ArgAssetRelationDetail{}
		resp  = &api.RespAssetRelationDetail{}
		asset *model.Asset
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if resp.RelationState, err = srv.AssetRelation(ctx, arg.MID, arg.OID, arg.OType); err != nil {
		ctx.JSON(nil, err)
		return
	}
	if asset, resp.AssetPlatformPrice, err = srv.AssetQuery(ctx, arg.OID, arg.OType, arg.Currency); err != nil {
		ctx.JSON(nil, err)
		return
	}
	resp.AssetPrice = asset.Price
	ctx.JSON(resp, err)
}
