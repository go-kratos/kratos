package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"strings"
)

func card(c *bm.Context) {
	c.JSON(artSrv.FindCard(c, c.Request.Form.Get("id")))
}

func cards(c *bm.Context) {
	var (
		params = c.Request.Form
	)
	ids := strings.Split(params.Get("ids"), ",")
	if len(ids) > 100 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(artSrv.FindCards(c, ids))
}
