package http

import (
	bm "go-common/library/net/http/blademaster"
)

func serverInfos(c *bm.Context) {
	c.JSON(srv.ServerInfos(c))
}

func serverList(c *bm.Context) {
	v := new(struct {
		Platform string `form:"platform" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(srv.ServerList(c, v.Platform), nil)
}

func serverWeight(c *bm.Context) {
	v := new(struct {
		IP string `form:"ip"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	nodes, region, province, err := srv.ServerWeight(c, v.IP)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{})
	data["nodes"] = nodes
	data["region"] = region
	data["province"] = province
	c.JSON(data, err)
}
