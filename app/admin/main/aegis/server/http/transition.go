package http

import (
	"fmt"
	"strconv"
	"strings"

	"go-common/app/admin/main/aegis/model/common"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func listTransition(c *bm.Context) {
	pm := new(net.ListNetElementParam)
	if err := c.Bind(pm); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if pm.Sort != "desc" && pm.Sort != "asc" {
		pm.Sort = "desc"
	}

	c.JSON(srv.GetTransitionList(c, pm))
}

func getTranByNet(c *bm.Context) {
	pm := c.Request.Form.Get("net_id")
	id, err := strconv.ParseInt(pm, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(srv.GetTranByNet(c, id))
}

func showTransition(c *bm.Context) {
	pm := c.Request.Form.Get("id")
	id, err := strconv.ParseInt(pm, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(srv.ShowTransition(c, id))
}

func preTran(pm *net.TransitionEditParam) (invalid bool) {
	invalid = true
	pm.Trigger = net.TriggerManual
	pm.ChName = common.FilterChname(pm.ChName)
	pm.Name = common.FilterName(pm.Name)
	if pm.ChName == "" || pm.Name == "" || (pm.Trigger == net.TriggerManual && pm.Limit > 50) {
		return
	}

	existBindMap := map[string]int{}
	for _, item := range pm.TokenList {
		item.TokenID = strings.TrimSpace(item.TokenID)
		item.ChName = strings.TrimSpace(item.ChName)
		if item.ID < 0 || item.TokenID == "" || item.Type == net.BindTypeFlow {
			return
		}

		kk := fmt.Sprintf("%s_%d", item.TokenID, item.Type)
		if existBindMap[kk] > 0 {
			log.Error("preTran bind(%+v) duplicated in tokenid+type", item)
			return
		}
		existBindMap[kk] = 1
	}
	invalid = false
	return
}

func addTransition(c *bm.Context) {
	pm := new(net.TransitionEditParam)
	if err := c.BindWith(pm, binding.JSON); err != nil || pm.NetID <= 0 {
		log.Error("addTransition bind params error(%v) body(%+v)", err, c.Request.Body)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if preTran(pm) {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	admin := uid(c)

	id, err, msg := srv.AddTransition(c, admin, pm)
	c.JSONMap(map[string]interface{}{
		"id":      id,
		"message": msg,
	}, err)
}

func updateTransition(c *bm.Context) {
	pm := new(net.TransitionEditParam)
	if err := c.BindWith(pm, binding.JSON); err != nil || pm.ID <= 0 {
		log.Error("updateTransition bind params error(%v) body(%+v)", err, c.Request.Body)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if preTran(pm) {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	admin := uid(c)
	err, msg := srv.UpdateTransition(c, admin, pm)
	c.JSONMap(map[string]interface{}{
		"message": msg,
	}, err)
}

func switchTransition(c *bm.Context) {
	pm := new(net.SwitchParam)
	if err := c.Bind(pm); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(nil, srv.SwitchTransition(c, pm.ID, pm.Disable))
}
