package http

import (
	"go-common/app/service/main/point/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func pointInfo(c *bm.Context) {
	mid, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(pointSvc.PointInfo(c, mid.(int64)))
}

func pointPage(c *bm.Context) {
	var err error
	mid, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	arg := new(model.ArgRPCPointHistory)
	if err = c.Bind(arg); err != nil {
		return
	}
	arg.Mid = mid.(int64)
	c.JSON(pointSvc.PointPage(c, arg))
}
