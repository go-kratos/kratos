package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const (
	_sessUnKey = "username"
)

func queryUserInfo(c *bm.Context) {
	var (
		username string
		err      error
	)
	if username, err = getUsername(c); err != nil {
		return
	}

	c.JSON(srv.QueryUserInfo(c, username))

}

func syncWechatContact(c *bm.Context) {
	c.JSON(nil, srv.HttpSyncWechatContacts(c))
}

func updateVisibleBugly(c *bm.Context) {
	var (
		v = new(struct {
			UpdateUsername string `form:"username"`
			VisibleBugly   bool   `form:"visible_bugly"`
		})
		err      error
		username string
	)

	if err = c.Bind(v); err != nil {
		return
	}

	if username, err = getUsername(c); err != nil {
		return
	}
	c.JSON(nil, srv.UpdateUserVisibleBugly(c, username, v.UpdateUsername, v.VisibleBugly))
}

func accessToBugly(c *bm.Context) {
	var (
		username string
		err      error
	)
	if username, err = getUsername(c); err != nil {
		return
	}

	if !srv.AccessToBugly(c, username) {
		c.JSON(nil, ecode.AccessDenied)
		c.Abort()
		return
	}
}

func getUsername(c *bm.Context) (username string, err error) {
	user, exist := c.Get(_sessUnKey)
	if !exist {
		err = ecode.AccessKeyErr
		c.JSON(nil, err)
		return
	}
	username = user.(string)
	return
}
