package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// user get user by sysid.
func user(c *bm.Context) {
	var (
		params = c.Request.Form
		system string
	)
	if system = params.Get("system"); system == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.User(c, system), nil)
}

// saveUser save user.
func saveUser(c *bm.Context) {
	var (
		params         = c.Request.Form
		roleID, userID int64
		name           string
		err            error
	)
	if name = params.Get("user_name"); name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	roleIDStr := params.Get("role_id")
	if roleID, err = strconv.ParseInt(roleIDStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	userIDStr := params.Get("user_id")
	userID, _ = strconv.ParseInt(userIDStr, 10, 64)
	system := params.Get("system")
	if userID == 0 && system == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svr.SaveUser(c, roleID, userID, system, name))
}

// delUser del user.
func delUser(c *bm.Context) {
	var (
		params = c.Request.Form
		userID int64
		err    error
	)
	userIDStr := params.Get("user_id")
	if userID, err = strconv.ParseInt(userIDStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svr.DelUser(c, userID))
}
