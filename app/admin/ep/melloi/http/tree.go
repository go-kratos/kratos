package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func queryUserTree(c *bm.Context) {
	session, err := c.Request.Cookie("_AJSESSIONID")
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryUserTree(c, session.Value))
}

func queryTreeAdmin(c *bm.Context) {
	v := new(struct {
		Path string `form:"path"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	session, err := c.Request.Cookie("_AJSESSIONID")
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryTreeAdmin(c, v.Path, session.Value))
}
