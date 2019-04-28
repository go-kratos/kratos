package http

import (
	"time"

	"go-common/app/interface/main/space/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func setTopDynamic(c *bm.Context) {
	v := new(struct {
		DyID int64 `form:"dy_id" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(nil, spcSvc.SetTopDynamic(c, mid, v.DyID))
}

func cancelTopDynamic(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(nil, spcSvc.CancelTopDynamic(c, mid, time.Now()))
}

func dynamicList(c *bm.Context) {
	v := new(model.DyListArg)
	if err := c.Bind(v); err != nil {
		return
	}
	if v.Pn > 1 && v.LastTime == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		v.Mid = midInter.(int64)
	}
	c.JSON(spcSvc.DynamicList(c, v))
}
