package http

import (
	"strconv"
	"strings"

	"go-common/app/admin/main/aegis/model/common"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func listToken(c *bm.Context) {
	pm := new(net.ListTokenParam)
	if err := c.Bind(pm); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if pm.Sort != "desc" && pm.Sort != "asc" {
		pm.Sort = "desc"
	}

	c.JSON(srv.GetTokenList(c, pm))
}

func tokenGroupByType(c *bm.Context) {
	pm := c.Request.Form.Get("net_id")
	netID, err := strconv.ParseInt(pm, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(srv.TokenGroupByType(c, netID))
}

func tokenByName(c *bm.Context) {
	name := strings.TrimSpace(c.Request.Form.Get("name"))
	busIDStr := c.Request.Form.Get("business_id")
	busID, err := strconv.ParseInt(busIDStr, 10, 64)
	if err != nil || busID <= 0 || name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(srv.TokenByName(c, busID, name))
}

func configToken(c *bm.Context) {
	data := map[string]interface{}{
		"compare":     net.TokenCompareDesc,
		"value_types": net.TokenValueTypeDesc,
	}
	c.JSONMap(data, nil)
}

func showToken(c *bm.Context) {
	pm := c.Request.Form.Get("id")
	id, err := strconv.ParseInt(pm, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(srv.ShowToken(c, id))
}

func preToken(pm *net.Token) (invalid bool) {
	compare := net.GetTokenCompare(pm.Compare)
	tp := net.GetTokenValueType(pm.Type)
	pm.Value = strings.TrimSpace(pm.Value)
	pm.ChName = common.FilterChname(pm.ChName)
	pm.Name = common.FilterName(pm.Name)
	if pm.ChName == "" || pm.Name == "" || compare == "" || tp == "" || pm.Value == "" {
		invalid = true
	}

	return
}

func addToken(c *bm.Context) {
	pm := new(net.Token)
	if err := c.Bind(pm); err != nil || pm.NetID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if preToken(pm) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pm.UID = uid(c)

	id, err, msg := srv.AddToken(c, pm)
	c.JSONMap(map[string]interface{}{
		"id":      id,
		"message": msg,
	}, err)
}
