package http

import (
	bm "go-common/library/net/http/blademaster"
)

func publish(c *bm.Context) {
	arg := new(struct {
		ResID int `form:"res_id" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(apsSvc.Publish(c, arg.ResID))
}
