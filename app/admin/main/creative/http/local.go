package http

import (
	"net/http"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func moPing(c *bm.Context) {
	var err error
	if err = svc.Ping(c); err != nil {
		log.Error("creative-admin ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
