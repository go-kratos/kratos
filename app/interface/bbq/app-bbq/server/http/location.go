package http

import (
	http "go-common/app/interface/bbq/app-bbq/api/http/v1"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

func locationAll(c *bm.Context) {
	arg := new(http.LocationRequest)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}

	c.JSON(srv.GetLocaitonAll(c, arg))
}

func location(c *bm.Context) {
	arg := new(http.LocationRequest)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(srv.GetLocationChild(c, arg))
}
