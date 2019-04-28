package http

import (
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"net/http"
)

func ping(c *bm.Context) {
	var err error
	if err = vdpSvc.Ping(c); err != nil {
		log.Error("videoup-interface ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
