package http

import (
	"strconv"

	"go-common/app/admin/main/aegis/model/common"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func listNet(c *bm.Context) {
	pm := new(net.ListNetParam)
	if err := c.Bind(pm); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if pm.Sort != "desc" && pm.Sort != "asc" {
		pm.Sort = "desc"
	}

	c.JSON(srv.GetNetList(c, pm))
}

func getNetByBusiness(c *bm.Context) {
	pm := c.Request.Form.Get("business_id")
	nid, err := strconv.ParseInt(pm, 10, 64)
	if err != nil || nid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(srv.GetNetByBusiness(c, nid))
}

func showNet(c *bm.Context) {
	pm := c.Request.Form.Get("id")
	nid, err := strconv.ParseInt(pm, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(srv.ShowNet(c, nid))
}

func addNet(c *bm.Context) {
	pm := new(net.Net)
	if err := c.Bind(pm); err != nil || pm.BusinessID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pm.ChName = common.FilterChname(pm.ChName)
	if pm.ChName == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	pm.UID = uid(c)
	id, err, msg := srv.AddNet(c, pm)
	c.JSONMap(map[string]interface{}{
		"id":      id,
		"message": msg,
	}, err)
}

func updateNet(c *bm.Context) {
	pm := new(net.NetEditParam)
	if err := c.Bind(pm); err != nil || pm.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pm.ChName = common.FilterChname(pm.ChName)
	if pm.ChName == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	admin := uid(c)
	err, msg := srv.UpdateNet(c, admin, pm)
	c.JSONMap(map[string]interface{}{
		"message": msg,
	}, err)
}

func switchNet(c *bm.Context) {
	pm := new(net.SwitchParam)
	if err := c.Bind(pm); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(nil, srv.SwitchNet(c, pm.ID, pm.Disable))
}
