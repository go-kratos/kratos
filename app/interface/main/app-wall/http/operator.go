package http

import (
	"time"

	bm "go-common/library/net/http/blademaster"
)

func reddot(c *bm.Context) {
	res := map[string]interface{}{
		"data": operatorSvc.Reddot(time.Now()),
	}
	returnDataJSON(c, res, nil)
}
