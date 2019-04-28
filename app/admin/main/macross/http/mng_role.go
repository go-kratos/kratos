package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// role get role by sysid.
func role(c *bm.Context) {
	var (
		params = c.Request.Form
		system string
	)
	if system = params.Get("system"); system == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.Role(c, system), nil)
}

// saveRole save role.
func saveRole(c *bm.Context) {
	var (
		params   = c.Request.Form
		roleID   int64
		roleName string
	)
	if roleName = params.Get("role_name"); roleName == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	roleIDStr := params.Get("role_id")
	roleID, _ = strconv.ParseInt(roleIDStr, 10, 64)
	system := params.Get("system")
	if roleID == 0 && system == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svr.SaveRole(c, roleID, system, roleName))
}

// DelRole del role.
func DelRole(c *bm.Context) {
	var (
		params = c.Request.Form
		roleID int64
		err    error
	)
	roleIDStr := params.Get("role_id")
	if roleID, err = strconv.ParseInt(roleIDStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svr.DelRole(c, roleID))
}
