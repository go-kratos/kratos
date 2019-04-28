package http

import (
	"net/http"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// ping check server ok.
func ping(c *bm.Context) {
	if err := rpSvr.Ping(c); err != nil {
		log.Error("reply interface ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
