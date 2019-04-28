package http

import (
	"strconv"

	"go-common/app/admin/main/aegis/model/common"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func listFlow(c *bm.Context) {
	pm := new(net.ListNetElementParam)
	if err := c.Bind(pm); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pm.Sort != "desc" && pm.Sort != "asc" {
		pm.Sort = "desc"
	}

	c.JSON(srv.GetFlowList(c, pm))
}

func getFlowByNet(c *bm.Context) {
	pm := c.Request.Form.Get("net_id")
	id, err := strconv.ParseInt(pm, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(srv.GetFlowByNet(c, id))
}

func showFlow(c *bm.Context) {
	pm := c.Request.Form.Get("id")
	id, err := strconv.ParseInt(pm, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(srv.ShowFlow(c, id))
}

func preFlow(pm *net.FlowEditParam) (invalid bool) {
	pm.ChName = common.FilterChname(pm.ChName)
	pm.Name = common.FilterName(pm.Name)
	invalid = pm.ChName == "" || pm.Name == ""
	return
}

func addFlow(c *bm.Context) {
	pm := new(net.FlowEditParam)
	if err := c.Bind(pm); err != nil || pm.NetID <= 0 {
		log.Error("addFlow bind params error(%v) pm(%+v)", err, pm)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if preFlow(pm) {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	admin := uid(c)
	id, err, msg := srv.AddFlow(c, admin, pm)
	c.JSONMap(map[string]interface{}{
		"id":      id,
		"message": msg,
	}, err)
}

func updateFlow(c *bm.Context) {
	pm := new(net.FlowEditParam)
	if err := c.Bind(pm); err != nil || pm.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if preFlow(pm) {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	admin := uid(c)
	err, msg := srv.UpdateFlow(c, admin, pm)
	c.JSONMap(map[string]interface{}{"message": msg}, err)
}

func switchFlow(c *bm.Context) {
	pm := new(net.SwitchParam)
	if err := c.Bind(pm); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(nil, srv.SwitchFlow(c, pm.ID, pm.Disable))
}
