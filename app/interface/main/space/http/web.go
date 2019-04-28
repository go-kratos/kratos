package http

import bm "go-common/library/net/http/blademaster"

func webIndex(c *bm.Context) {
	var mid int64
	v := new(struct {
		VMid int64 `form:"mid" validate:"min=1"`
		Pn   int32 `form:"pn" default:"1" validate:"min=1"`
		Ps   int32 `form:"ps" default:"20" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(spcSvc.WebIndex(c, mid, v.VMid, v.Pn, v.Ps))
}
