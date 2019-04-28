package http

import (
	bm "go-common/library/net/http/blademaster"
)

func attentions(c *bm.Context) {
	var mid int64
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	c.JSON(webSvc.Attentions(c, mid))
}

func card(c *bm.Context) {
	var loginID int64
	v := new(struct {
		Mid      int64 `form:"mid" validate:"min=1"`
		TopPhoto bool  `form:"photo"`
		Article  bool  `form:"article"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	// login mid
	if loginIDStr, ok := c.Get("mid"); ok {
		loginID = loginIDStr.(int64)
	}
	c.JSON(webSvc.Card(c, v.Mid, loginID, v.TopPhoto, v.Article))
}
