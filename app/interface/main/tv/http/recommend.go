package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// get season's recommend data from pgc
func recommend(c *bm.Context) {
	req := c.Request.Form
	sid := req.Get("season_id")
	stype := req.Get("season_type")
	// param check
	if atoi(sid) == 0 || atoi(stype) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(tvSvc.RecomFilter(sid, stype))
}
