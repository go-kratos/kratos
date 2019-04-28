package http

import (
	"go-common/app/admin/ep/saga/model"
	bm "go-common/library/net/http/blademaster"
)

// @params TasksReq
// @router get /ep/admin/saga/v1/tasks/project
// @response TasksResp
func projectTasks(ctx *bm.Context) {
	var (
		req = &model.TasksReq{}
		err error
	)
	if err = ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(srv.MergeTasks(ctx, req))
}
