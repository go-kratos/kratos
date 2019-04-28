package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func kv(c *bm.Context) {
	c.JSON(webSvc.Kv(c))
}

func cmtbox(c *bm.Context) {
	var (
		id  int64
		err error
	)
	params := c.Request.Form
	idStr := params.Get("id")
	if id, err = strconv.ParseInt(idStr, 10, 64); err != nil || id <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(webSvc.CmtBox(c, id))
}

func abServer(c *bm.Context) {
	var (
		mid   int64
		buvid string
	)
	v := new(struct {
		Channel  string `form:"channel"`
		Platform int    `form:"platform"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if ck, err := c.Request.Cookie("buvid3"); err == nil {
		buvid = ck.Value
	}
	if buvid == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(webSvc.AbServer(c, mid, v.Platform, v.Channel, buvid))
}
