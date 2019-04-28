package http

import (
	"go-common/app/admin/ep/saga/model"
	bm "go-common/library/net/http/blademaster"
)

// @params queryProjectJob
// @router get /ep/admin/saga/v1/data/project/job
// @response TeamDataResp
func queryProjectJob(c *bm.Context) {
	var (
		req = &model.ProjectJobRequest{}
		err error
	)
	if err = c.Bind(req); err != nil {
		return
	}

	if req.Username, err = getUsername(c); err != nil {
		c.JSON(nil, err)
		return
	}
	//c.JSON(srv.QueryProjectJob(c, req))
	c.JSON(srv.QueryProjectJobNew(c, req))
}
