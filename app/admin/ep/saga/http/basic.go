package http

import (
	"go-common/app/admin/ep/saga/model"
	bm "go-common/library/net/http/blademaster"
)

func queryProjectStatus(c *bm.Context) {
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
	c.JSON(srv.QueryProjectStatus(c, req), nil)
}

func queryProjectTypes(c *bm.Context) {
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
	c.JSON(srv.QueryProjectTypes(c, req), nil)
}
