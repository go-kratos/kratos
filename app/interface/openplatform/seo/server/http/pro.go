package http

import (
	"go-common/app/interface/openplatform/seo/model"
	bm "go-common/library/net/http/blademaster"
)

func proList(c *bm.Context) {
	c.String(0, "project list")
}

func proInfo(c *bm.Context) {
	logUA(c)

	arg := new(model.ArgProID)
	if err := c.Bind(arg); err != nil {
		return
	}

	bot := isBot(c)
	res, err := srv.GetPro(c, arg.ID, bot)
	if err != nil {
		c.String(0, "%v", err)
		return
	}

	if setCache(c, res) {
		return
	}
	c.String(0, string(res))
}
