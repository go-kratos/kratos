package http

import (
	"go-common/app/service/main/msm/model"
	bm "go-common/library/net/http/blademaster"
)

func codes(c *bm.Context) {
	var (
		err   error
		code  *model.Codes
		param = new(struct {
			Ver int64 `form:"ver"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if code, err = svr.Codes(c, param.Ver); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(code, nil)
}

func codesLangs(c *bm.Context) {
	var (
		err   error
		code  *model.CodesLangs
		param = new(struct {
			Ver int64 `form:"ver"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if code, err = svr.CodesLangs(c, param.Ver); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(code, nil)
}
