package http

import (
	"strconv"

	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func listDirection(c *bm.Context) {
	pm := new(net.ListDirectionParam)
	if err := c.Bind(pm); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pm.Sort != "desc" && pm.Sort != "asc" {
		pm.Sort = "desc"
	}
	//direction只有在确立了flow_id/transition_id后才能选择
	if pm.FlowID == 0 && pm.TransitionID == 0 {
		pm.Direction = 0
	}

	c.JSON(srv.GetDirectionList(c, pm))
}

func showDirection(c *bm.Context) {
	pm := c.Request.Form.Get("id")
	id, err := strconv.ParseInt(pm, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(srv.ShowDirection(c, id))
}

func preDir(pm *net.DirEditParam) (invalid bool) {
	if pm.Direction == net.DirInput && pm.Output != "" {
		pm.Output = ""
	}
	invalid = pm.Order == net.DirOrderOrSplit && pm.Guard == ""
	return
}

func addDirection(c *bm.Context) {
	pm := new(net.DirEditParam)
	if err := c.Bind(pm); err != nil || pm.NetID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if preDir(pm) {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	admin := uid(c)
	id, err, msg := srv.AddDirection(c, admin, pm)
	c.JSONMap(map[string]interface{}{
		"id":      id,
		"message": msg,
	}, err)
}

func updateDirection(c *bm.Context) {
	pm := new(net.DirEditParam)
	if err := c.Bind(pm); err != nil || pm.ID <= 0 {
		log.Error("updateDirection bind error(%+v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if preDir(pm) {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	admin := uid(c)

	err, msg := srv.UpdateDirection(c, admin, pm)
	c.JSONMap(map[string]interface{}{
		"message": msg,
	}, err)
}

func switchDirection(c *bm.Context) {
	pm := new(net.SwitchParam)
	if err := c.Bind(pm); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	err, msg := srv.SwitchDirection(c, pm.ID, pm.Disable)
	c.JSONMap(map[string]interface{}{
		"message": msg,
	}, err)
}
