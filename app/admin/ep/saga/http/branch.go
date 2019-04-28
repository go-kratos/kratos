package http

import (
	"go-common/app/admin/ep/saga/model"
	bm "go-common/library/net/http/blademaster"
)

func queryProjectBranchList(c *bm.Context) {
	var (
		req = &model.ProjectDataReq{}
		err error
	)
	if err = c.Bind(req); err != nil {
		return
	}

	if req.Username, err = getUsername(c); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryProjectBranchList(c, req))
}

// @params queryBranchDiffWith
// @router get /ep/admin/saga/v1/data/branch/report
// @response BranchDiffWithRequest
func queryBranchDiffWith(c *bm.Context) {
	var (
		req = &model.BranchDiffWithRequest{}
		err error
	)
	if err = c.Bind(req); err != nil {
		return
	}

	if req.Username, err = getUsername(c); err != nil {
		return
	}
	c.JSON(srv.QueryBranchDiffWith(c, req))
}
