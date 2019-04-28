package http

import (
	"go-common/app/interface/main/report-click/model"
	bm "go-common/library/net/http/blademaster"
)

func errReport(c *bm.Context) {
	v := new(model.ErrReport)
	if err := c.Bind(v); err != nil {
		return
	}
	clickSvr.ErrReport(c, v)
	c.JSON(nil, nil)
}
