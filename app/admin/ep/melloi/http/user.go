package http

import (
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func queryUser(c *bm.Context) {
	// 验证登录sessionID
	session, err := c.Request.Cookie("_AJSESSIONID")
	if err != nil {
		c.JSON(nil, ecode.AccessKeyErr)
		return
	}

	token, _ := srv.QueryServiceTreeToken(c, session.Value)
	if token == "" {
		c.JSON(nil, ecode.AccessKeyErr)
		return
	}

	// 获取用户名
	userName, err := c.Request.Cookie("username")
	if err != nil {
		c.JSON(nil, ecode.AccessKeyErr)
		return
	}
	c.JSON(srv.QueryUser(userName.Value))
}

func updateUser(c *bm.Context) {
	user := model.User{}
	if err := c.BindWith(&user, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, srv.UpdateUser(&user))
}
