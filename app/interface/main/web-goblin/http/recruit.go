package http

import (
	"go-common/app/interface/main/web-goblin/model/web"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func recruit(c *bm.Context) {
	var (
		param = c.Request.Form
		err   error
		v     = &web.Params{}
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Mode == "social" || v.Mode == "campus" {
		c.JSON(srvWeb.Recruit(c, param, v))
		return
	}
	c.JSON("mode 只能为 社招(social) 或者 校招(campus）", ecode.RequestErr)
}
