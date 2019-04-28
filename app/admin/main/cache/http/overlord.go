package http

import (
	"go-common/app/admin/main/cache/model"
	bm "go-common/library/net/http/blademaster"
)

// @params OverlordReq
// @router get /x/admin/cache/overlord/clusters
// @response OverlordResp
func overlordClusters(ctx *bm.Context) {
	req := new(model.OverlordReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(srv.OverlordClusters(ctx, req))
}

// @params OverlordReq
// @router get /x/admin/cache/overlord/ops/names
// @response EmpResp
func overlordOpsClusterNames(ctx *bm.Context) {
	req := new(model.OverlordReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(srv.OpsClusterNames(ctx, req))
}

// @params OverlordReq
// @router get /x/admin/cache/overlord/ops/nodes
// @response EmpResp
func overlordOpsNodes(ctx *bm.Context) {
	req := new(model.OverlordReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(srv.OpsClusterNodes(ctx, req))
}

// @params OverlordReq
// @router get /x/admin/cache/overlord/import/ops/cluster
// @response EmpResp
func overlordImportCluster(ctx *bm.Context) {
	req := new(model.OverlordReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(srv.ImportOpsCluster(ctx, req))
}

// @params OverlordReq
// @router get /x/admin/cache/overlord/new/ops/node
// @response EmpResp
func overlordClusterNewNode(ctx *bm.Context) {
	req := new(model.OverlordReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(srv.ImportOpsNode(ctx, req))
}

// @params OverlordReq
// @router get /x/admin/cache/overlord/replace/ops/node
// @response EmpResp
func overlordClusterReplaceNode(ctx *bm.Context) {
	req := new(model.OverlordReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(srv.ReplaceOpsNode(ctx, req))
}

// @params OverlordReq
// @router get /x/admin/cache/overlord/cluster/del
// @response EmpResp
func overlordDelCluster(ctx *bm.Context) {
	req := new(model.OverlordReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(srv.DelOverlordCluster(ctx, req))
}

// @params OverlordReq
// @router get /x/admin/cache/overlord/node/del
// @response EmpResp
func overlordDelNode(ctx *bm.Context) {
	req := new(model.OverlordReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(srv.DelOverlordNode(ctx, req))
}

// @params OverlordReq
// @router get /x/admin/cache/overlord/app/clusters
// @response OverlordResp
func overlordAppClusters(ctx *bm.Context) {
	req := new(model.OverlordReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	req.Cookie = ctx.Request.Header.Get("Cookie")
	ctx.JSON(srv.OverlordAppClusters(ctx, req))
}

// @params OverlordReq
// @router get /x/admin/cache/overlord/app/can/bind/clusters
// @response OverlordResp
func overlordAppNeedClusters(ctx *bm.Context) {
	req := new(model.OverlordReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(srv.OverlordAppCanBindClusters(ctx, req))
}

// @params OverlordReq
// @router get /x/admin/cache/overlord/app/cluster/bind
// @response OverlordResp
func overlordAppClusterBind(ctx *bm.Context) {
	req := new(model.OverlordReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	req.Cookie = ctx.Request.Header.Get("Cookie")
	ctx.JSON(srv.OverlordAppClusterBind(ctx, req))
}

// @params OverlordReq
// @router get /x/admin/cache/overlord/app/cluster/del
// @response OverlordResp
func overlordAppClusterDel(ctx *bm.Context) {
	req := new(model.OverlordReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	req.Cookie = ctx.Request.Header.Get("Cookie")
	ctx.JSON(srv.OverlordAppClusterDel(ctx, req))
}

// @params OverlordReq
// @router get /x/admin/cache/overlord/app/appids
// @response OverlordResp
func overlordAppAppIDs(ctx *bm.Context) {
	req := new(model.OverlordReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	req.Cookie = ctx.Request.Header.Get("Cookie")
	ctx.JSON(srv.OverlordAppAppIDs(ctx, req))
}

func overlordToml(ctx *bm.Context) {
	req := new(model.OverlordReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	resp, err := srv.OverlordToml(ctx, req)
	if err != nil {
		ctx.Status(500)
		return
	}
	ctx.Writer.Write(resp)
}
