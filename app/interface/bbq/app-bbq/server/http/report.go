package http

import (
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

func reportConfig(c *bm.Context) {
	arg := new(v1.Base)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	res := &v1.ReportConfigResponse{
		Report:  model.Reports,
		Reasons: model.Reasons,
	}
	c.JSON(res, nil)
}

func reportReport(c *bm.Context) {
	arg := new(v1.ReportRequest)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	accessKey := c.Request.Form.Get("access_key")
	c.JSON(srv.Report(c, arg, mid.(int64), accessKey))
}
