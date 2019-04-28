package http

import (
	bm "go-common/library/net/http/blademaster"
)

func userInfo(c *bm.Context) {
	arg := new(struct {
		MID int64 `form:"mid" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(tvSrv.UserInfo(c, arg.MID))

}
