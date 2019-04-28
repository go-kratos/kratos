package http

import (
	bm "go-common/library/net/http/blademaster"
)

func saveBugly2Tapd(c *bm.Context) {
	c.JSON(nil, srv.AsyncBuglyInsertTapd(c))
}

func updateBuglyStatusInTapd(c *bm.Context) {
	c.JSON(nil, srv.AsyncUpdateBuglyStatusInTapd(c))
}

func updateTitleInTapd(c *bm.Context) {
	c.JSON(nil, srv.AsyncUpdateBugInTapd(c))
}
