package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func arcStaff(c *bm.Context) {
	v := new(struct {
		AID   int64 `form:"aid" validate:"required"`
		Cache bool  `form:"cache"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	res, err := pubSvc.StaffList(c, v.AID, v.Cache)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(res, nil)
}
