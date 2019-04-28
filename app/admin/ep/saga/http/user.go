package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const (
	_sessUnKey = "username"
)

// @params EmptyReq
// @router get /ep/admin/saga/v1/user
// @response User
func queryUserInfo(ctx *bm.Context) {
	var (
		userName string
		err      error
	)
	if userName, err = getUsername(ctx); err != nil {
		return
	}
	ctx.JSON(srv.UserInfo(userName), nil)
}

func getUsername(c *bm.Context) (string, error) {
	if user, err := c.Request.Cookie(_sessUnKey); err == nil {
		return user.Value, nil
	}
	if value, exist := c.Get(_sessUnKey); exist {
		return value.(string), nil
	}
	return "", ecode.AccessKeyErr
}
