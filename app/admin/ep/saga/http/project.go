package http

import (
	"go-common/app/admin/ep/saga/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

// @params EmptyReq
// @router get /ep/admin/saga/v1/projects/favorite
// @response FavoriteProjectsResp
func favoriteProjects(ctx *bm.Context) {
	var (
		req      = &model.Pagination{}
		err      error
		userName string
	)
	if err = ctx.Bind(req); err != nil {
		ctx.JSON(nil, err)
		return
	}
	if userName, err = getUsername(ctx); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(srv.FavoriteProjects(ctx, req, userName))
}

// @params EditFavoriteReq
// @router post /ep/admin/saga/v1/projects/favorite/edit
// @response EmptyResp
func editFavorite(ctx *bm.Context) {
	var (
		req      = &model.EditFavoriteReq{}
		err      error
		userName string
	)
	if err = ctx.BindWith(req, binding.JSON); err != nil {
		ctx.JSON(nil, err)
		return
	}
	if userName, err = getUsername(ctx); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(srv.EditFavorite(ctx, req, userName))
}

func queryCommonProjects(ctx *bm.Context) {
	ctx.JSON(srv.QueryCommonProjects(ctx))
}

// @params QueryProjectInfoRequest
// @router get /ep/admin/saga/v1/data/project
// @response ProjectsResp
func queryProjectInfo(c *bm.Context) {
	var (
		req = &model.ProjectInfoRequest{}
		err error
	)
	if err = c.Bind(req); err != nil {
		return
	}
	if err = req.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryProjectInfo(c, req))
}
