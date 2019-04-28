package http

import (
	"fmt"
	"net/url"

	"go-common/app/admin/main/tv/model"
	bm "go-common/library/net/http/blademaster"
)

func epResult(c *bm.Context) {
	var (
		req   = c.Request.Form
		err   error
		page  int
		order int
	)
	if page, order, err = paramFilter(req); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(tvSrv.EpResult(req, page, order))
}

func seasonResult(c *bm.Context) {
	var (
		req   = c.Request.Form
		err   error
		page  int
		order int
	)
	if page, order, err = paramFilter(req); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(tvSrv.SeasonResult(req, page, order))
}

// filter the params: page & order
func paramFilter(req url.Values) (page int, order int, err error) {
	page = atoi(req.Get("page"))
	order = atoi(req.Get("order"))
	if page == 0 {
		page = 1
	}
	if order == 0 {
		order = 1
	}
	if order != 1 && order != 2 {
		err = fmt.Errorf("Param Order %d is incorrect", order)
		return
	}
	return
}

func arcResult(c *bm.Context) {
	v := new(model.ReqArcCons)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(tvSrv.ArcResult(c, v))
}

func videoResult(c *bm.Context) {
	v := new(model.ReqVideoCons)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(tvSrv.VideoResult(c, v))
}
