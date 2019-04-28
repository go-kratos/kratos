package http

import (
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func logByRpID(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			RpID int64 `form:"rpid" validate:"required"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	list, err := rpSvc.LogsByRpID(c, v.RpID)
	if err != nil {
		log.Error("svc.ActionInfo(%+v) error(%v)", v, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(list, nil)
}
