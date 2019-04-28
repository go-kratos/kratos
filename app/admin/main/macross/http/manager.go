package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// getAuths get user.
func getAuths(c *bm.Context) {
	var userName string
	username, _ := c.Get("username")
	userName, ok := username.(string)
	if !ok || userName == "" {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(svr.GetAuths(c, userName))
}
