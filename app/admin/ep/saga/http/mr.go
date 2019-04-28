package http

import (
	"go-common/app/admin/ep/saga/model"
	bm "go-common/library/net/http/blademaster"
)

// @params ProjectDataReq
// @router get /ep/admin/saga/v1/data/project/report
// @response ProjectDataResp
func queryProjectMr(c *bm.Context) {
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
	c.JSON(srv.QueryProjectMr(c, req))
}

// @params TeamDataRequest
// @router get /ep/admin/saga/v1/data/mr/report
// @response TeamDataResp
func queryTeamMr(c *bm.Context) {
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
	c.JSON(srv.QueryTeamMr(c, req))
}

//@params ProjectMrReportReq
//@router get /ep/admin/saga/v1/data/project/mr/report
//@response ProjectMrReportResp
func queryProjectMrReport(c *bm.Context) {
	var (
		req = &model.ProjectMrReportReq{}
		err error
	)
	if err = c.Bind(req); err != nil {
		return
	}
	if req.Username, err = getUsername(c); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryProjectMrReport(c, req))
}
