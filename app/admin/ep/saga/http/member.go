package http

import (
	"go-common/app/admin/ep/saga/model"
	bm "go-common/library/net/http/blademaster"
)

func queryProjectMembers(c *bm.Context) {
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
	c.JSON(srv.QueryProjectMembers(c, req))
}
