package http

import (
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
)

func addAuthority(c *blademaster.Context) {
	v := new(struct {
		MIDS []int64 `form:"mids,split" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	num, err := svr.AddAuthority(c, v.MIDS)
	if err != nil {
		log.Error("svr.AddAuthority error(%v) arg(%v)", err, v)
		c.JSON(nil, err)
		return
	}
	c.JSON(num, nil)
}

func removeAuthority(c *blademaster.Context) {
	v := new(struct {
		MIDS []int64 `form:"mids,split" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	num, err := svr.RemoveAuthority(c, v.MIDS)
	if err != nil {
		log.Error("svr.RemoveAuthority error(%v) arg(%v)", err, v)
		c.JSON(nil, err)
		return
	}
	c.JSON(num, nil)
}
