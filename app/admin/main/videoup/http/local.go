package http

import (
	"net/http"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// ping check server ok.
func ping(c *bm.Context) {

	if vdaSvc.Ping(c) != nil {
		log.Error("videoup-admin service ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
