package http

import (
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func showActivity(c *bm.Context) {
	v := new(struct {
		ActivityID int64 `form:"activity_id" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	data, err := svc.ShowActivity(c, mid, v.ActivityID)
	if err != nil {
		log.Error("svc.ShowActivity error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func signUpActivity(c *bm.Context) {
	var err error
	v := new(struct {
		ActivityID int64 `form:"activity_id" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	if err = svc.SignUpActivity(c, mid, v.ActivityID); err != nil {
		log.Error("svc.SignUpActivity error(%v)", err)
	}
	c.JSON(nil, err)
}
