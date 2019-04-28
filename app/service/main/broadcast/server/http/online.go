package http

import bm "go-common/library/net/http/blademaster"

func onlineTop(c *bm.Context) {
	v := new(struct {
		Business string `form:"business" validate:"required"`
		Num      int    `form:"num" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(srv.OnlineTop(c, v.Business, v.Num))
}

func onlineRoom(c *bm.Context) {
	v := new(struct {
		Business string   `form:"business" validate:"required"`
		Rooms    []string `form:"rooms" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(srv.OnlineRoom(c, v.Business, v.Rooms))
}

func onlineTotal(c *bm.Context) {
	res := make(map[string]int64)
	res["ip_count"], res["conn_count"] = srv.OnlineTotal(c)
	c.JSON(res, nil)
}
