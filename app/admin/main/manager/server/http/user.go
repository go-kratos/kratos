package http

import (
	"net/http"
	"strconv"
	"strings"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

const _logout = "http://dashboard-mng.bilibili.co/logout?caller=manager-admin"

func authUser(c *bm.Context) {
	var (
		username string
	)
	if un, ok := c.Get("username"); ok {
		username = un.(string)
	} else {
		c.JSON(nil, ecode.Unauthorized)
		return
	}
	c.JSON(mngSvc.Auth(c, username))
}

func logout(c *bm.Context) {
	// purge mid cache
	c.Redirect(http.StatusFound, _logout)
}

func permissions(c *bm.Context) {
	var username string
	if username = c.Request.Form.Get("username"); username == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(mngSvc.Permissions(c, username))
}

func users(c *bm.Context) {
	var (
		err    error
		pn, ps int
		params = c.Request.Form
		pnStr  = params.Get("pn")
		psStr  = params.Get("ps")
	)
	if pn, err = strconv.Atoi(pnStr); err != nil || pn <= 0 {
		pn = 1
	}
	if ps, err = strconv.Atoi(psStr); err != nil || ps <= 0 {
		ps = 20
	}
	c.JSON(mngSvc.Users(c, pn, ps))
}

func usersTotal(c *bm.Context) {
	var (
		err   error
		total int64
	)
	if total, err = mngSvc.UsersTotal(c); err != nil {
		log.Error("mngSvc.UsersTotal error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(struct {
		Total int64 `json:"total"`
	}{total}, nil)
}

func heartbeat(c *bm.Context) {
	un, ok := c.Get("username")
	if !ok {
		log.Error("username not found in context")
		return
	}
	c.JSON(nil, mngSvc.Heartbeat(c, un.(string)))
}

// batch check unames
func usersNames(c *bm.Context) {
	var (
		err    error
		params = c.Request.Form
		uids   string
		uidsV  []int64
		items  map[int64]string
	)
	if uids = params.Get("uids"); uids == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if uidsV, err = xstr.SplitInts(uids); err != nil {
		log.Error("mngSvc.Unames(%s) error(%v)", uids, err)
		c.JSON(nil, err)
		return
	}
	items = mngSvc.Unames(c, uidsV)
	c.JSON(items, nil)
}

// batch check users' departments
func usersDepts(c *bm.Context) {
	var (
		err    error
		params = c.Request.Form
		uids   string
		uidsV  []int64
		items  map[int64]string
	)
	if uids = params.Get("uids"); uids == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if uidsV, err = xstr.SplitInts(uids); err != nil {
		log.Error("mngSvc.Udepts(%s) error(%v)", uids, err)
		c.JSON(nil, err)
		return
	}
	items = mngSvc.Udepts(c, uidsV)
	c.JSON(items, nil)
}

func userIds(c *bm.Context) {
	var (
		params  = c.Request.Form
		unames  string
		unamesV []string
		items   map[string]int64
	)
	if unames = params.Get("unames"); unames == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	unamesV = strings.Split(unames, ",")
	items = mngSvc.UIds(c, unamesV)
	c.JSON(items, nil)
}
