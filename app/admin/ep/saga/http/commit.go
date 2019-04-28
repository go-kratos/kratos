package http

import (
	"go-common/app/admin/ep/saga/model"
	bm "go-common/library/net/http/blademaster"
)

// @params ProjectDataReq
// @router get /ep/admin/saga/v1/data/project/report
// @response ProjectDataResp
func queryProjectCommit(c *bm.Context) {
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
	c.JSON(srv.QueryProjectCommit(c, req))
}

// @params TeamDataRequest
// @router get /ep/admin/saga/v1/data/commit/report
// @response TeamDataResp
func queryTeamCommit(c *bm.Context) {
	var (
		req = &model.TeamDataRequest{}
		err error
	)
	if err = c.Bind(req); err != nil {
		return
	}

	if req.Username, err = getUsername(c); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryTeamCommit(c, req))
}

// @params CommitRequest
// @router get /ep/admin/saga/v1/data/commit
// @response CommitResp
func queryCommit(c *bm.Context) {
	var (
		req = &model.CommitRequest{}
		err error
	)
	if err = c.Bind(req); err != nil {
		return
	}

	if req.Username, err = getUsername(c); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryCommit(c, req))
}
