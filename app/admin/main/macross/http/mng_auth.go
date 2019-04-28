package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// auth get auth by sysid.
func auth(c *bm.Context) {
	var (
		params = c.Request.Form
		system string
	)
	if system = params.Get("system"); system == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.Auth(c, system), nil)
}

// saveAuth save auth.
func saveAuth(c *bm.Context) {
	var (
		params             = c.Request.Form
		authID             int64
		authName, authFlag string
	)
	if authName = params.Get("auth_name"); authName == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	authIDStr := params.Get("auth_id")
	authID, _ = strconv.ParseInt(authIDStr, 10, 64)
	system := params.Get("system")
	if authID == 0 && system == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	authFlag = params.Get("auth_flag")
	if authID == 0 && system != "" && authFlag == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svr.SaveAuth(c, authID, system, authName, authFlag))
}

// delAuth del auth.
func delAuth(c *bm.Context) {
	var (
		params = c.Request.Form
		authID int64
		err    error
	)
	authIDStr := params.Get("auth_id")
	if authID, err = strconv.ParseInt(authIDStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svr.DelAuth(c, authID))
}

// authRelation update authRelation.
func authRelation(c *bm.Context) {
	var (
		params         = c.Request.Form
		roleID, authID int64
		state          int
		err            error
	)
	roleIDStr := params.Get("role_id")
	if roleID, err = strconv.ParseInt(roleIDStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	authIDStr := params.Get("auth_id")
	if authID, err = strconv.ParseInt(authIDStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	stateStr := params.Get("state")
	if state, err = strconv.Atoi(stateStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svr.AuthRelation(c, roleID, authID, state))
}
