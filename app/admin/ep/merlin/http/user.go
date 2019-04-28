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

	c.JSON(svc.QueryUserInfo(c, username))

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
