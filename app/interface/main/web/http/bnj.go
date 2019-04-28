package http

import bm "go-common/library/net/http/blademaster"

func bnj2019(c *bm.Context) {
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(webSvc.Bnj2019(c, mid))
}

func bnj2019Aids(c *bm.Context) {
	data := make(map[string]interface{}, 1)
	data["list"] = webSvc.Bnj2019Aids(c)
	c.JSON(data, nil)
}

func timeline(c *bm.Context) {
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	list, err := webSvc.Timeline(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 1)
	data["list"] = list
	c.JSON(data, nil)
}
