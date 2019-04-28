package http

import bm "go-common/library/net/http/blademaster"

func view(c *bm.Context) {
	var (
		mid int64
		v   = new(struct {
			Oid  int64 `form:"oid" validate:"required"`
			Type int32 `form:"type" validate:"required"`
			Aid  int64 `form:"aid"`
			Plat int32 `form:"plat"`
		})
	)
	iMid, ok := c.Get("mid")
	if ok {
		mid = iMid.(int64)
	}
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(dmSvc.View(c, mid, v.Aid, v.Oid, v.Type, v.Plat))
}
