package http

import (
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func errors(c *bm.Context) {
	type result struct {
		Error         string            `json:"error"`
		InstanceError map[string]string `json:"instance_error"`
	}
	res := result{
		Error:         cs.Error(),
		InstanceError: cs.Errors(),
	}
	c.JSON(res, nil)
}

func checkMaster(c *bm.Context) {
	arg := new(struct {
		Addr     string `form:"addr" validate:"required"`
		User     string `form:"user" validate:"required"`
		Password string `form:"password" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	name, pos, err := cs.CheckMaster(arg.Addr, arg.User, arg.Password)
	if err != nil {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	res := map[string]interface{}{"name": name, "pos:": pos}
	c.JSON(res, nil)
}

func syncPos(c *bm.Context) {
	arg := new(struct {
		Addr string `form:"addr" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		log.Error("syncpos params err %v", err)
		return
	}
	c.JSON(nil, cs.PosSync(arg.Addr))
}
