package http

import (
	"go-common/app/interface/openplatform/seo/model"
	bm "go-common/library/net/http/blademaster"
)

func itemInfo(c *bm.Context) {
	logUA(c)

	arg := new(model.ArgItemID)
	if err := c.Bind(arg); err != nil {
		return
	}

	bot := isBot(c)
	res, err := srv.GetItem(c, arg.ID, bot)
	if err != nil {
		c.String(0, "%v", err)
		return
	}

	if setCache(c, res) {
		return
	}
	c.String(0, string(res))
}
