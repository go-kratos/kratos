package http

import (
	bm "go-common/library/net/http/blademaster"
)

func albumIndex(c *bm.Context) {
	v := new(struct {
		Mid int64 `form:"mid" validate:"min=1"`
		Ps  int   `form:"ps" default:"8" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(spcSvc.AlbumIndex(c, v.Mid, v.Ps))
}
