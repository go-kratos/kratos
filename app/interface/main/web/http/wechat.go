package http

import (
	"go-common/app/interface/main/web/model"
	bm "go-common/library/net/http/blademaster"
)

func wxHot(c *bm.Context) {
	v := new(struct {
		Pn int `form:"pn" default:"1" validate:"min=1"`
		Ps int `form:"ps" default:"100" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	list, count, err := webSvc.WxHot(c, v.Pn, v.Ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"pn":    v.Pn,
		"ps":    v.Ps,
		"count": count,
	}
	data["data"] = list
	data["page"] = page
	c.JSONMap(data, nil)
}

func wxSearchAll(c *bm.Context) {
	var (
		mid   int64
		buvid string
		err   error
	)
	v := new(model.SearchAllArg)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Pn <= 0 {
		v.Pn = 1
	}
	if ck, err := c.Request.Cookie("buvid3"); err == nil {
		buvid = ck.Value
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(webSvc.SearchAll(c, mid, v, buvid, c.Request.Header.Get("User-Agent"), model.WxSearchType))
}
