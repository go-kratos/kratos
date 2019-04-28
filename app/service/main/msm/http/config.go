package http

import (
	bm "go-common/library/net/http/blademaster"
)

// push config update
func push(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			Ver  int64  `form:"version" validate:"gte=0"`
			App  string `form:"service" validate:"required"`
			BVer string `form:"build_ver" validate:"required"`
			Env  string `form:"environment" validate:"required"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, svr.Push(c, param.App, param.BVer, param.Env, param.Ver))
}

func setToken(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			App   string `form:"service" validate:"required"`
			Token string `form:"token" validate:"required"`
			Env   string `form:"environment" validate:"required"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, svr.SetToken(c, param.App, param.Env, param.Token))
}
