package http

import (
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func joinAv(c *bm.Context) {
	v := new(struct {
		AccountType int `form:"account_type"`
		SignType    int `form:"sign_type"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	if err = svc.JoinAv(c, v.AccountType, mid, v.SignType); err != nil {
		log.Error("svc.Join accountType(%d) mid(%d) signType(%d) error(%v)", v.AccountType, mid, v.SignType, err)
	}
	c.JSON(nil, err)
}

func joinBgm(c *bm.Context) {
	v := new(struct {
		AccountType int `form:"account_type"`
		SignType    int `form:"sign_type"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	if err = svc.JoinBgm(c, mid, v.AccountType, v.SignType); err != nil {
		log.Error("svc.JoinBgm accountType(%d) mid(%d) signType(%d) error(%v)", v.AccountType, mid, v.SignType, err)
	}
	c.JSON(nil, err)
}
