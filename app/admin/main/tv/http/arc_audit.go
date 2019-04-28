package http

import (
	bm "go-common/library/net/http/blademaster"
)

//add archive
func arcAdd(c *bm.Context) {
	v := new(struct {
		Aids []int64 `form:"aids,split" validate:"required,min=1,dive,gt=0"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(tvSrv.AddArcs(v.Aids))
}
