package http

import bm "go-common/library/net/http/blademaster"

// walls get GetWall
func walls(c *bm.Context) {
	res := map[string]interface{}{
		"data": wallSvc.Wall(),
	}
	returnDataJSON(c, res, nil)
}
