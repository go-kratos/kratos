package http

import bm "go-common/library/net/http/blademaster"

func channel(c *bm.Context) {
	var (
		mid   int64
		buvid string
	)
	v := new(struct {
		ID int64 `form:"id" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if ck, err := c.Request.Cookie("buvid3"); err == nil {
		buvid = ck.Value
	}
	c.JSON(srvWeb.Channel(c, v.ID, mid, buvid))
}
