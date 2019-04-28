package http

import (
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func addWhite(c *bm.Context) {
	v := new(struct {
		MID  int64 `form:"mid" validate:"required"`
		Type int   `form:"type" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.InsertWhite(c, v.MID, v.Type)
	if err != nil {
		log.Error("growup svr.AddWhite error(%v)", err)
	}
	c.JSON(nil, err)
}
