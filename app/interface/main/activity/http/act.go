package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func redDot(c *bm.Context) {
	v := new(struct {
		Mid int64 `form:"mid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(likeSvc.RedDot(c, v.Mid))
}

func clearRedDot(c *bm.Context) {
	var loginMid int64
	v := new(struct {
		Mid int64 `form:"mid"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		loginMid = midInter.(int64)
		if loginMid != 0 {
			v.Mid = loginMid
		}
	}
	if v.Mid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, likeSvc.ClearRetDot(c, v.Mid))
}
