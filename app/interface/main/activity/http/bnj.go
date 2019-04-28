package http

import bm "go-common/library/net/http/blademaster"

func previewInfo(c *bm.Context) {
	var loginMid int64
	if midInter, ok := c.Get("mid"); ok {
		loginMid = midInter.(int64)
	}
	c.JSON(bnjSvc.PreviewInfo(c, loginMid), nil)
}

func timeline(c *bm.Context) {
	var loginMid int64
	if midInter, ok := c.Get("mid"); ok {
		loginMid = midInter.(int64)
	}
	c.JSON(bnjSvc.Timeline(c, loginMid), nil)
}

func reset(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	cd, err := bnjSvc.TimeReset(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]int64{"cd": cd}, nil)
}

func reward(c *bm.Context) {
	v := new(struct {
		Step int `form:"step" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(nil, likeSvc.Reward(c, mid, v.Step))
}

func delTime(c *bm.Context) {
	v := new(struct {
		Key string `form:"key" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, bnjSvc.DelTime(c, v.Key))
}

func fail(c *bm.Context) {
	c.JSON(nil, nil)
}
