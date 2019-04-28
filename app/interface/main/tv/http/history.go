package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// get history
func history(c *bm.Context) {
	v := new(struct {
		AccessKey string `form:"access_key"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if v.AccessKey == "" {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	if mid, ok := c.Get("mid"); !ok { // if not login, we don't call follow data
		c.JSON(nil, ecode.NoLogin)
	} else {
		c.JSON(hisSvc.GetHistory(c, mid.(int64)))
	}
}
