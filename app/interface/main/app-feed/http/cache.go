package http

import (
	"encoding/json"

	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func upRcmd(c *bm.Context) {
	params := c.Request.Form
	item := params.Get("item")
	var is []*ai.Item
	if err := json.Unmarshal([]byte(item), &is); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, feedSvc.UpRcmdCache(c, is))
}
